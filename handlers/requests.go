package handlers

import (
	"net/url"

	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

// Validate will transform validator.ValidateErrors into an ApiError
func Validate(data interface{}) error {
	err := validate.Struct(data)
	if err != nil {
		if es, ok := err.(validator.ValidationErrors); ok {
			return ValidationErrors2(es)
		}
		return err
	}

	return nil
}

func Parse(dst interface{}, values url.Values) error {
	if err := decoder.Decode(dst, values); err != nil {
		if es, ok := err.(form.DecodeErrors); ok {
			return DecoderErrors(es)
		}

		return err
	}
	return nil
}
