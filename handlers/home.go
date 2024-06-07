package handlers

import (
	"elf/internal/core"
	"elf/views/components"
	"elf/views/home"
	"errors"
	"net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request) error {
	if openModalExists := r.URL.Query().Has("openModal"); openModalExists {
		switch r.URL.Query().Get("openModal") {
		case "newWishlist":
			return Render(w, r, components.Modal())
		default:
			return errors.New("bad request")
		}
	}

	ws := []core.Wishlist{
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
	}

	return Render(w, r, home.Index(ws))
}
