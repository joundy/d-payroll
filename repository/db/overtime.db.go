package repository

import (
	"context"
	internalerror "d-payroll/internal-error"
	"d-payroll/repository/db/models"
	"errors"

	"gorm.io/gorm"
)

type OvertimeDB interface {
	CreateOvertime(ctx context.Context, overtime *models.UserOvertime) error
	ApproveOvertime(ctx context.Context, overtimeID uint, approvedByUserID uint) error
	GetOvertimesByUserID(ctx context.Context, userID uint) ([]*models.UserOvertime, error)
	GetOvertimeByID(ctx context.Context, overtimeID uint) (*models.UserOvertime, error)
}

type overtimeDB struct {
	DB *gorm.DB
}

func NewOvertimeDB(db *gorm.DB) OvertimeDB {
	return &overtimeDB{DB: db}
}

func (o *overtimeDB) CreateOvertime(ctx context.Context, overtime *models.UserOvertime) error {
	return o.DB.WithContext(ctx).Create(overtime).Error
}

func (o *overtimeDB) ApproveOvertime(ctx context.Context, overtimeID uint, approvedByUserID uint) error {
	return o.DB.WithContext(ctx).Model(&models.UserOvertime{}).Where("id = ?", overtimeID).Update("approved_by_user_id", approvedByUserID).Error
}

func (o *overtimeDB) GetOvertimesByUserID(ctx context.Context, userID uint) ([]*models.UserOvertime, error) {
	var overtimes []*models.UserOvertime
	result := o.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&overtimes)
	if result.Error != nil {
		return nil, result.Error
	}
	return overtimes, nil
}

func (o *overtimeDB) GetOvertimeByID(ctx context.Context, overtimeID uint) (*models.UserOvertime, error) {
	var overtime *models.UserOvertime

	result := o.DB.WithContext(ctx).Where("id = ?", overtimeID).First(&overtime)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &internalerror.NotFoundError{}
		}
		return nil, result.Error
	}

	return overtime, nil
}
