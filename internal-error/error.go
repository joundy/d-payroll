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
