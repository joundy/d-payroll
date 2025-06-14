package overtimeservice

import (
	"context"
	"d-payroll/config"
	"d-payroll/entity"
	internalerror "d-payroll/internal-error"
	repository "d-payroll/repository/db"
	"d-payroll/repository/db/models"
	attendanceservice "d-payroll/service/attendance"
	"d-payroll/utils"
)

type OvertimeService interface {
	CreateOvertime(ctx context.Context, overtime *entity.UserOvertime) (*entity.UserOvertime, error)
	ApproveOvertime(ctx context.Context, overtimeID uint, approvedByUserID uint) error
	GetOvertimesByUserID(ctx context.Context, userID uint) ([]*entity.UserOvertime, error)
}

type overtimeService struct {
	config        *config.Config
	overtimeDB    repository.OvertimeDB
	attendanceSvc attendanceservice.AttendanceService
}

func NewOvertimeService(config *config.Config, overtimeDB repository.OvertimeDB, attendanceSvc attendanceservice.AttendanceService) OvertimeService {
	return &overtimeService{
		config:        config,
		overtimeDB:    overtimeDB,
		attendanceSvc: attendanceSvc,
	}
}

func (s *overtimeService) CreateOvertime(ctx context.Context, overtime *entity.UserOvertime) (*entity.UserOvertime, error) {
	if !utils.IsWeekend() {

		isCheckedOut, err := s.attendanceSvc.IsCheckedOut(ctx, overtime.UserID)
		if err != nil {
			return nil, err
		}

		if !isCheckedOut {
			return nil, &internalerror.OvertimeSubmitBeforeCheckoutError{}
		}
	}

	thisDayOvertimes, err := s.overtimeDB.GetThisDayOvertimeByUserID(ctx, overtime.UserID)
	if err != nil {
		return nil, err
	}
	totalMilis := 0
	for _, overtime := range thisDayOvertimes {
		totalMilis += overtime.DurationMilis
	}

	if totalMilis+overtime.DurationMilis > s.config.Overtime.MaxDurationPerDayMilis {
		return nil, &internalerror.OvertimeExceedsLimitError{}
	}

	overtimeModel := &models.UserOvertime{}
	overtimeModel.FromOvertimeEntity(overtime)

	err = s.overtimeDB.CreateOvertime(ctx, overtimeModel)
	if err != nil {
		return nil, err
	}

	return overtimeModel.ToOvertimeEntity(), nil
}

func (s *overtimeService) ApproveOvertime(ctx context.Context, overtimeID uint, approvedByUserID uint) error {
	overtime, err := s.overtimeDB.GetOvertimeByID(ctx, overtimeID)
	if err != nil {
		return err
	}
	if overtime.IsApproved {
		return &internalerror.OvertimeAlreadyApprovedError{}
	}

	return s.overtimeDB.ApproveOvertime(ctx, overtimeID, approvedByUserID)
}

func (s *overtimeService) GetOvertimesByUserID(ctx context.Context, userID uint) ([]*entity.UserOvertime, error) {
	overtimeModels, err := s.overtimeDB.GetOvertimesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	overtimes := make([]*entity.UserOvertime, len(overtimeModels))
	for i, model := range overtimeModels {
		overtimes[i] = model.ToOvertimeEntity()
	}

	return overtimes, nil
}
