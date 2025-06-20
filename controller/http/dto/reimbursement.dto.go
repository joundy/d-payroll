package dto

import (
	"d-payroll/entity"
	"time"
)

type CreateReimbursementBodyDto struct {
	Description string `json:"description" validate:"required"`
	Amount      int    `json:"amount" validate:"required,min=1"`
}

func (c *CreateReimbursementBodyDto) ToReimbursementEntity(userID uint) *entity.UserReimbursement {
	return &entity.UserReimbursement{
		UserID:      userID,
		Description: c.Description,
		Amount:      c.Amount,
	}
}

type ReimbursementResponseDto struct {
	ID              *uint      `json:"id"`
	UserID          uint       `json:"user_id"`
	Description     string     `json:"description"`
	Amount          int        `json:"amount"`
	IsApproved      bool       `json:"is_approved"`
	UpdatedByUserID *uint      `json:"updated_by_user_id"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

func (r *ReimbursementResponseDto) FromReimbursementEntity(reimbursement *entity.UserReimbursement) {
	r.ID = reimbursement.ID
	r.UserID = reimbursement.UserID
	r.Description = reimbursement.Description
	r.Amount = reimbursement.Amount
	r.IsApproved = reimbursement.IsApproved
	r.UpdatedByUserID = reimbursement.UpdatedByUserID
	r.CreatedAt = reimbursement.CreatedAt
	r.UpdatedAt = reimbursement.UpdatedAt
}
