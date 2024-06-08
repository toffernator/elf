package main

import (
	"context"
	"elf/handlers"
	"elf/internal/auth"
	"elf/internal/config"
	"elf/internal/store"
	"elf/middleware"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	_ "github.com/glebarez/go-sqlite"
)

type Services struct {
	SecureCookie  *securecookie.SecureCookie
	Sessions      sessions.Store
	Authenticator *auth.Authenticator

	Wishlists *store.Wishlist
	Users     *store.User

	Db *sqlx.DB
}

func ConfigureServices(cfg config.Config) (*Services, error) {
	ConfigureLogger(cfg)

	secureCookie := ConfigureSecureCookie(cfg.SecureCookie)

	sessions, err := ConfigureSessionStore(cfg.Session)
	if err != nil {
		return nil, err
	}

	authenticator, err := ConfigureAuthenticator(cfg.Auth0)
	if err != nil {
		return nil, err
	}

	db, err := ConfigureDb(cfg.Db)
	if err != nil {
		return nil, err
	}

	wishlists := ConfigureWishlistStore(db)
	users := ConfigureUserStore(db)

	return &Services{
		SecureCookie:  secureCookie,
		Sessions:      sessions,
		Authenticator: authenticator,

		Wishlists: wishlists,
		Users:     users,

		Db: db,
	}, nil
}

func ConfigureLogger(cfg config.Config) {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)
}

func ConfigureSecureCookie(cfg config.SecureCookie) *securecookie.SecureCookie {
	// TODO: Look into key rotation: https://github.com/gorilla/securecookie?tab=readme-ov-file#key-rotation
	hashKey := []byte(cfg.HashKey)
	blockKey := []byte(cfg.BlockKey)

	s := securecookie.New(hashKey, blockKey)
	return s
}

func ConfigureSessionStore(cfg config.Session) (s sessions.Store, err error) {
	s = sessions.NewCookieStore([]byte(cfg.Secret))
	if s == nil {
		return s, errors.New("Session store configuration error")
	}
	return s, nil
}

func ConfigureAuthenticator(cfg config.Auth0) (a *auth.Authenticator, err error) {
	a, err = auth.NewAuthenticator(context.Background(), cfg)
	if err != nil {
		return a, fmt.Errorf("Authenticator configuration error: %w", err)
	}
	return a, nil
}

func ConfigureDb(cfg config.Db) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite", cfg.Name)
	return db, err
}

func ConfigureWishlistStore(db *sqlx.DB) *store.Wishlist {
	return store.NewWishlist(db)
}

func ConfigureUserStore(db *sqlx.DB) *store.User {
	return store.NewUser(db)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	cfg := config.MustLoadConfig()

	services, err := ConfigureServices(*cfg)
	if err != nil {
		panic(err)
	}

	router := chi.NewMux()

	middlewareAuthServices := &middleware.AuthServices{
		Sessions: services.Sessions,
	}
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)
	router.Use(middleware.Make(middleware.AddUserToContext(cfg, middlewareAuthServices)))

	router.Get("/", handlers.Make(handlers.Index))
	router.Handle("/*", public())

	homeServices := &handlers.HomeServices{
		Wishlists: services.Wishlists,
	}
	router.Get("/home", handlers.Make(handlers.Home(cfg, homeServices)))

	handlerAuthServices := &handlers.AuthServices{
		SecureCookie:  services.SecureCookie,
		Sessions:      services.Sessions,
		Authenticator: services.Authenticator,
		Users:         services.Users,
	}
	router.Get("/login", handlers.Make(handlers.Login(cfg, handlerAuthServices)))
	router.Get("/login/callback", handlers.Make(handlers.LoginCallback(cfg, handlerAuthServices)))
	router.Get("/logout", handlers.Make(handlers.Logout(cfg, handlerAuthServices)))
	router.Get("/logout/callback", handlers.Make(handlers.LogoutCallback(cfg, handlerAuthServices)))

	handlerWishlistServices := &handlers.WishlistServices{
		WishlistCreator: services.Wishlists,
		WishlistReader:  services.Wishlists,
	}
	router.Post("/wishlist", handlers.Make(handlers.NewWishlist(cfg, handlerWishlistServices)))
	router.Get("/wishlist/{id}", handlers.Make(handlers.GetWishlist(cfg, handlerWishlistServices)))
	router.Get("/wishlist/{id}/page", handlers.Make(handlers.GetWishlistPage(cfg, handlerWishlistServices)))

	router.Get("/ping", handlers.Make(handlers.Ping))
	router.Get("/teapot", handlers.Make(handlers.IAmATeapot))

	listenAddr := cfg.ListenAddr
	slog.Info("HTTP server started", "listenAddr", listenAddr)
	http.ListenAndServe(cfg.ListenAddr, router)
}
