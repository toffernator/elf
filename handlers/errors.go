package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
)

type ApiError struct {
	StatusCode int
	Msg        any
}

func (err ApiError) Error() string {
	return fmt.Sprintf("%d: %s", err.StatusCode, err.Msg)
}

type Field string

type Location string

const QUERY_PARAM_LOCATION Location = "URL query parameters"
const PATH_PARAM_LOCATION Location = "URL path parameters"
const FORM_LOCATION = "form"

type FieldError struct {
	Location Location
	Value    any
	Reason   string
}

func ValidationErrors(errors map[Field]FieldError) ApiError {
	return ApiError{
		StatusCode: http.StatusUnprocessableEntity,
		Msg:        errors,
	}
}

func ValidatioNErrors2(es validator.ValidationErrors) ApiError {
	return ApiError{
		StatusCode: http.StatusUnprocessableEntity,
		Msg:        es,
	}
}

func DecoderErrors(es form.DecodeErrors) ApiError {
	return ApiError{
		StatusCode: http.StatusUnprocessableEntity,
		Msg:        es,
	}
}

func Unauthenticated() ApiError {
	return ApiError{
		StatusCode: http.StatusUnauthorized,
	}
}
