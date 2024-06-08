package middleware

import "github.com/gorilla/sessions"

type Params struct {
	Sessions             sessions.Store
	SessionCookieName    string
	SessionCookieUserKey string
}
