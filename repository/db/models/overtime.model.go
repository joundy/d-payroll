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
	DurationMilis    int
	ApprovedByUserID *uint
	ApprovedByUser   *User `gorm:"foreignKey:ApprovedByUserID"`
}

func (u *UserOvertime) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return
}

func (u *UserOvertime) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return
}

func (o *UserOvertime) ToOvertimeEntity() *entity.UserOvertime {
	return &entity.UserOvertime{
		ID:               &o.ID,
		UserID:           o.UserID,
		Description:      o.Description,
		OvertimeAt:       o.OvertimeAt,
		DurationMilis:    o.DurationMilis,
		ApprovedByUserID: o.ApprovedByUserID,
		CreatedAt:        &o.CreatedAt,
		UpdatedAt:        &o.UpdatedAt,
	}
}

func (o *UserOvertime) FromOvertimeEntity(overtime *entity.UserOvertime) {
	o.UserID = overtime.UserID
	o.Description = overtime.Description
	o.OvertimeAt = overtime.OvertimeAt
	o.DurationMilis = overtime.DurationMilis
	o.ApprovedByUserID = overtime.ApprovedByUserID

	if overtime.CreatedAt != nil {
		o.CreatedAt = *overtime.CreatedAt
	}

	if overtime.UpdatedAt != nil {
		o.UpdatedAt = *overtime.UpdatedAt
	}
}
