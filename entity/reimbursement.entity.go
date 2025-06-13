package entity

import "time"

type UserReimbursement struct {
	ID               *uint
	UserID           uint
	Description      string
	Amount           int
	ApprovedByUserID *uint
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
}
