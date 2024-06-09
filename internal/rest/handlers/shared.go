package handlers

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

type HTTPHandler func(w http.ResponseWriter, r *http.Request) error

func Make(h HTTPHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if e, ok := err.(ApiError); ok {
				http.Error(w, e.Error(), e.StatusCode)
			} else {
				http.Error(w, "We did an oopsie!", http.StatusInternalServerError)
			}

			slog.Error("HTTP handler error", "err", err, "path", r.URL.Path)
		}
	}
}

func Render(w http.ResponseWriter, r *http.Request, c templ.Component) error {
	return c.Render(r.Context(), w)
}