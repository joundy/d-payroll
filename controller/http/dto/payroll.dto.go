package dto

import (
	"d-payroll/entity"
	"time"
)

type CreatePayrollBodyDto struct {
	Name      string    `json:"name" validate:"required"`
	StartedAt time.Time `json:"started_at" validate:"required"`
	EndedAt   time.Time `json:"ended_at" validate:"required"`
}

func (c *CreatePayrollBodyDto) ToPayrollEntity(userID uint) *entity.Payroll {
	return &entity.Payroll{
		Name:            c.Name,
		StartedAt:       c.StartedAt,
		EndedAt:         c.EndedAt,
		CreatedByUserID: &userID,
	}
}

type PayrollResponseDto struct {
	ID              *uint      `json:"id"`
	Name            string     `json:"name"`
	StartedAt       time.Time  `json:"started_at"`
	EndedAt         time.Time  `json:"ended_at"`
	IsRolled        *bool      `json:"is_rolled"`
	UpdatedByUserID *uint      `json:"updated_by_user_id"`
	CreatedByUserID *uint      `json:"created_by_user_id"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

func (p *PayrollResponseDto) FromPayrollEntity(payroll *entity.Payroll) {
	p.ID = payroll.ID
	p.Name = payroll.Name
	p.StartedAt = payroll.StartedAt
	p.EndedAt = payroll.EndedAt
	p.IsRolled = payroll.IsRolled
	p.UpdatedByUserID = payroll.UpdatedByUserID
	p.CreatedByUserID = payroll.CreatedByUserID
	p.CreatedAt = payroll.CreatedAt
	p.UpdatedAt = payroll.UpdatedAt
}
