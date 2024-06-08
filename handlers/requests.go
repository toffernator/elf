package handlers

import (
	"log/slog"
	"net/url"

	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

var decoder *form.Decoder = form.NewDecoder()

// Validate will transform validator.ValidateErrors into an ApiError
func Validate(data interface{}) error {
	slog.Info("called with args", "data", data)

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
	slog.Info("called with args", "dst", dst, "values", values)

	err := decoder.Decode(dst, values)
	if err != nil {
		if es, ok := err.(form.DecodeErrors); ok {
			return DecoderErrors(es)
		}

		return err
	}

	slog.Info("ends with value", "dst", dst)
	return nil
}
