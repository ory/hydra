package x

import (
	"context"
	"net/url"

	"github.com/ory/fosite"
)

type redirectConfiguration interface {
	InsecureRedirects(context.Context) []string
}

func IsRedirectURISecure(rc redirectConfiguration) func(context.Context, *url.URL) bool {
	return func(ctx context.Context, redirectURI *url.URL) bool {
		if fosite.IsRedirectURISecure(ctx, redirectURI) {
			return true
		}

		for _, allowed := range rc.InsecureRedirects(ctx) {
			if redirectURI.String() == allowed {
				return true
			}
		}

		return false
	}
}
