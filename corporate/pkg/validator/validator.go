package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type (
	ErrorResponse struct {
		FailedField string      `json:"failed_field"`
		Tag         string      `json:"tag"`
		Value       interface{} `json:"value"`
	}

	StructValidator struct {
		validator *validator.Validate
	}
)

// NewStructValidator creates a new instance of StructValidator
func NewStructValidator() *StructValidator {
	v := validator.New()

	return &StructValidator{
		validator: v,
	}
}

func (v *StructValidator) ValidateStruct(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		var errors []ErrorResponse
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ErrorResponse{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Value:       err.Value(),
			})
		}
		return fmt.Errorf(formatValidationErrors(errors))
	}
	return nil
}

func formatValidationErrors(errors []ErrorResponse) string {
	var errMsgs []string
	for _, err := range errors {
		errMsgs = append(errMsgs, fmt.Sprintf("[%s]: '%v' | Needs to implement '%s'", err.FailedField, err.Value, err.Tag))
	}
	return strings.Join(errMsgs, " and ")
}

func (v *StructValidator) Engine() interface{} {
	return v.validator
}
