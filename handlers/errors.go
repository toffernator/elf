package handlers

import (
	"fmt"
	"net/http"
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

func ValidationError(errors map[Field]FieldError) ApiError {
	return ApiError{
		StatusCode: http.StatusUnprocessableEntity,
		Msg:        errors,
	}
}

func Unauthenticated() ApiError {
	return ApiError{
		StatusCode: http.StatusUnauthorized,
	}
}
