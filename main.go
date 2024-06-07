package main

import (
	"context"
	"elf/handlers"
	"elf/internal/auth"
	"elf/internal/config"
	"elf/middleware"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var users auth.AuthenticatedUserStore = &auth.ArrayAuthenticatedUserStore{}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	cfg := config.MustLoadConfig()

	var sessionStore = sessions.NewCookieStore([]byte(os.Getenv("GORILLA_SESSIONS_SECRET")))

	// TODO: Look into key rotation: https://github.com/gorilla/securecookie?tab=readme-ov-file#key-rotation
	var hashKey = []byte(os.Getenv("GORILLA_SECURE_COOKIES_HASHKEY"))
	var blockKey = []byte(os.Getenv("GORILLA_SECURE_COOKIES_BLOCKKEY"))
	var secureCookies = securecookie.New(hashKey, blockKey)

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{AddSource: true}))
	slog.SetDefault(logger)

	authenticator, err := auth.NewAuthenticator(context.Background(), cfg.Auth0)
	if err != nil {
		slog.Error("The authenticator failed initialization", "err", err.Error())
		return
	}

	router := chi.NewMux()

	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)
	router.Use(middleware.Make(middleware.AddUserToContext(sessionStore, cfg.Auth.SessionCookieName, cfg.Auth.SessionCookieUserKey)))

	router.Handle("/*", public())
	router.Get("/home", handlers.Make(handlers.HandleHome))
	router.Get("/login", handlers.Make(handlers.Login(authenticator, secureCookies, cfg.OAuth.StateLength, cfg.OAuth.StateCookieName)))
	router.Get("/login/callback", handlers.Make(handlers.LoginCallback(authenticator, sessionStore, secureCookies, users, cfg.OAuth.StateCookieName, cfg.Auth.SessionCookieName, cfg.Auth.SessionCookieUserKey, cfg.Auth0.SessionCookieAccessTokenKey)))
	router.Get("/logout", handlers.Make(handlers.Logout(authenticator, cfg.Auth0.LogoutCallbackUrl)))
	router.Get("/logout/callback", handlers.Make(handlers.LogoutCallback(sessionStore, cfg.Auth.SessionCookieName, cfg.Auth.SessionCookieUserKey)))

	router.Post("/wishlist", handlers.Make(handlers.Ping))

	router.Get("/ping", handlers.Make(handlers.Ping))
	router.Get("/teapot", handlers.Make(handlers.IAmATeapot))

	listenAddr := cfg.ListenAddr
	slog.Info("HTTP server started", "listenAddr", listenAddr)
	http.ListenAndServe(listenAddr, router)
}
