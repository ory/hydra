package x

import (
	"context"
	"github.com/ory/fosite"
	"net/http"
	"net/url"
)

type ClientAuthenticatorProvider interface {
	ClientAuthenticator() ClientAuthenticator
}

type ClientAuthenticator interface {
	AuthenticateClient(ctx context.Context, r *http.Request, form url.Values) (fosite.Client, error)
}
