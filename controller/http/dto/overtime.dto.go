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
	ID               *uint      `json:"id,omitempty"`
	UserID           uint       `json:"user_id"`
	Description      string     `json:"description"`
	OvertimeAt       time.Time  `json:"overtime_at"`
	DurationMilis    int        `json:"duration_milis"`
	ApprovedByUserID *uint      `json:"approved_by_user_id"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}

func (o *OvertimeResponseDto) FromOvertimeEntity(overtime *entity.UserOvertime) {
	o.ID = overtime.ID
	o.UserID = overtime.UserID
	o.Description = overtime.Description
	o.OvertimeAt = overtime.OvertimeAt
	o.DurationMilis = overtime.DurationMilis
	o.ApprovedByUserID = overtime.ApprovedByUserID
	o.CreatedAt = overtime.CreatedAt
	o.UpdatedAt = overtime.UpdatedAt
}
