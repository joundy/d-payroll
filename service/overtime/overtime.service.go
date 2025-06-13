package overtimeservice

import (
	"context"
	"d-payroll/entity"
	internalerror "d-payroll/internal-error"
	repository "d-payroll/repository/db"
	"d-payroll/repository/db/models"
)

type OvertimeService interface {
	CreateOvertime(ctx context.Context, overtime *entity.UserOvertime) (*entity.UserOvertime, error)
	ApproveOvertime(ctx context.Context, overtimeID uint, approvedByUserID uint) error
	GetOvertimesByUserID(ctx context.Context, userID uint) ([]*entity.UserOvertime, error)
}

type overtimeService struct {
	overtimeDB repository.OvertimeDB
}

func NewOvertimeService(overtimeDB repository.OvertimeDB) OvertimeService {
	return &overtimeService{
		overtimeDB: overtimeDB,
	}
}

func (s *overtimeService) CreateOvertime(ctx context.Context, overtime *entity.UserOvertime) (*entity.UserOvertime, error) {

	overtimeModel := &models.UserOvertime{}
	overtimeModel.FromOvertimeEntity(overtime)

	err := s.overtimeDB.CreateOvertime(ctx, overtimeModel)
	if err != nil {
		return nil, err
	}

	return overtimeModel.ToOvertimeEntity(), nil
}

func (s *overtimeService) ApproveOvertime(ctx context.Context, overtimeID uint, approvedByUserID uint) error {
	overtime, err := s.overtimeDB.GetOvertimeByID(ctx, overtimeID)
	if err != nil {
		return err
	}
	if overtime.ApprovedByUserID != nil {
		return &internalerror.OvertimeAlreadyApprovedError{}
	}

	return s.overtimeDB.ApproveOvertime(ctx, overtimeID, approvedByUserID)
}

func (s *overtimeService) GetOvertimesByUserID(ctx context.Context, userID uint) ([]*entity.UserOvertime, error) {
	overtimeModels, err := s.overtimeDB.GetOvertimesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	overtimes := make([]*entity.UserOvertime, len(overtimeModels))
	for i, model := range overtimeModels {
		overtimes[i] = model.ToOvertimeEntity()
	}

	return overtimes, nil
}
