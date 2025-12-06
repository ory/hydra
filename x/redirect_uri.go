// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"
	"net/url"

	"github.com/ory/hydra/v2/fosite"
)

type redirectConfiguration interface {
	IsDevelopmentMode(context.Context) bool
}

func IsRedirectURISecure(rc redirectConfiguration) func(context.Context, *url.URL) bool {
	return func(ctx context.Context, redirectURI *url.URL) bool {
		if rc.IsDevelopmentMode(ctx) {
			return true
		}

		if fosite.IsRedirectURISecure(ctx, redirectURI) {
			return true
		}

		return false
	}
}
