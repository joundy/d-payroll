package utils

import (
	"d-payroll/entity"
	internalerror "d-payroll/internal-error"

	"github.com/go-playground/validator/v10"
)

var v = validator.New()

func ValidateStruct(data interface{}) error {
	validationErrors := []entity.ValidationErrorField{}

	errs := v.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem entity.ValidationErrorField

			elem.Field = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()

			validationErrors = append(validationErrors, elem)
		}
	}

	if len(validationErrors) > 0 {
		return &internalerror.ValidationError{
			Fields: validationErrors,
		}
	}

	return nil
}
