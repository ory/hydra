// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:generate ../.bin/mockgen -package mock -destination ../internal/mock/config_cookie.go . CookieConfigProvider

package x

import (
	"context"
	"net/http"
)

type CookieConfigProvider interface {
	CookieDomain(ctx context.Context) string
	IsDevelopmentMode(ctx context.Context) bool
	CookieSameSiteMode(ctx context.Context) http.SameSite
	CookieSameSiteLegacyWorkaround(ctx context.Context) bool
	CookieSecure(ctx context.Context) bool
}
