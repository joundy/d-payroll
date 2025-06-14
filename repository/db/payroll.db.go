package repository

import (
	"context"
	"errors"

	internalerror "d-payroll/internal-error"
	"d-payroll/repository/db/models"

	"gorm.io/gorm"
)

type PayrollDB interface {
	CreatePayroll(ctx context.Context, payroll *models.Payroll) error
	GetPayrollByID(ctx context.Context, payrollID uint) (*models.Payroll, error)
	GetPayrolls(ctx context.Context) ([]*models.Payroll, error)
	RollPayroll(ctx context.Context, payrollID uint, userID uint) error
}

type payrollDB struct {
	DB *gorm.DB
}

func NewPayrollDB(db *gorm.DB) PayrollDB {
	return &payrollDB{DB: db}
}

func (p *payrollDB) CreatePayroll(ctx context.Context, payroll *models.Payroll) error {
	return p.DB.WithContext(ctx).Create(payroll).Error
}

func (p *payrollDB) GetPayrollByID(ctx context.Context, payrollID uint) (*models.Payroll, error) {
	var payroll *models.Payroll

	result := p.DB.WithContext(ctx).Where("id = ?", payrollID).First(&payroll)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &internalerror.NotFoundError{}
		}
		return nil, result.Error
	}

	return payroll, nil
}

func (p *payrollDB) GetPayrolls(ctx context.Context) ([]*models.Payroll, error) {
	var payrolls []*models.Payroll
	if err := p.DB.WithContext(ctx).Find(&payrolls).Error; err != nil {
		return nil, err
	}

	return payrolls, nil
}

func (p *payrollDB) RollPayroll(ctx context.Context, payrollID uint, userID uint) error {
	return p.DB.WithContext(ctx).Model(&models.Payroll{}).
		Where("id = ?", payrollID).
		Updates(map[string]interface{}{
			"is_rolled":          true,
			"updated_by_user_id": userID,
		}).Error
}
