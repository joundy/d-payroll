package entity

import "time"

type AttendanceType string

const (
	AttendanceTypeCheckIn  AttendanceType = "CHECKIN"
	AttendanceTypeCheckOut AttendanceType = "CHECKOUT"
)

type UserAttendance struct {
	ID        *uint
	UserID    uint
	Type      AttendanceType
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
