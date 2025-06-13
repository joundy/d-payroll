package dto

import (
	"d-payroll/entity"
	"time"
)

type AttendanceResponseDto struct {
	Id        *uint      `json:"id,omitempty"`
	Type      string     `json:"type"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (a *AttendanceResponseDto) FromUserAttendanceEntity(attendance *entity.UserAttendance) {
	a.Id = attendance.ID
	a.Type = string(attendance.Type)
	a.CreatedAt = attendance.CreatedAt
	a.UpdatedAt = attendance.UpdatedAt
}
