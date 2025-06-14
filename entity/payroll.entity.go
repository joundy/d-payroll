package entity

import "time"

type Payroll struct {
	ID              *uint
	Name            string
	StartedAt       time.Time
	EndedAt         time.Time
	IsRolled        *bool
	UpdatedByUserID *uint
	CreatedByUserID *uint
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

type UserPayslipSummary struct {
	ID               *uint
	PayrollID        uint
	UserID           uint
	TotalTakeHomePay int
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
}
