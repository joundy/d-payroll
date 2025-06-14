package dto

import (
	"d-payroll/entity"
	"time"
)

type CreateOvertimeBodyDto struct {
	Description   string    `json:"description" validate:"required"`
	OvertimeAt    time.Time `json:"overtime_at" validate:"required"`
	DurationMilis int       `json:"duration_milis" validate:"required,min=1"`
}

func (c *CreateOvertimeBodyDto) ToOvertimeEntity(userID uint) *entity.UserOvertime {
	return &entity.UserOvertime{
		UserID:        userID,
		Description:   c.Description,
		OvertimeAt:    c.OvertimeAt,
		DurationMilis: c.DurationMilis,
	}
}

type OvertimeResponseDto struct {
	ID              *uint      `json:"id"`
	UserID          uint       `json:"user_id"`
	Description     string     `json:"description"`
	OvertimeAt      time.Time  `json:"overtime_at"`
	DurationMilis   int        `json:"duration_milis"`
	IsApproved      bool       `json:"is_approved"`
	UpdatedByUserID *uint      `json:"updated_by_user_id"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

func (o *OvertimeResponseDto) FromOvertimeEntity(overtime *entity.UserOvertime) {
	o.ID = overtime.ID
	o.UserID = overtime.UserID
	o.Description = overtime.Description
	o.OvertimeAt = overtime.OvertimeAt
	o.DurationMilis = overtime.DurationMilis
	o.IsApproved = overtime.IsApproved
	o.UpdatedByUserID = overtime.UpdatedByUserID
	o.CreatedAt = overtime.CreatedAt
	o.UpdatedAt = overtime.UpdatedAt
}
