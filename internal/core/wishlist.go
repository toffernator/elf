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
	OwnerId int64  `validate:"required"`
	Name    string `validate:"required"`
	Image   string
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
	Take    int
	Limit   int
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
	Id    int64 `validate:"required"`
	Name  *string
	Image *string
}

func (w WishlistUpdateParams) Validate() (err error) {
	err = validate.Struct(&w)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	return ValidationErrorsFromValidatorErrors(errs)
}

type WishlistAddProductParams struct {
	Id        int64 `validate:"required"`
	ProductId int64 `validate:"required"`
}

func (w WishlistAddProductParams) Validate() (err error) {
	err = validate.Struct(&w)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	return ValidationErrorsFromValidatorErrors(errs)
}
