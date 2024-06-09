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

type WishlistReadParams struct {
	Id      *int64
	OwnerId *int64
}

func (w *WishlistReadParams) IsZero() bool {
	return w.Id == nil && w.OwnerId == nil
}

func (w WishlistReadParams) Validate() (err error) {
	err = validate.Struct(&w)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	return ValidationErrorsFromValidatorErrors(errs)
}

type WishlistUpdateParams struct {
	OwnerId  *int64
	Name     *string
	Image    *string
	Products *[]ProductCreateParams
}

func (w WishlistUpdateParams) Validate() (err error) {
	err = validate.Struct(&w)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	return ValidationErrorsFromValidatorErrors(errs)
}
