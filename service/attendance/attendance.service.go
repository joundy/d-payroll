package attendanceservice

import (
	"context"
	"d-payroll/entity"
	internalerror "d-payroll/internal-error"
	repository "d-payroll/repository/db"
	"d-payroll/repository/db/models"
	"d-payroll/utils"
	"errors"
)

// TODO:
// - Possible race condition, checkin and checkout at the same time
// - Fix using transaction or mutex lock

type AttendanceService interface {
	Checkin(ctx context.Context, userID uint) (*entity.UserAttendance, error)
	Checkout(ctx context.Context, userID uint) (*entity.UserAttendance, error)
	IsCheckedOut(ctx context.Context, userID uint) (bool, error)

	GetAttendancesByUserID(ctx context.Context, userID uint) ([]*entity.UserAttendance, error)
}

type attendanceService struct {
	attendanceDB repository.AttendanceDB
}

func NewAttendanceService(attendanceDB repository.AttendanceDB) AttendanceService {
	return &attendanceService{attendanceDB: attendanceDB}
}

func (s *attendanceService) Checkin(ctx context.Context, userID uint) (*entity.UserAttendance, error) {
	if utils.IsWeekend() {
		return nil, &internalerror.AttendanceWeekendError{}
	}

	attendanceModel := &models.UserAttendance{
		UserID: userID,
		Type:   models.AttendanceTypeCheckIn,
	}

	thisDayCheckin, err := s.attendanceDB.GetThisDayAttendanceByUserID(ctx, userID, models.AttendanceTypeCheckIn)
	if err != nil {
		if !errors.Is(err, &internalerror.NotFoundError{}) {
			return nil, err
		}
	}
	if thisDayCheckin != nil {
		return nil, &internalerror.AttendanceAlreadyCheckedInError{}
	}

	err = s.attendanceDB.CreateAttendance(ctx, attendanceModel)
	if err != nil {
		return nil, err
	}

	return attendanceModel.ToAttendanceEntity(), nil
}

func (s *attendanceService) Checkout(ctx context.Context, userID uint) (*entity.UserAttendance, error) {
	_, err := s.attendanceDB.GetThisDayAttendanceByUserID(ctx, userID, models.AttendanceTypeCheckIn)
	if err != nil {
		if errors.Is(err, &internalerror.NotFoundError{}) {
			return nil, &internalerror.AttendanceCannotCheckedOutError{}
		}

		return nil, err
	}

	thisDayCheckout, err := s.attendanceDB.GetThisDayAttendanceByUserID(ctx, userID, models.AttendanceTypeCheckOut)
	if err != nil {
		if !errors.Is(err, &internalerror.NotFoundError{}) {
			return nil, err
		}
	}
	if thisDayCheckout != nil {
		return nil, &internalerror.AttendanceAlreadyCheckedOutError{}
	}

	attendanceModel := &models.UserAttendance{
		UserID: userID,
		Type:   models.AttendanceTypeCheckOut,
	}

	err = s.attendanceDB.CreateAttendance(ctx, attendanceModel)
	if err != nil {
		return nil, err
	}

	return attendanceModel.ToAttendanceEntity(), nil
}

func (s *attendanceService) IsCheckedOut(ctx context.Context, userID uint) (bool, error) {
	thisDayCheckout, err := s.attendanceDB.GetThisDayAttendanceByUserID(ctx, userID, models.AttendanceTypeCheckOut)
	if err != nil {
		if !errors.Is(err, &internalerror.NotFoundError{}) {
			return false, err
		}
	}
	return thisDayCheckout != nil, nil
}

func (s *attendanceService) GetAttendancesByUserID(ctx context.Context, userID uint) ([]*entity.UserAttendance, error) {
	attendances, err := s.attendanceDB.GetAttendancesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var userAttendances []*entity.UserAttendance
	for _, attendance := range attendances {
		userAttendances = append(userAttendances, attendance.ToAttendanceEntity())
	}

	return userAttendances, nil
}
