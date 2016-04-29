package oauth2

import "github.com/ory-am/fosite/handler/oidc/strategy"

type Session struct {
	Subject                  string `json:"sub"`
	*strategy.IDTokenSession `json:"idToken"`
}
