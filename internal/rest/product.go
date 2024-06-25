package rest

import (
	component "elf/internal/rest/views/product"
	"net/http"
	"strconv"
)

func (s *Server) HandleProductNew(w http.ResponseWriter, r *http.Request) (err error) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return &DecodingError{Field: "id", Value: idStr, Expectation: "be an integer"}
	}

	return Render(w, r, component.NewProduct(id))
}
