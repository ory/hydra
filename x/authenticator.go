package x

import (
	"context"
	"net/http"
	"net/url"

	"github.com/ory/fosite"
)

type ClientAuthenticatorProvider interface {
	ClientAuthenticator() ClientAuthenticator
}

type ClientAuthenticator interface {
	AuthenticateClient(ctx context.Context, r *http.Request, form url.Values) (fosite.Client, error)
}
