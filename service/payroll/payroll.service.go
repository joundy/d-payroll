package payrollservice

import (
	"context"
	"d-payroll/config"
	"d-payroll/entity"
	internalerror "d-payroll/internal-error"
	repository "d-payroll/repository/db"
	"d-payroll/repository/db/models"
)

type PayrollService interface {
	CreatePayroll(ctx context.Context, payroll *entity.Payroll) (*entity.Payroll, error)
	GetPayrolls(ctx context.Context) ([]*entity.Payroll, error)
	RollPayroll(ctx context.Context, payrollID uint, userID uint) error
}

type payrollService struct {
	config    *config.Config
	payrollDB repository.PayrollDB
}

func NewPayrollService(config *config.Config, payrollDB repository.PayrollDB) PayrollService {
	return &payrollService{
		config:    config,
		payrollDB: payrollDB,
	}
}

func (s *payrollService) CreatePayroll(ctx context.Context, payroll *entity.Payroll) (*entity.Payroll, error) {
	payrollModel := &models.Payroll{}
	payrollModel.FromPayrollEntity(payroll)

	err := s.payrollDB.CreatePayroll(ctx, payrollModel)
	if err != nil {
		return nil, err
	}

	return payrollModel.ToPayrollEntity(), nil
}

func (s *payrollService) GetPayrolls(ctx context.Context) ([]*entity.Payroll, error) {
	payrollModels, err := s.payrollDB.GetPayrolls(ctx)
	if err != nil {
		return nil, err
	}

	payrolls := make([]*entity.Payroll, len(payrollModels))
	for i, payrollModel := range payrollModels {
		payrolls[i] = payrollModel.ToPayrollEntity()
	}

	return payrolls, nil
}

func (s *payrollService) RollPayroll(ctx context.Context, payrollID uint, userID uint) error {
	payroll, err := s.payrollDB.GetPayrollByID(ctx, payrollID)
	if err != nil {
		return err
	}

	if payroll.IsRolled != nil && *payroll.IsRolled {
		return &internalerror.PayrollAlreadyRolledError{}
	}

	return s.payrollDB.RollPayroll(ctx, payrollID, userID)
}
