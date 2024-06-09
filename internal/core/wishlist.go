package core

import "github.com/go-playground/validator/v10"

type Wishlist struct {
	Id       int64  `validate:"required"`
	OwnerId  int64  `validate:"required"`
	Name     string `validate:"required"`
	Image    string
	Products []Product
}

func (w Wishlist) Validate() (err error) {
	err = validate.Struct(&w)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	return ValidationErrorsFromValidatorErrors(errs)
}

type WishlistCreateParams struct {
	OwnerId  int64  `validate:"required"`
	Name     string `validate:"required"`
	Image    string
	Products []ProductCreateParams
}

func (w WishlistCreateParams) Validate() (err error) {
	err = validate.Struct(&w)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	return ValidationErrorsFromValidatorErrors(errs)
}

type WishlistReadByParams struct {
	OwnerId int64 `validate:"required"`
}

func (w WishlistReadByParams) Validate() (err error) {
	err = validate.Struct(&w)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	return ValidationErrorsFromValidatorErrors(errs)
}

type WishlistUpdateParams struct {
	Id       int64 `validate:"required"`
	OwnerId  *int64
	Name     *string
	Image    *string
	Products []ProductCreateParams
}

func (w WishlistUpdateParams) Validate() (err error) {
	err = validate.Struct(&w)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	return ValidationErrorsFromValidatorErrors(errs)
}
