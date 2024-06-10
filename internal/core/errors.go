package core

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Error struct {
	msg any
}

func (err *Error) Error() string {
	return fmt.Sprintf("core.Error: %v", err.msg)
}

type ValidationError struct {
	Namespace   string
	Value       any
	Expectation string
}

func (err *ValidationError) Error() string {
	return fmt.Sprintf("'core.ValidationError: '%s' should '%s'", err.Value, err.Expectation)
}

func NewValidationError(Namespace string, Value any, Expectation string) *ValidationError {
	return &ValidationError{Namespace, Value, Expectation}
}

type ValidationErrors struct {
	errs []ValidationError
}

func ValidationErrorsFromValidatorErrors(errs validator.ValidationErrors) *ValidationErrors {
	vs := make([]ValidationError, len(errs))
	for i, err := range errs {
		vs[i] = *NewValidationError(err.StructField(), err.Param(), err.Error())
	}
	return &ValidationErrors{vs}
}

func (err *ValidationErrors) Error() string {
	errs := make(map[string]string, len(err.errs))
	for _, e := range err.errs {
		errs[e.Namespace] = e.Error()
	}

	return fmt.Sprintf("core.ValidationErrors: %v", errs)
}

func (err *ValidationErrors) Append(e ValidationError) {
	err.errs = append(err.errs, e)
}
