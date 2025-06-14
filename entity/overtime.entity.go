package entity

import "time"

type UserOvertime struct {
	ID               *uint
	UserID           uint
	Description      string
	OvertimeAt       time.Time
	DurationMilis    int
	ApprovedByUserID *uint
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
}
