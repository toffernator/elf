package auth

import (
	"context"
	"errors"
	"net/url"
	"os"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

const (
	AUTH0_DOMAIN    = "AUTH0_DOMAIN"
	AUTH0_CLIENT_ID = ""
)

var (
	ErrNoIdTokenField = errors.New("no id_token field in oauth2 token")
)

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
}

// NewNewAuthenticator returns an Authenticator with a the default Auth0 OIDC
// provider, and the default Auth0 OAuth configuration for Elf.
func NewAuthenticator(ctx context.Context, auth0Issuer url.URL) (auth *Authenticator, err error) {
	provider, err := oidc.NewProvider(ctx, auth0Issuer.String())
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	auth = &Authenticator{Provider: provider, Config: conf}
	return auth, nil
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
// Returns an ErrNoIdTokenField if the *oauth2.Token is not valid.
func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, ErrNoIdTokenField
	}

	oidcConfig := &oidc.Config{
		ClientID: a.Config.ClientID,
	}

	return a.Provider.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}

func (a *Authenticator) LogoutUrl() string {

	q := url.Values{}
	q.Add("client_id", os.Getenv(AUTH0_CLIENT_ID))
	q.Add("returnTo", "http://127.0.0.1:7331/logout/callback")

	u := url.URL{
		Scheme:   "https",
		Host:     os.Getenv(AUTH0_DOMAIN),
		Path:     "/v2/logout",
		RawQuery: q.Encode(),
	}

	return u.String()
}
