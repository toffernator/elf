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

func InvalidQueryParameters(errors map[string]string) ApiError {
	return ApiError{
		StatusCode: http.StatusBadRequest,
		Msg:        errors,
	}
}

func Unauthenticated() ApiError {
	return ApiError{
		StatusCode: http.StatusUnauthorized,
	}
}
