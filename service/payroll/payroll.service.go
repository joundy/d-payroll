package payrollservice

import (
	"context"
	"d-payroll/config"
	"d-payroll/entity"
	internalerror "d-payroll/internal-error"
	repository "d-payroll/repository/db"
	"d-payroll/repository/db/models"
	attendanceservice "d-payroll/service/attendance"
	overtimeservice "d-payroll/service/overtime"
	reimbursementservice "d-payroll/service/reimbursement"
	userservice "d-payroll/service/user"
	"fmt"
	"time"
)

type PayrollService interface {
	CreatePayroll(ctx context.Context, payroll *entity.Payroll) (*entity.Payroll, error)
	GetPayrolls(ctx context.Context) ([]*entity.Payroll, error)
	RollPayroll(ctx context.Context, payrollID uint, userID uint) error

	GeneratePayslip(ctx context.Context, payrollID uint, userID uint) (*entity.Payslip, error)
}

type payrollService struct {
	config    *config.Config
	payrollDB repository.PayrollDB

	userservice          userservice.UserService
	attendanceService    attendanceservice.AttendanceService
	reimbursementService reimbursementservice.ReimbursementService
	overtimeService      overtimeservice.OvertimeService
}

func NewPayrollService(config *config.Config, payrollDB repository.PayrollDB, userservice userservice.UserService, attendanceService attendanceservice.AttendanceService, reimbursementService reimbursementservice.ReimbursementService, overtimeService overtimeservice.OvertimeService) PayrollService {
	return &payrollService{
		config:    config,
		payrollDB: payrollDB,

		userservice:          userservice,
		attendanceService:    attendanceService,
		reimbursementService: reimbursementService,
		overtimeService:      overtimeService,
	}
}

func (s *payrollService) CreatePayroll(ctx context.Context, payroll *entity.Payroll) (*entity.Payroll, error) {
	payrollModel := &models.Payroll{}
	payrollModel.FromPayrollEntity(payroll)

	err := s.payrollDB.CreatePayroll(ctx, payrollModel)
	if err != nil {
		return nil, err
	}

	return payrollModel.ToPayrollEntity(), nil
}

func (s *payrollService) GetPayrolls(ctx context.Context) ([]*entity.Payroll, error) {
	payrollModels, err := s.payrollDB.GetPayrolls(ctx)
	if err != nil {
		return nil, err
	}

	payrolls := make([]*entity.Payroll, len(payrollModels))
	for i, payrollModel := range payrollModels {
		payrolls[i] = payrollModel.ToPayrollEntity()
	}

	return payrolls, nil
}

func (s *payrollService) RollPayroll(ctx context.Context, payrollID uint, userID uint) error {
	payroll, err := s.payrollDB.GetPayrollByID(ctx, payrollID)
	if err != nil {
		return err
	}

	if payroll.IsRolled != nil && *payroll.IsRolled {
		return &internalerror.PayrollAlreadyRolledError{}
	}

	return s.payrollDB.RollPayroll(ctx, payrollID, userID)
}

// TODO: this should be cached, not ideal, shoud lock the database (maybe SHARE restriction is enough)
func (s *payrollService) GeneratePayslip(ctx context.Context, payrollID uint, userID uint) (*entity.Payslip, error) {
	payroll, err := s.payrollDB.GetPayrollByID(ctx, payrollID)
	if err != nil {
		return nil, err
	}

	if payroll.IsRolled == nil {
		return nil, &internalerror.PayrollNotRolledError{}
	}

	user, err := s.userservice.GetUserById(ctx, userID)
	if err != nil {
		return nil, err
	}

	attendancesGroup, err := s.attendanceService.GetAttendancesByUserIDAndDateBetweenGroupByDate(ctx, userID, payroll.StartedAt, payroll.UpdatedAt)
	if err != nil {
		return nil, err
	}

	var reimbursementDetails []*entity.PayslipReimburseDetail
	reimbursements, err := s.reimbursementService.GetReimbursementsByUserIDAndDateBetween(ctx, userID, payroll.StartedAt, payroll.UpdatedAt)
	if err != nil {
		return nil, err
	}
	for _, reimbursement := range reimbursements {
		if reimbursement.IsApproved {
			reimbursementDetails = append(reimbursementDetails, &entity.PayslipReimburseDetail{
				Description: reimbursement.Description,
				Amount:      reimbursement.Amount,
				CreatedAt:   *reimbursement.CreatedAt,
			})
		}
	}

	var overtimeDetails []*entity.PayslipOvertimeDetail
	overtimes, err := s.overtimeService.GetOvertimesByUserIDAndDateBetween(ctx, userID, payroll.StartedAt, payroll.UpdatedAt)
	if err != nil {
		return nil, err
	}
	for _, overtime := range overtimes {
		if overtime.IsApproved {
			overtimeDetails = append(overtimeDetails, &entity.PayslipOvertimeDetail{
				OvertimeAt:    *overtime.CreatedAt,
				Description:   overtime.Description,
				DurationMilis: overtime.DurationMilis,
				CreatedAt:     *overtime.CreatedAt,
			})

		}
	}

	if user.UserInfo == nil {
		return nil, fmt.Errorf("User info not found")
	}
	salary := user.UserInfo.MonthlySalary
	// TODO: precision issuee heree..
	proRateMilis := float32(*salary) / float32(s.config.Payroll.DayPerMonthProrate*s.config.Payroll.MaxWorkingMilisPerDay)

	attendanceDetails := []*entity.PayslipAttendanceDetail{}
	for _, attendance := range attendancesGroup {
		// checkinAt is the time of checkin or the end of the payroll if checkin is nil
		// it's possible the started at of the payroll is after the checkin (edge case)
		checkinAt := payroll.UpdatedAt
		if attendance.CheckIn != nil {
			checkinAt = *attendance.CheckIn.CreatedAt
		}

		// if the checkout is nil, it means the user forgot to checkout, so we use max working milis per day
		durationMilis := s.config.Payroll.MaxWorkingMilisPerDay
		var checkoutAt *time.Time
		if attendance.CheckOut != nil {
			checkoutAt = attendance.CheckOut.CreatedAt
			if attendance.CheckIn != nil && attendance.CheckOut != nil && attendance.CheckOut.CreatedAt != nil {
				durationMilis = int(attendance.CheckOut.CreatedAt.Sub(checkinAt).Milliseconds())
			}
		}

		// cap the duration to max working milis per day
		if durationMilis > s.config.Payroll.MaxWorkingMilisPerDay {
			durationMilis = s.config.Payroll.MaxWorkingMilisPerDay
		}

		attendanceDetails = append(attendanceDetails, &entity.PayslipAttendanceDetail{
			CheckinAt:     checkinAt,
			CheckoutAt:    checkoutAt,
			DurationMilis: durationMilis,
		})
	}

	attendanceTotalDurationMilis := 0
	var attendanceTotalAmount float32
	for _, attendance := range attendanceDetails {
		attendanceTotalDurationMilis += attendance.DurationMilis
		attendanceTotalAmount += float32(attendance.DurationMilis) * proRateMilis
	}
	attendance := &entity.PayslipAttendance{
		Details:            attendanceDetails,
		TotalDurationMilis: attendanceTotalDurationMilis,
		TotalAmount:        attendanceTotalAmount,
	}

	var reimburseTotalAmount float32
	for _, reimbursement := range reimbursements {
		reimburseTotalAmount += float32(reimbursement.Amount)
	}
	reimburse := &entity.PayslipReimburse{
		Details:     reimbursementDetails,
		TotalAmount: reimburseTotalAmount,
	}

	overtimeTotalDurationMilis := 0
	var overtimeTotalAmount float32
	for _, overtime := range overtimes {
		overtimeTotalDurationMilis += overtime.DurationMilis
		overtimeTotalAmount += float32(overtime.DurationMilis) * proRateMilis
	}
	overtime := &entity.PayslipOvertime{
		Details:            overtimeDetails,
		TotalDurationMilis: overtimeTotalDurationMilis,
		TotalAmount:        overtimeTotalAmount,
	}

	payslip := &entity.Payslip{
		PayrollID:   payroll.ID,
		UserID:      userID,
		Salary:      *salary,
		ProRate:     proRateMilis,
		Attendance:  attendance,
		Overtime:    overtime,
		Reimburse:   reimburse,
		TakeHomePay: attendanceTotalAmount + overtimeTotalAmount + reimburseTotalAmount,
	}

	return payslip, nil
}
