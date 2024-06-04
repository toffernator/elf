package main

import (
	"elf/handlers"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	router := chi.NewMux()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Handle("/*", public())
	router.Get("/", handlers.Make(handlers.HandleHome))
	router.Get("/login", handlers.Make(handlers.HandleLoginIndex))
	router.Get("/ping", handlers.Make(handlers.Ping))
	router.Get("/teapot", handlers.Make(handlers.IAmATeapot))

	listenAddr := os.Getenv("LISTEN_ADDR")
	slog.Info("HTTP server started", "listenAddr", listenAddr)
	http.ListenAndServe(listenAddr, router)
}
