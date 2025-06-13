package dto

import (
	"d-payroll/entity"
	"time"
)

type CreateOvertimeBodyDto struct {
	Description     string    `json:"description" validate:"required"`
	OvertimeAt      time.Time `json:"overtime_at" validate:"required"`
	DurationMinutes int       `json:"duration_minutes" validate:"required,min=1"`
}

func (c *CreateOvertimeBodyDto) ToOvertimeEntity(userID uint) *entity.UserOvertime {
	return &entity.UserOvertime{
		UserID:          userID,
		Description:     c.Description,
		OvertimeAt:      c.OvertimeAt,
		DurationMinutes: c.DurationMinutes,
	}
}

type OvertimeResponseDto struct {
	ID               *uint      `json:"id,omitempty"`
	UserID           uint       `json:"user_id"`
	Description      string     `json:"description"`
	OvertimeAt       time.Time  `json:"overtime_at"`
	DurationMinutes  int        `json:"duration_minutes"`
	ApprovedByUserID *uint      `json:"approved_by_user_id"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}

func (o *OvertimeResponseDto) FromOvertimeEntity(overtime *entity.UserOvertime) {
	o.ID = overtime.ID
	o.UserID = overtime.UserID
	o.Description = overtime.Description
	o.OvertimeAt = overtime.OvertimeAt
	o.DurationMinutes = overtime.DurationMinutes
	o.ApprovedByUserID = overtime.ApprovedByUserID
	o.CreatedAt = overtime.CreatedAt
	o.UpdatedAt = overtime.UpdatedAt
}
