package auth

import (
	"context"
	"elf/internal/config"
	"errors"
	"net/url"

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
	domain   string
}

// NewNewAuthenticator returns an Authenticator with a the default Auth0 OIDC
// provider, and the default Auth0 OAuth configuration for Elf.
func NewAuthenticator(ctx context.Context, cfg config.Auth0) (auth *Authenticator, err error) {
	auth0Issuer := url.URL{
		Scheme: "https",
		Host:   cfg.Domain,
		Path:   "/",
	}

	provider, err := oidc.NewProvider(ctx, auth0Issuer.String())
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.LoginCallbackUrl,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	auth = &Authenticator{Provider: provider, Config: conf, domain: cfg.Domain}
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

func (a *Authenticator) LogoutUrl(returnTo string) string {

	q := url.Values{}
	q.Add("client_id", a.Config.ClientID)
	q.Add("returnTo", returnTo)

	u := url.URL{
		Scheme:   "https",
		Host:     a.domain,
		Path:     "/v2/logout",
		RawQuery: q.Encode(),
	}

	return u.String()
}
