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

func ValidationErrors(es validator.ValidationErrors) ApiError {
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
