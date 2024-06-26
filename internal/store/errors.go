package store

import (
	"fmt"
)

type EntityDoesNotExistError struct {
	Entity string
	Id     int64
}

func (e *EntityDoesNotExistError) Error() string {
	return fmt.Sprintf("The %s with the id %d does not exist", e.Entity, e.Id)
}

func (e EntityDoesNotExistError) Is(target error) bool {
	switch err := target.(type) {
	case *EntityDoesNotExistError:
		return err.Entity == e.Entity && err.Id == e.Id
	default:
		return false
	}
}

func NewEntityDoesNotExistError(entity string, id int64) *EntityDoesNotExistError {
	return &EntityDoesNotExistError{
		Entity: entity,
		Id:     id,
	}
}
