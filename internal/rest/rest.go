package rest

import (
	"context"
	"elf/internal/auth/auth0"
	"elf/internal/config"
	"elf/internal/core"
	"encoding/gob"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

func init() {
	gob.Register(core.User{})
}

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

type HTTPHandler func(w http.ResponseWriter, r *http.Request) error

func MakeHandler(h HTTPHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			if e, ok := err.(*DecodingErrors); ok {
				http.Error(w, e.Error(), http.StatusBadRequest)
			} else if e, ok := err.(*core.ValidationError); ok {
				http.Error(w, e.Error(), http.StatusUnprocessableEntity)
			} else {
				http.Error(w, "We did an oopsie!", http.StatusInternalServerError)
			}
			slog.Error("HTTP handler error", "path", r.URL.Path, "err", err)
		}
	}
}

type MiddlewareFunc func(w http.ResponseWriter, r *http.Request, next http.Handler) error

func MakeMiddleware(m MiddlewareFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := m(w, r, next); err != nil {
				slog.Error("Middleware handler error", "err", err, "path", r.URL.Path)
				next.ServeHTTP(w, r)
			}
		})
	}
}

type Server struct {
	Router chi.Router
	Config *config.Config

	// TODO: Wrap these in interfaces
	SecureCookies *securecookie.SecureCookie
	Sessions      sessions.Store
	Authenticator *auth0.Authenticator

	Users     UserService
	Wishlists WishlistService
	Products  ProductService
}

func (s *Server) RegisterRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	if s.Config.IsDevelop() {
		router.Use(MakeMiddleware(func(w http.ResponseWriter, r *http.Request, next http.Handler) error {
			next.ServeHTTP(w, r)
			w.Header().Set("Cache-Control", "no-cache")
			return nil
		}))
	}
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Use(MakeMiddleware(s.AddUserToContext))

	router.Get("/", MakeHandler(s.HandleHome))
	router.Get("/ping", MakeHandler(Ping))
	router.Get("/teapot", MakeHandler(IAmATeapot))

	router.Get("/login", MakeHandler(s.Login))
	router.Get("/login/callback", MakeHandler(s.LoginCallback))
	router.Get("/logout", MakeHandler(s.Logout))
	router.Get("/logout/callback", MakeHandler(s.LogoutCallback))

	router.Route("/wishlist", func(r chi.Router) {
		// TODO: Add EnsureAuthenticated middleware
		r.Post("/", MakeHandler(s.HandleWishlistCreate))
		// TODO: Return the HTML head in a PWA style
		r.Get("/new", MakeHandler(s.HandleWishlistNew))

		r.Get("/{id}", MakeHandler(s.HandleWishlistRead))
	})

	s.Router = router
}

type ProductService interface {
	Create(ctx context.Context, p core.ProductCreateParams) (core.Product, error)
}

type UserService interface {
	Create(ctx context.Context, p core.UserCreateParams) (core.User, error)
	Read(ctx context.Context, id int64) (core.User, error)
}

type WishlistService interface {
	Create(ctx context.Context, p core.WishlistCreateParams) (core.Wishlist, error)
	Read(ctx context.Context, id int64) (core.Wishlist, error)
	ReadBy(ctx context.Context, p core.WishlistReadByParams) ([]core.Wishlist, error)
	Update(ctx context.Context, p core.WishlistUpdateParams) (core.Wishlist, error)
}

var decoder *form.Decoder = form.NewDecoder()

func Decode(v interface{}, r *http.Request) (err error) {
	err = r.ParseForm()
	if err != nil {
		return err
	}

	err = decoder.Decode(v, r.Form)
	errs, ok := err.(form.DecodeErrors)
	if !ok {
		return err
	}

	return DecodingErrorsFromDecoderErrors(errs)
}
