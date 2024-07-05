package rest

import (
	"fmt"
	"net/http"

	"github.com/go-playground/form"
)

type Error struct {
	StatusCode int
	Reason     any
}

func (err *Error) Error() string {
	return fmt.Sprintf("rest.Error: %v", err.Reason)
}

type DecodingError struct {
	Field       string
	Value       any
	Expectation string
}

func (err *DecodingError) Error() string {
	return fmt.Sprintf("'rest.DecodingError: '%s' should '%s'", err.Value, err.Expectation)
}

func NewValidationError(Field string, Value any, Expectation string) *DecodingError {
	return &DecodingError{Field, Value, Expectation}
}

type DecodingErrors struct {
	errs map[string]DecodingError
}

func DecodingErrorsFromDecoderErrors(errs form.DecodeErrors) *DecodingErrors {
	ds := make(map[string]DecodingError, len(errs))
	for field, err := range errs {
		ds[field] = *NewValidationError("TODO", err, err.Error())
	}
	return &DecodingErrors{ds}
}

func (err *DecodingErrors) Error() string {
	errs := make(map[string]string, len(err.errs))
	for field, e := range err.errs {
		errs[field] = e.Error()
	}

	return fmt.Sprintf("rest.DecodingErrors: %v", errs)
}

func (err *DecodingErrors) Set(field string, e DecodingError) {
	err.errs[field] = e
}

func NewAuthenticationError(reason string) *Error {
	return &Error{
		StatusCode: http.StatusUnauthorized,
		Reason:     reason,
	}
}
