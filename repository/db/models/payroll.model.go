package models

import (
	"d-payroll/entity"
	"d-payroll/utils"
	"time"

	"gorm.io/gorm"
)

type Payroll struct {
	gorm.Model

	Name            string
	StartedAt       time.Time
	EndedAt         time.Time
	IsRolled        *bool `gorm:"default:false"`
	UpdatedByUserID *uint
	UpdatedByUser   *User `gorm:"foreignKey:UpdatedByUserID"`
	CreatedByUserID *uint
	CreatedByUser   *User `gorm:"foreignKey:CreatedByUserID"`
}

func (p *Payroll) BeforeCreate(tx *gorm.DB) (err error) {
	p.CreatedAt = utils.TimeNow()
	p.UpdatedAt = utils.TimeNow()
	return
}

func (p *Payroll) BeforeUpdate(tx *gorm.DB) (err error) {
	p.UpdatedAt = utils.TimeNow()
	return
}

func (p *Payroll) ToPayrollEntity() *entity.Payroll {
	return &entity.Payroll{
		ID:              &p.ID,
		Name:            p.Name,
		StartedAt:       p.StartedAt,
		EndedAt:         p.EndedAt,
		IsRolled:        p.IsRolled,
		UpdatedByUserID: p.UpdatedByUserID,
		CreatedByUserID: p.CreatedByUserID,
		CreatedAt:       &p.CreatedAt,
		UpdatedAt:       &p.UpdatedAt,
	}
}

func (p *Payroll) FromPayrollEntity(payroll *entity.Payroll) {
	p.Name = payroll.Name
	p.StartedAt = payroll.StartedAt
	p.EndedAt = payroll.EndedAt
	p.IsRolled = payroll.IsRolled
	p.UpdatedByUserID = payroll.UpdatedByUserID
	p.CreatedByUserID = payroll.CreatedByUserID

	if payroll.CreatedAt != nil {
		p.CreatedAt = *payroll.CreatedAt
	}

	if payroll.UpdatedAt != nil {
		p.UpdatedAt = *payroll.UpdatedAt
	}
}

type UserPayslipSummary struct {
	gorm.Model

	PayrollID        uint
	Payroll          *Payroll `gorm:"foreignKey:PayrollID"`
	UserID           uint
	User             *User `gorm:"foreignKey:UserID"`
	TotalTakeHomePay int
}

func (u *UserPayslipSummary) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = utils.TimeNow()
	u.UpdatedAt = utils.TimeNow()
	return
}

func (u *UserPayslipSummary) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = utils.TimeNow()
	return
}

func (u *UserPayslipSummary) ToUserPayslipSummaryEntity() *entity.UserPayslipSummary {
	return &entity.UserPayslipSummary{
		ID:               &u.ID,
		PayrollID:        u.PayrollID,
		UserID:           u.UserID,
		TotalTakeHomePay: u.TotalTakeHomePay,
		CreatedAt:        &u.CreatedAt,
		UpdatedAt:        &u.UpdatedAt,
	}
}

func (u *UserPayslipSummary) FromUserPayslipSummaryEntity(summary *entity.UserPayslipSummary) {
	u.PayrollID = summary.PayrollID
	u.UserID = summary.UserID
	u.TotalTakeHomePay = summary.TotalTakeHomePay

	if summary.CreatedAt != nil {
		u.CreatedAt = *summary.CreatedAt
	}

	if summary.UpdatedAt != nil {
		u.UpdatedAt = *summary.UpdatedAt
	}
}

func (UserPayslipSummary) TableName() string {
	return "user_payslip_summaries"
}
