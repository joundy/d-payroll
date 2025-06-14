package entity

import "time"

type UserOvertime struct {
	ID              *uint
	UserID          uint
	Description     string
	OvertimeAt      time.Time
	DurationMilis   int
	IsApproved      bool
	UpdatedByUserID *uint
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}
