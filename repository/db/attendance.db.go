package repository

import (
	"context"
	internalerror "d-payroll/internal-error"
	"d-payroll/repository/db/models"
	"d-payroll/utils"
	"errors"

	"gorm.io/gorm"
)

type AttendanceDB interface {
	CreateAttendance(ctx context.Context, attendance *models.UserAttendance) error
	GetThisDayAttendanceByUserID(ctx context.Context, userID uint, attenanceType models.AttendanceType) (*models.UserAttendance, error)
	GetAttendancesByUserID(ctx context.Context, userID uint) ([]*models.UserAttendance, error)
}

type attendanceDB struct {
	DB *gorm.DB
}

func NewAttendanceDB(db *gorm.DB) AttendanceDB {
	return &attendanceDB{DB: db}
}

func (e *attendanceDB) CreateAttendance(ctx context.Context, attendance *models.UserAttendance) error {
	return e.DB.WithContext(ctx).Create(attendance).Error
}

func (e *attendanceDB) GetThisDayAttendanceByUserID(ctx context.Context, userID uint, attendanceType models.AttendanceType) (*models.UserAttendance, error) {
	var attendance *models.UserAttendance
	result := e.DB.WithContext(ctx).Where("user_id = ? AND type = ? AND created_at BETWEEN ? AND ?", userID, attendanceType, utils.GetStartOfDay(), utils.GetEndOfDay()).First(&attendance)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &internalerror.NotFoundError{}
		}

		return nil, result.Error
	}

	return attendance, nil
}

func (e *attendanceDB) GetAttendancesByUserID(ctx context.Context, userID uint) ([]*models.UserAttendance, error) {
	var attendances []*models.UserAttendance
	result := e.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&attendances)
	if result.Error != nil {
		return nil, result.Error
	}
	return attendances, nil
}
