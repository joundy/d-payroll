package models

import (
	"d-payroll/entity"
	"time"

	"gorm.io/gorm"
)

type UserReimbursement struct {
	gorm.Model

	UserID          uint
	User            *User `gorm:"foreignKey:UserID"`
	Description     string
	Amount          int
	IsApproved      bool `gorm:"default:false"`
	UpdatedByUserID *uint
	UpdatedByUser   *User `gorm:"foreignKey:UpdatedByUserID"`
}

func (u *UserReimbursement) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return
}

func (u *UserReimbursement) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return
}

func (r *UserReimbursement) ToReimbursementEntity() *entity.UserReimbursement {
	return &entity.UserReimbursement{
		ID:              &r.ID,
		UserID:          r.UserID,
		Description:     r.Description,
		Amount:          r.Amount,
		IsApproved:      r.IsApproved,
		UpdatedByUserID: r.UpdatedByUserID,
		CreatedAt:       &r.CreatedAt,
		UpdatedAt:       &r.UpdatedAt,
	}
}

func (r *UserReimbursement) FromReimbursementEntity(reimbursement *entity.UserReimbursement) {
	r.UserID = reimbursement.UserID
	r.Description = reimbursement.Description
	r.Amount = reimbursement.Amount
	r.IsApproved = reimbursement.IsApproved
	r.UpdatedByUserID = reimbursement.UpdatedByUserID

	if reimbursement.CreatedAt != nil {
		r.CreatedAt = *reimbursement.CreatedAt
	}

	if reimbursement.UpdatedAt != nil {
		r.UpdatedAt = *reimbursement.UpdatedAt
	}
}
