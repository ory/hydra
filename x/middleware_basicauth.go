package x

import (
	"crypto/subtle"
	"net/http"
)

type BasicAuthMiddleware struct {
	basicAuthUsername string
	basicAuthPassword string
}

func NewBasicAuthMiddleware(username, password string) *BasicAuthMiddleware {
	return &BasicAuthMiddleware{
		basicAuthUsername: username,
		basicAuthPassword: password,
	}
}

func (m *BasicAuthMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	gotUsername, gotPassword, hasAuth := r.BasicAuth()
	validUsername := subtle.ConstantTimeCompare([]byte(gotUsername), []byte(m.basicAuthUsername)) == 1
	validPassword := subtle.ConstantTimeCompare([]byte(gotPassword), []byte(m.basicAuthPassword)) == 1
	if hasAuth && validUsername && validPassword {
		next(rw, r)
		return
	}
	http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}
