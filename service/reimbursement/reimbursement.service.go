package reimbursementservice

import (
	"context"
	"d-payroll/entity"
	internalerror "d-payroll/internal-error"
	repository "d-payroll/repository/db"
	"d-payroll/repository/db/models"
)

type ReimbursementService interface {
	CreateReimbursement(ctx context.Context, reimbursement *entity.UserReimbursement) (*entity.UserReimbursement, error)
	ApproveReimbursement(ctx context.Context, reimbursementID uint, approvedByUserID uint) error
	GetReimbursementsByUserID(ctx context.Context, userID uint) ([]*entity.UserReimbursement, error)
}

type reimbursementService struct {
	reimbursementDB repository.ReimbursementDB
}

func NewReimbursementService(reimbursementDB repository.ReimbursementDB) ReimbursementService {
	return &reimbursementService{
		reimbursementDB: reimbursementDB,
	}
}

func (s *reimbursementService) CreateReimbursement(ctx context.Context, reimbursement *entity.UserReimbursement) (*entity.UserReimbursement, error) {
	reimbursementModel := &models.UserReimbursement{}
	reimbursementModel.FromReimbursementEntity(reimbursement)

	err := s.reimbursementDB.CreateReimbursement(ctx, reimbursementModel)
	if err != nil {
		return nil, err
	}

	return reimbursementModel.ToReimbursementEntity(), nil
}

func (s *reimbursementService) ApproveReimbursement(ctx context.Context, reimbursementID uint, approvedByUserID uint) error {
	reimbursement, err := s.reimbursementDB.GetReimbursementByID(ctx, reimbursementID)
	if err != nil {
		return err
	}
	if reimbursement.ApprovedByUserID != nil {
		return &internalerror.ReimbursementAlreadyApprovedError{}
	}

	return s.reimbursementDB.ApproveReimbursement(ctx, reimbursementID, approvedByUserID)
}

func (s *reimbursementService) GetReimbursementsByUserID(ctx context.Context, userID uint) ([]*entity.UserReimbursement, error) {
	reimbursementModels, err := s.reimbursementDB.GetReimbursementsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	reimbursements := make([]*entity.UserReimbursement, len(reimbursementModels))
	for i, model := range reimbursementModels {
		reimbursements[i] = model.ToReimbursementEntity()
	}

	return reimbursements, nil
}
