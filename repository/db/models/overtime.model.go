package models

import (
	"d-payroll/entity"
	"time"

	"gorm.io/gorm"
)

type UserOvertime struct {
	gorm.Model

	UserID           uint
	User             *User `gorm:"foreignKey:UserID"`
	Description      string
	OvertimeAt       time.Time
	DurationMinutes  int
	ApprovedByUserID *uint
	ApprovedByUser   *User `gorm:"foreignKey:ApprovedByUserID"`
}

func (o *UserOvertime) ToOvertimeEntity() *entity.UserOvertime {
	return &entity.UserOvertime{
		ID:               &o.ID,
		UserID:           o.UserID,
		Description:      o.Description,
		OvertimeAt:       o.OvertimeAt,
		DurationMinutes:  o.DurationMinutes,
		ApprovedByUserID: o.ApprovedByUserID,
		CreatedAt:        &o.CreatedAt,
		UpdatedAt:        &o.UpdatedAt,
	}
}

func (o *UserOvertime) FromOvertimeEntity(overtime *entity.UserOvertime) {
	o.UserID = overtime.UserID
	o.Description = overtime.Description
	o.OvertimeAt = overtime.OvertimeAt
	o.DurationMinutes = overtime.DurationMinutes
	o.ApprovedByUserID = overtime.ApprovedByUserID

	if overtime.CreatedAt != nil {
		o.CreatedAt = *overtime.CreatedAt
	}

	if overtime.UpdatedAt != nil {
		o.UpdatedAt = *overtime.UpdatedAt
	}
}
