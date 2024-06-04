package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/sessions"
	"github.com/toffernator/elf/auth"
	"github.com/toffernator/elf/internal/core"
)

func Root(store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := ComputeAuthenticatedUser(r, store)
		if errors.Is(err, ErrNoUserInSession) {
			fmt.Fprintln(w, "Elf")
			return
		}

		tpl, err := os.ReadFile("templates/index.html")
		if err != nil {
			slog.Error("", "err", err.Error())
			http.Error(w, "We did an oopsie!", http.StatusInternalServerError)
			return
		}

		data := struct {
			IsAuthenticated bool
			Profile         auth.Profile
			Wishlists       []core.Wishlist
		}{
			IsAuthenticated: user != nil,
			Profile:         user.Profile,
			Wishlists: []core.Wishlist{
				{
					Id:       1,
					OwnerId:  1,
					Image:    "https://i1.wp.com/stpatricklincolnschool.com/wp-content/uploads/2019/01/wish-list-1.jpg?fit=1122%2C1200&ssl=1",
					Name:     "Christmas 2023",
					Products: []core.Product{},
				},
				{
					Id:       2,
					Name:     "Birthday 2021",
					Image:    "https://i1.wp.com/stpatricklincolnschool.com/wp-content/uploads/2019/01/wish-list-1.jpg?fit=1122%2C1200&ssl=1",
					OwnerId:  1,
					Products: []core.Product{},
				},
				{
					Id:       3,
					Name:     "Christmas 2020",
					Image:    "https://i1.wp.com/stpatricklincolnschool.com/wp-content/uploads/2019/01/wish-list-1.jpg?fit=1122%2C1200&ssl=1",
					OwnerId:  1,
					Products: []core.Product{},
				},
				{
					Id:       4,
					Name:     "Christmas 2019",
					Image:    "https://i1.wp.com/stpatricklincolnschool.com/wp-content/uploads/2019/01/wish-list-1.jpg?fit=1122%2C1200&ssl=1",
					OwnerId:  1,
					Products: []core.Product{},
				},
			},
		}

		t, err := template.New("webpage").Parse(string(tpl))
		t.Execute(w, data)
	}
}

func Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Pong!")
	}
}

func IAmATeapot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "I am a teapot!", http.StatusTeapot)
	}
}
