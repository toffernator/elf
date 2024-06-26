package core

import (
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
)

const (
	CurrencyEur Currency = iota
	CurrencyChf
	CurrencyNone
)

type Product struct {
	Id       int64  `validate:"required"`
	Name     string `validate:"required"`
	Url      string `validate:"url"`
	Price    int
	Currency Currency
}

func (p Product) Validate() (err error) {
	err = validate.Struct(&p)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	coreErrs := ValidationErrorsFromValidatorErrors(errs)
	currencyErr := p.Currency.validate()
	if currencyErr != nil {
		coreErrs.Append(*currencyErr)
	}

	return coreErrs
}

type Currency int16

func (c Currency) Validate() (err error) {
	return c.validate()
}

func (c Currency) validate() *ValidationError {
	switch c {
	case CurrencyEur, CurrencyChf, CurrencyNone:
		return nil
	}

	expectation := fmt.Sprintf("be in the range [%d, %d]", CurrencyEur, CurrencyNone)
	return NewValidationError("Currency", c, expectation)
}

type ProductCreateParams struct {
	BelongsToId int64  `validate:"required"`
	Name        string `validate:"required"`
	Url         string `validate:"url"`
	Price       int
	Currency    Currency
}

func (p ProductCreateParams) Validate() (err error) {
	slog.Info("ProductCreateParams.Validate is called", "p", p)
	err = validate.Struct(&p)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	return ValidationErrorsFromValidatorErrors(errs)
}
