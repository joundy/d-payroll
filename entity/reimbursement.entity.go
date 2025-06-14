package entity

import "time"

type UserReimbursement struct {
	ID              *uint
	UserID          uint
	Description     string
	Amount          int
	IsApproved      bool
	UpdatedByUserID *uint
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}
