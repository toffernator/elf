//go:build dev

package main

import (
	"log/slog"
	"net/http"
	"os"
)

func public() http.Handler {
	slog.Info("building static files for development")
	return http.StripPrefix("/public/", http.FileServerFS(os.DirFS("cmd/server/public")))
}
