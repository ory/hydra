// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/mapx"
)

func setCSRFCookie(ctx context.Context, w http.ResponseWriter, r *http.Request, conf x.CookieConfigProvider, store sessions.Store, name, csrfValue string, maxAge time.Duration) error {
	// Errors can be ignored here, because we always get a session back. Error typically means that the
	// session doesn't exist yet.
	session, _ := store.Get(r, name)

	sameSite := conf.CookieSameSiteMode(ctx)
	if isLegacyCSRFCookieName(name) {
		sameSite = 0
	}

	session.Values["csrf"] = csrfValue
	session.Options.HttpOnly = true
	session.Options.Secure = conf.CookieSecure(ctx)
	session.Options.SameSite = sameSite
	session.Options.Domain = conf.CookieDomain(ctx)
	session.Options.MaxAge = int(maxAge.Seconds())
	if err := session.Save(r, w); err != nil {
		return errors.WithStack(err)
	}

	if sameSite == http.SameSiteNoneMode && conf.CookieSameSiteLegacyWorkaround(ctx) {
		return setCSRFCookie(ctx, w, r, conf, store, legacyCSRFCookieName(name), csrfValue, maxAge)
	}

	return nil
}

func validateCSRFCookie(ctx context.Context, r *http.Request, conf x.CookieConfigProvider, store sessions.Store, name, expectedCSRF string) error {
	if cookie, err := getCSRFCookie(ctx, r, store, conf, name); err != nil {
		return errors.WithStack(fosite.ErrRequestForbidden.WithHint("CSRF session cookie could not be decoded."))
	} else if csrf, err := mapx.GetString(cookie.Values, "csrf"); err != nil {
		return errors.WithStack(fosite.ErrRequestForbidden.WithHint("No CSRF value available in the session cookie."))
	} else if csrf != expectedCSRF {
		return errors.WithStack(fosite.ErrRequestForbidden.WithHint("The CSRF value from the token does not match the CSRF value from the data store."))
	}

	return nil
}

func getCSRFCookie(ctx context.Context, r *http.Request, store sessions.Store, conf x.CookieConfigProvider, name string) (*sessions.Session, error) {
	cookie, err := store.Get(r, name)
	if !isLegacyCSRFCookieName(name) &&
		conf.CookieSameSiteMode(ctx) == http.SameSiteNoneMode &&
		conf.CookieSameSiteLegacyWorkaround(ctx) &&
		(err != nil || len(cookie.Values) == 0) {
		return store.Get(r, legacyCSRFCookieName(name))
	}
	return cookie, err
}

func legacyCSRFCookieName(name string) string { return name + "_legacy" }
func isLegacyCSRFCookieName(name string) bool { return strings.HasSuffix(name, "_legacy") }
