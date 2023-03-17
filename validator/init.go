package validator

import "github.com/go-playground/validator/v10"

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func validationErrorFormat(err error) []*ErrorResponse {
	var errors []*ErrorResponse
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = e.StructNamespace()
			element.Tag = e.Tag()
			element.Value = e.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
