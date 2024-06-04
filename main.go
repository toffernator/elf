package main

import (
	"context"
	"elf/handlers"
	"elf/internal/auth"
	"elf/internal/handler"
	"elf/internal/store"
	"elf/middleware"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var wishlists handler.WishlistCreatorReader = &store.ArrayWishlist{}
var users auth.AuthenticatedUserStore = &auth.ArrayAuthenticatedUserStore{}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	var sessionStore = sessions.NewCookieStore([]byte(os.Getenv("GORILLA_SESSIONS_SECRET")))

	// TODO: Look into key rotation: https://github.com/gorilla/securecookie?tab=readme-ov-file#key-rotation
	var hashKey = []byte(os.Getenv("GORILLA_SECURE_COOKIES_HASHKEY"))
	var blockKey = []byte(os.Getenv("GORILLA_SECURE_COOKIES_BLOCKKEY"))
	var secureCookies = securecookie.New(hashKey, blockKey)

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{AddSource: true}))
	slog.SetDefault(logger)

	auth0Issuer := url.URL{
		Scheme: "https",
		Host:   os.Getenv(auth.AUTH0_DOMAIN),
		Path:   "/",
	}
	authenticator, err := auth.NewAuthenticator(context.Background(), auth0Issuer)
	if err != nil {
		slog.Error("The authenticator failed initialization", "err", err.Error())
		return
	}

	router := chi.NewMux()

	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)
	router.Use(middleware.Make(middleware.AddUserToContext(sessionStore, handler.AUTH0_SESSION_NAME)))

	router.Handle("/*", public())
	router.Get("/", handlers.Make(handlers.HandleHome))
	router.Get("/login", handlers.Make(handlers.Login(authenticator, secureCookies)))
	router.Get("/login/callback", handlers.Make(handlers.LoginCallback(authenticator, sessionStore, secureCookies, users)))
	router.Get("/logout", handlers.Make(handlers.Logout(authenticator)))
	router.Get("/logout/callback", handlers.Make(handlers.LogoutCallback(sessionStore)))

	router.Get("/ping", handlers.Make(handlers.Ping))
	router.Get("/teapot", handlers.Make(handlers.IAmATeapot))

	listenAddr := os.Getenv("LISTEN_ADDR")
	slog.Info("HTTP server started", "listenAddr", listenAddr)
	http.ListenAndServe(listenAddr, router)
}
