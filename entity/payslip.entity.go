package entity

import "time"

type PayslipOvertimeDetail struct {
	OvertimeAt    time.Time
	Description   string
	DurationMilis int
	CreatedAt     time.Time
}

type PayslipOvertime struct {
	Details            []*PayslipOvertimeDetail
	TotalDurationMilis int
	TotalAmount        float32
}

type PayslipReimburseDetail struct {
	Description string
	Amount      int
	CreatedAt   time.Time
}

type PayslipReimburse struct {
	Details     []*PayslipReimburseDetail
	TotalAmount float32
}

type PayslipAttendanceDetail struct {
	CheckinAt     time.Time
	CheckoutAt    *time.Time
	DurationMilis int
}

type PayslipAttendance struct {
	Details            []*PayslipAttendanceDetail
	TotalDurationMilis int
	TotalAmount        float32
}

type Payslip struct {
	PayrollID   uint
	UserID      uint
	Salary      int
	ProRate     float32
	Attendance  *PayslipAttendance
	Overtime    *PayslipOvertime
	Reimburse   *PayslipReimburse
	TakeHomePay float32
}
