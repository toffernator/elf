package core

import "github.com/go-playground/validator/v10"

type User struct {
	Id   int64  `validate:"required"`
	Name string `validate:"required"`
}

func (u *User) Validate() (err error) {
	err = validate.Struct(&u)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	return ValidationErrorsFromValidatorErrors(errs)
}

type UserCreateParams struct {
	Name string `validate:"required"`
}

func (u *UserCreateParams) Validate() (err error) {
	err = validate.Struct(&u)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	return ValidationErrorsFromValidatorErrors(errs)
}

type UserReadParams struct {
	Id int64 `validate:"required"`
}

func (u *UserReadParams) Validate() (err error) {
	err = validate.Struct(&u)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	return ValidationErrorsFromValidatorErrors(errs)
}
