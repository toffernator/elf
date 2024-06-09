package handlers

import (
	"errors"
	"fmt"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "Pong!")
	return nil
}

func IAmATeapot(w http.ResponseWriter, r *http.Request) error {
	return errors.New("I'm a teapot!")
}
