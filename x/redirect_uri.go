package x

import (
	"context"
	"net/url"

	"github.com/ory/fosite"
)

type redirectConfiguration interface {
	InsecureRedirects() []string
}

func IsRedirectURISecure(rc redirectConfiguration) func(ctx context.Context, redirectURI *url.URL) bool {
	return func(ctx context.Context, redirectURI *url.URL) bool {
		if fosite.IsRedirectURISecure(nil, redirectURI) {
			return true
		}

		for _, allowed := range rc.InsecureRedirects() {
			if redirectURI.String() == allowed {
				return true
			}
		}

		return false
	}
}
