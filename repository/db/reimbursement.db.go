package repository

import (
	"context"
	internalerror "d-payroll/internal-error"
	"d-payroll/repository/db/models"
	"errors"

	"gorm.io/gorm"
)

type ReimbursementDB interface {
	CreateReimbursement(ctx context.Context, reimbursement *models.UserReimbursement) error
	ApproveReimbursement(ctx context.Context, reimbursementID uint, approvedByUserID uint) error
	GetReimbursementsByUserID(ctx context.Context, userID uint) ([]*models.UserReimbursement, error)
	GetReimbursementByID(ctx context.Context, reimbursementID uint) (*models.UserReimbursement, error)
}

type reimbursementDB struct {
	DB *gorm.DB
}

func NewReimbursementDB(db *gorm.DB) ReimbursementDB {
	return &reimbursementDB{DB: db}
}

func (r *reimbursementDB) CreateReimbursement(ctx context.Context, reimbursement *models.UserReimbursement) error {
	return r.DB.WithContext(ctx).Create(reimbursement).Error
}

func (r *reimbursementDB) ApproveReimbursement(ctx context.Context, reimbursementID uint, updatedByUserID uint) error {
	return r.DB.WithContext(ctx).
		Model(&models.UserReimbursement{}).
		Where("id = ?", reimbursementID).
		Updates(map[string]interface{}{
			"is_approved":        true,
			"updated_by_user_id": updatedByUserID,
		}).Error
}

func (r *reimbursementDB) GetReimbursementsByUserID(ctx context.Context, userID uint) ([]*models.UserReimbursement, error) {
	var reimbursements []*models.UserReimbursement
	result := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&reimbursements)
	if result.Error != nil {
		return nil, result.Error
	}
	return reimbursements, nil
}

func (r *reimbursementDB) GetReimbursementByID(ctx context.Context, reimbursementID uint) (*models.UserReimbursement, error) {
	var reimbursement *models.UserReimbursement

	result := r.DB.WithContext(ctx).Where("id = ?", reimbursementID).First(&reimbursement)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &internalerror.NotFoundError{}
		}
		return nil, result.Error
	}

	return reimbursement, nil
}
