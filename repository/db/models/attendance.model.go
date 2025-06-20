package models

import (
	"d-payroll/entity"
	"d-payroll/utils"

	"gorm.io/gorm"
)

type AttendanceType string

const (
	AttendanceTypeCheckIn  AttendanceType = "CHECKIN"
	AttendanceTypeCheckOut AttendanceType = "CHECKOUT"
)

type UserAttendance struct {
	gorm.Model

	UserID uint
	User   *User          `gorm:"foreignKey:UserID"`
	Type   AttendanceType `gorm:"type:attendance"`
}

func (u *UserAttendance) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = utils.TimeNow()
	u.UpdatedAt = utils.TimeNow()
	return
}

func (u *UserAttendance) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = utils.TimeNow()
	return
}

func (a *UserAttendance) ToAttendanceEntity() *entity.UserAttendance {
	return &entity.UserAttendance{
		ID:        &a.ID,
		UserID:    a.UserID,
		Type:      entity.AttendanceType(a.Type),
		CreatedAt: &a.CreatedAt,
		UpdatedAt: &a.UpdatedAt,
	}
}

func (a *UserAttendance) FromAttendanceEntity(attendance *entity.UserAttendance) {
	a.UserID = attendance.UserID
	a.Type = AttendanceType(attendance.Type)

	if attendance.CreatedAt != nil {
		a.CreatedAt = *attendance.CreatedAt
	}

	if attendance.UpdatedAt != nil {
		a.UpdatedAt = *attendance.UpdatedAt
	}

}
