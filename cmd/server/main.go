package main

import (
	"context"
	"elf/internal/auth/auth0"
	"elf/internal/config"
	"elf/internal/rest"
	"elf/internal/service"
	"elf/internal/store"
	"elf/internal/store/sqlite"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	_ "github.com/glebarez/go-sqlite"
)

var (
	logger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))
)

func ConfigureSecureCookie(cfg config.SecureCookie) *securecookie.SecureCookie {
	// TODO: Look into key rotation: https://github.com/gorilla/securecookie?tab=readme-ov-file#key-rotation
	hashKey := []byte(cfg.HashKey)
	blockKey := []byte(cfg.BlockKey)

	s := securecookie.New(hashKey, blockKey)
	return s
}

func ConfigureSessionStore(cfg config.Session) (s sessions.Store, err error) {
	s = sessions.NewCookieStore([]byte(cfg.Secret))
	return s, nil
}

func ConfigureAuthenticator(cfg config.Auth0) (a *auth0.Authenticator, err error) {
	a, err = auth0.NewAuthenticator(context.Background(), cfg)
	if err != nil {
		return a, fmt.Errorf("Authenticator configuration error: %w", err)
	}
	return a, nil
}

func ConfigureDb(cfg config.Db) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite", cfg.Name)
	return db, err
}

func ConfigureWishlistStore(db *sqlx.DB) service.WishlistStore {
	return sqlite.NewWishlistStore(db)
}

func ConfigureUserStore(db *sqlx.DB) service.UserStore {
	return sqlite.NewUserStore(db)
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	cfg := config.MustLoadConfig()

	db, err := ConfigureDb(cfg.Db)
	if err != nil {
		panic(err)
	}

	s := rest.Server{
		Config: cfg,

		SecureCookies: ConfigureSecureCookie(cfg.SecureCookie),
		Sessions: func() sessions.Store {
			s, err := ConfigureSessionStore(cfg.Session)
			if err != nil {
				panic(err)
			}
			return s
		}(),
		Authenticator: func() *auth0.Authenticator {
			a, err := ConfigureAuthenticator(cfg.Auth0)
			if err != nil {
				panic(err)
			}
			return a
		}(),

		Users: func() rest.UserService {
			sqliteStore := sqlite.NewUserStore(db)
			loggedStore := store.NewLoggedUserStore(sqliteStore, logger)
			return service.NewUserService(loggedStore)
		}(),
		Wishlists: func() rest.WishlistService {
			sqliteStore := sqlite.NewWishlistStore(db)
			loggedStore := store.NewLoggedWishlistStore(sqliteStore, logger)
			return service.NewWishlistService(loggedStore)
		}(),
		Products: func() rest.ProductService {
			sqliteStore := sqlite.NewProductStore(db)
			loggedStore := store.NewLoggedProductStore(sqliteStore, logger)
			return service.NewProductService(loggedStore)
		}(),
	}

	s.RegisterRoutes()
	s.Router.Handle("/*", public())

	listenAddr := cfg.ListenAddr
	slog.Info("HTTP server started", "listenAddr", listenAddr)
	http.ListenAndServe(cfg.ListenAddr, s.Router)
	for {
	}
}
