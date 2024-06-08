package handlers

import (
	"log/slog"
	"net/url"

	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
)

var validateInstance *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

var decoder *form.Decoder = form.NewDecoder()

type ApiRequest interface {
	parse() (url.Values, error)
	data() interface{}
}

func Validate(r ApiRequest) error {
	values, err := r.parse()
	if err != nil {
		return err
	}
	err = parse(r.data(), values)
	if err != nil {
		return err
	}

	err = validate(r.data())
	if err != nil {
		return err
	}

	return nil
}

// validate will transform validator.ValidateErrors into an ApiError
func validate(data interface{}) error {
	slog.Info("called with args", "data", data)

	err := validateInstance.Struct(data)
	if err != nil {
		if es, ok := err.(validator.ValidationErrors); ok {
			return ValidationErrors(es)
		}
		return err
	}

	return nil
}

func parse(dst interface{}, values url.Values) error {
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
