package dto

import (
	"d-payroll/entity"
	"time"
)

type AttendanceResponseDto struct {
	Id        *uint      `json:"id"`
	Type      string     `json:"type"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (a *AttendanceResponseDto) FromUserAttendanceEntity(attendance *entity.UserAttendance) {
	a.Id = attendance.ID
	a.Type = string(attendance.Type)
	a.CreatedAt = attendance.CreatedAt
	a.UpdatedAt = attendance.UpdatedAt
}
