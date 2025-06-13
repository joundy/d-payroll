package internalerror

import (
	"d-payroll/entity"
	"fmt"
)

type ValidationError struct {
	Fields []entity.ValidationErrorField `json:"fields"`
}

func (v *ValidationError) toString() string {
	var str string

	for _, field := range v.Fields {
		str += fmt.Sprintf("\n%s: %s", field.Field, field.Tag)
	}

	return str
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("Validation error: %s", v.toString())
}

type NotFoundError struct{}

func (d *NotFoundError) Error() string {
	return "Data not found"
}

type InvalidCredentialsError struct{}

func (i *InvalidCredentialsError) Error() string {
	return "Invalid credentials"
}

type AttendanceWeekendError struct{}

func (a *AttendanceWeekendError) Error() string {
	return "Attendance cannot checked in on weekend"
}

type AttendanceAlreadyCheckedInError struct{}

func (a *AttendanceAlreadyCheckedInError) Error() string {
	return "Attendance already checked in"
}

type AttendanceAlreadyCheckedOutError struct{}

func (a *AttendanceAlreadyCheckedOutError) Error() string {
	return "Attendance already checked out"
}

type AttendanceCannotCheckedOutError struct{}

func (a *AttendanceCannotCheckedOutError) Error() string {
	return "Attendance cannot checked out because it is not checked in"
}

type ReimbursementAlreadyApprovedError struct{}

func (r *ReimbursementAlreadyApprovedError) Error() string {
	return "Reimbursement already approved"
}
