package dto

import (
	"d-payroll/entity"
	"time"
)

type PayslipAttendanceDetailDto struct {
	CheckinAt     time.Time  `json:"checkin_at"`
	CheckoutAt    *time.Time `json:"checkout_at"`
	DurationMilis int        `json:"duration_milis"`
}

func (p *PayslipAttendanceDetailDto) FromPayslipAttendanceDetailEntity(attendance *entity.PayslipAttendanceDetail) {
	p.CheckinAt = attendance.CheckinAt
	p.CheckoutAt = attendance.CheckoutAt
	p.DurationMilis = attendance.DurationMilis
}

type PayslipAttendanceDto struct {
	Details            []*PayslipAttendanceDetailDto `json:"details"`
	TotalDurationMilis int                           `json:"total_duration_milis"`
	TotalAmount        float32                       `json:"total_amount"`
}

func (p *PayslipAttendanceDto) FromPayslipAttendanceEntity(attendance *entity.PayslipAttendance) {
	p.TotalDurationMilis = attendance.TotalDurationMilis
	p.TotalAmount = attendance.TotalAmount

	details := make([]*PayslipAttendanceDetailDto, len(attendance.Details))
	for i, detail := range attendance.Details {
		dto := &PayslipAttendanceDetailDto{}
		dto.FromPayslipAttendanceDetailEntity(detail)
		details[i] = dto
	}
	p.Details = details
}

type PayslipOvertimeDetailDto struct {
	OvertimeAt    time.Time `json:"overtime_at"`
	Description   string    `json:"description"`
	DurationMilis int       `json:"duration_milis"`
	CreatedAt     time.Time `json:"created_at"`
}

func (p *PayslipOvertimeDetailDto) FromPayslipOvertimeDetailEntity(overtime *entity.PayslipOvertimeDetail) {
	p.OvertimeAt = overtime.OvertimeAt
	p.Description = overtime.Description
	p.DurationMilis = overtime.DurationMilis
	p.CreatedAt = overtime.CreatedAt
}

type PayslipOvertimeDto struct {
	Details            []*PayslipOvertimeDetailDto `json:"details"`
	TotalDurationMilis int                         `json:"total_duration_milis"`
	TotalAmount        float32                     `json:"total_amount"`
}

func (p *PayslipOvertimeDto) FromPayslipOvertimeEntity(overtime *entity.PayslipOvertime) {
	p.TotalDurationMilis = overtime.TotalDurationMilis
	p.TotalAmount = overtime.TotalAmount

	details := make([]*PayslipOvertimeDetailDto, len(overtime.Details))
	for i, detail := range overtime.Details {
		dto := &PayslipOvertimeDetailDto{}
		dto.FromPayslipOvertimeDetailEntity(detail)
		details[i] = dto
	}
	p.Details = details
}

type PayslipReimburseDetailDto struct {
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
}

func (p *PayslipReimburseDetailDto) FromPayslipReimburseDetailEntity(reimburse *entity.PayslipReimburseDetail) {
	p.Description = reimburse.Description
	p.Amount = reimburse.Amount
	p.CreatedAt = reimburse.CreatedAt
}

type PayslipReimburseDto struct {
	Details     []*PayslipReimburseDetailDto `json:"details"`
	TotalAmount float32                      `json:"total_amount"`
}

func (p *PayslipReimburseDto) FromPayslipReimburseEntity(reimburse *entity.PayslipReimburse) {
	p.TotalAmount = reimburse.TotalAmount

	details := make([]*PayslipReimburseDetailDto, len(reimburse.Details))
	for i, detail := range reimburse.Details {
		dto := &PayslipReimburseDetailDto{}
		dto.FromPayslipReimburseDetailEntity(detail)
		details[i] = dto
	}
	p.Details = details
}

type PayslipDto struct {
	PayrollID   uint                  `json:"payroll_id"`
	UserID      uint                  `json:"user_id"`
	Salary      int                   `json:"salary"`
	ProRate     float32               `json:"pro_rate"`
	Attendance  *PayslipAttendanceDto `json:"attendance"`
	Overtime    *PayslipOvertimeDto   `json:"overtime"`
	Reimburse   *PayslipReimburseDto  `json:"reimburse"`
	TakeHomePay float32               `json:"take_home_pay"`
}

func (p *PayslipDto) FromPayslipEntity(payslip *entity.Payslip) {
	p.PayrollID = payslip.PayrollID
	p.UserID = payslip.UserID
	p.Salary = payslip.Salary
	p.ProRate = payslip.ProRate

	if payslip.Attendance != nil {
		p.Attendance = &PayslipAttendanceDto{}
		p.Attendance.FromPayslipAttendanceEntity(payslip.Attendance)
	}
	if payslip.Overtime != nil {
		p.Overtime = &PayslipOvertimeDto{}
		p.Overtime.FromPayslipOvertimeEntity(payslip.Overtime)
	}
	if payslip.Reimburse != nil {
		p.Reimburse = &PayslipReimburseDto{}
		p.Reimburse.FromPayslipReimburseEntity(payslip.Reimburse)
	}
	p.TakeHomePay = payslip.TakeHomePay
}

type UserPayslipSummaryDto struct {
	PayrollID        uint      `json:"payroll_id"`
	UserID           uint      `json:"user_id"`
	TotalTakeHomePay int       `json:"total_take_home_pay"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (u *UserPayslipSummaryDto) FromUserPayslipSummaryEntity(summary *entity.UserPayslipSummary) {
	u.PayrollID = summary.PayrollID
	u.UserID = summary.UserID
	u.TotalTakeHomePay = summary.TotalTakeHomePay
	u.CreatedAt = *summary.CreatedAt
	u.UpdatedAt = *summary.UpdatedAt
}
