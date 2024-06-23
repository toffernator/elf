package sqlite

import (
	"fmt"
)

type Error struct {
	Msg    string
	Entity string
}

func (e *Error) Error() string {
	return e.Msg
}

func NewEntityDoesNotExistError(entity string, id int64) *Error {
	return &Error{
		Msg:    fmt.Sprintf("The %s with the id %d does not exist", entity, id),
		Entity: entity,
	}
}
