package rest

import (
	component "elf/internal/rest/views/product"
	"net/http"
)

func (s *Server) HandleProductNew(w http.ResponseWriter, r *http.Request) (err error) {
	return Render(w, r, component.NewProduct())
}
