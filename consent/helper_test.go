// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/internal/mock"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/client"
)

func TestSanitizeClient(t *testing.T) {
	c := &client.Client{
		Secret: "some-secret",
	}
	ar := &fosite.AuthorizeRequest{
		Request: fosite.Request{
			Client: c,
		},
	}
	got := sanitizeClientFromRequest(ar)
	assert.Empty(t, got.Secret)
	assert.NotEmpty(t, c.Secret)
}

func TestMatchScopes(t *testing.T) {
	for k, tc := range []struct {
		granted, requested []string
		expected           bool
	}{{
		granted:   []string{"foo", "bar"},
		requested: []string{"foo", "bar"},
		expected:  true,
	}, {
		granted:   []string{"foo", "bar"},
		requested: []string{"foo", "bar", "baz"},
		expected:  false,
	}, {
		granted:   []string{"foo", "bar"},
		requested: []string{"foo"},
		expected:  true,
	}, {
		granted:   []string{"foo", "bar"},
		requested: []string{"zab", "baz"},
		expected:  false,
	}} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.Equal(t, tc.expected, matchScopes(fosite.ExactScopeStrategy, tc.granted, tc.requested))
		})
	}
}

func TestValidateCsrfSession(t *testing.T) {
	const name = "oauth2_authentication_csrf"

	type cookie struct {
		name      string
		csrfValue string
		sameSite  http.SameSite
	}
	for k, tc := range []struct {
		cookies                  []cookie
		csrfValue                string
		sameSiteLegacyWorkaround bool
		expectError              bool
		sameSite                 http.SameSite
	}{
		{
			cookies:                  []cookie{},
			csrfValue:                "CSRF-VALUE",
			sameSiteLegacyWorkaround: false,
			expectError:              true,
		},
		{
			cookies:                  []cookie{},
			csrfValue:                "CSRF-VALUE",
			sameSiteLegacyWorkaround: true,
			expectError:              true,
		},
		{
			cookies: []cookie{
				{
					name:      name,
					csrfValue: "WRONG-CSRF-VALUE",
					sameSite:  http.SameSiteDefaultMode,
				},
			},
			csrfValue:                "CSRF-VALUE",
			sameSiteLegacyWorkaround: false,
			expectError:              true,
		},
		{
			cookies: []cookie{
				{
					name:      name,
					csrfValue: "WRONG-CSRF-VALUE",
					sameSite:  http.SameSiteDefaultMode,
				},
			},
			csrfValue:                "CSRF-VALUE",
			sameSiteLegacyWorkaround: true,
			expectError:              true,
		},
		{
			cookies: []cookie{
				{
					name:      name,
					csrfValue: "CSRF-VALUE",
					sameSite:  http.SameSiteDefaultMode,
				},
			},
			csrfValue:                "CSRF-VALUE",
			sameSiteLegacyWorkaround: false,
			expectError:              false,
		},
		{
			cookies: []cookie{
				{
					name:      name,
					csrfValue: "CSRF-VALUE",
					sameSite:  http.SameSiteDefaultMode,
				},
			},
			csrfValue:                "CSRF-VALUE",
			sameSiteLegacyWorkaround: true,
			expectError:              false,
		},
		{
			cookies: []cookie{
				{
					name:      legacyCSRFSessionName(name),
					csrfValue: "CSRF-VALUE",
					sameSite:  http.SameSiteDefaultMode,
				},
			},
			csrfValue:                "CSRF-VALUE",
			sameSiteLegacyWorkaround: false,
			expectError:              true,
		},
		{
			cookies: []cookie{
				{
					name:      legacyCSRFSessionName(name),
					csrfValue: "CSRF-VALUE",
					sameSite:  http.SameSiteDefaultMode,
				},
			},
			csrfValue:                "CSRF-VALUE",
			sameSiteLegacyWorkaround: true,
			expectError:              false,
			sameSite:                 http.SameSiteNoneMode,
		},
		{
			cookies: []cookie{
				{
					name:      legacyCSRFSessionName(name),
					csrfValue: "CSRF-VALUE",
					sameSite:  http.SameSiteDefaultMode,
				},
			},
			csrfValue:                "CSRF-VALUE",
			sameSiteLegacyWorkaround: true,
			expectError:              true,
			sameSite:                 http.SameSiteLaxMode,
		},
		{
			cookies: []cookie{
				{
					name:      name,
					csrfValue: "CSRF-VALUE",
					sameSite:  http.SameSiteNoneMode,
				},
				{
					name:      legacyCSRFSessionName(name),
					csrfValue: "CSRF-VALUE",
					sameSite:  http.SameSiteDefaultMode,
				},
			},
			csrfValue:                "CSRF-VALUE",
			sameSiteLegacyWorkaround: false,
			expectError:              false,
		},
		{
			cookies: []cookie{
				{
					name:      name,
					csrfValue: "CSRF-VALUE",
					sameSite:  http.SameSiteNoneMode,
				},
				{
					name:      legacyCSRFSessionName(name),
					csrfValue: "CSRF-VALUE",
					sameSite:  http.SameSiteDefaultMode,
				},
			},
			csrfValue:                "CSRF-VALUE",
			sameSiteLegacyWorkaround: true,
			expectError:              false,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			config := mock.NewMockCookieConfigProvider(ctrl)
			config.EXPECT().CookieSameSiteLegacyWorkaround(gomock.Any()).Return(tc.sameSiteLegacyWorkaround).AnyTimes()
			config.EXPECT().IsDevelopmentMode(gomock.Any()).Return(false).AnyTimes()
			config.EXPECT().CookieSecure(gomock.Any()).Return(false).AnyTimes()

			sameSite := http.SameSiteDefaultMode
			if tc.sameSite > 0 {
				sameSite = tc.sameSite
			}

			config.EXPECT().CookieSameSiteMode(gomock.Any()).Return(sameSite).AnyTimes()

			store := sessions.NewCookieStore(securecookie.GenerateRandomKey(16))
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			for _, c := range tc.cookies {
				session, _ := store.Get(r, c.name)
				session.Values["csrf"] = c.csrfValue
				session.Options.HttpOnly = true
				session.Options.Secure = true
				session.Options.SameSite = c.sameSite
				err := session.Save(r, w)
				assert.NoError(t, err, "failed to save cookie %s", c.name)
			}

			err := ValidateCSRFSession(t.Context(), r, config, store, name, tc.csrfValue, new(flow.Flow))
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateCsrfSession(t *testing.T) {
	type cookie struct {
		httpOnly bool
		secure   bool
		domain   string
		sameSite http.SameSite
		maxAge   int
	}
	for _, tc := range []struct {
		name                     string
		secure                   bool
		domain                   string
		sameSite                 http.SameSite
		maxAge                   time.Duration
		sameSiteLegacyWorkaround bool
		expectedCookies          map[string]cookie
	}{
		{
			name:                     "csrf_default",
			secure:                   true,
			sameSite:                 http.SameSiteDefaultMode,
			maxAge:                   10 * time.Second,
			sameSiteLegacyWorkaround: false,
			expectedCookies: map[string]cookie{
				"csrf_default": {
					httpOnly: true,
					secure:   true,
					sameSite: 0, // see https://golang.org/doc/go1.16#net/http
					maxAge:   10,
				},
			},
		},
		{
			name:                     "csrf_lax_insecure",
			secure:                   false,
			sameSite:                 http.SameSiteLaxMode,
			maxAge:                   20 * time.Second,
			sameSiteLegacyWorkaround: false,
			expectedCookies: map[string]cookie{
				"csrf_lax_insecure": {
					httpOnly: true,
					secure:   false,
					sameSite: http.SameSiteLaxMode,
					maxAge:   20,
				},
			},
		},
		{
			name:                     "csrf_none",
			secure:                   true,
			sameSite:                 http.SameSiteNoneMode,
			maxAge:                   30 * time.Second,
			sameSiteLegacyWorkaround: false,
			expectedCookies: map[string]cookie{
				"csrf_none": {
					httpOnly: true,
					secure:   true,
					sameSite: http.SameSiteNoneMode,
					maxAge:   30,
				},
			},
		},
		{
			name:                     "csrf_none_fallback",
			secure:                   true,
			sameSite:                 http.SameSiteNoneMode,
			maxAge:                   40 * time.Second,
			sameSiteLegacyWorkaround: true,
			expectedCookies: map[string]cookie{
				"csrf_none_fallback": {
					httpOnly: true,
					secure:   true,
					sameSite: http.SameSiteNoneMode,
					maxAge:   40,
				},
				"csrf_none_fallback_legacy": {
					httpOnly: true,
					secure:   true,
					sameSite: 0,
					maxAge:   40,
				},
			},
		},
		{
			name:                     "csrf_strict_fallback_ignored",
			secure:                   true,
			sameSite:                 http.SameSiteStrictMode,
			maxAge:                   50 * time.Second,
			sameSiteLegacyWorkaround: true,
			expectedCookies: map[string]cookie{
				"csrf_strict_fallback_ignored": {
					httpOnly: true,
					secure:   true,
					sameSite: http.SameSiteStrictMode,
					maxAge:   50,
				},
			},
		},
		{
			name:     "csrf_domain",
			secure:   true,
			domain:   "foobar.com",
			sameSite: http.SameSiteNoneMode,
			expectedCookies: map[string]cookie{
				"csrf_domain": {
					httpOnly: true,
					secure:   true,
					domain:   "foobar.com",
					sameSite: http.SameSiteNoneMode,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store := sessions.NewCookieStore(securecookie.GenerateRandomKey(16))
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			config := mock.NewMockCookieConfigProvider(ctrl)
			config.EXPECT().CookieSameSiteMode(gomock.Any()).Return(tc.sameSite).AnyTimes()
			config.EXPECT().CookieSameSiteLegacyWorkaround(gomock.Any()).Return(tc.sameSiteLegacyWorkaround).AnyTimes()
			config.EXPECT().CookieSecure(gomock.Any()).Return(tc.secure).AnyTimes()
			config.EXPECT().CookieDomain(gomock.Any()).Return(tc.domain).AnyTimes()

			err := createCSRFSession(t.Context(), rr, req, config, store, tc.name, "value", tc.maxAge)
			assert.NoError(t, err)

			cookies := make(map[string]cookie)
			for _, c := range rr.Result().Cookies() {
				cookies[c.Name] = cookie{
					httpOnly: c.HttpOnly,
					secure:   c.Secure,
					sameSite: c.SameSite,
					maxAge:   c.MaxAge,
					domain:   c.Domain,
				}
			}
			assert.Equal(t, tc.expectedCookies, cookies)
		})
	}
}

func TestCaseInsensitiveFilterParam(t *testing.T) {
	for k, tc := range []struct {
		requestedQuery string
		key            string

		expectedQuery url.Values
	}{
		{
			requestedQuery: "key=value",
			key:            "key2",
			expectedQuery:  url.Values{"key": []string{"value"}},
		},
		{
			requestedQuery: "KeY=value",
			key:            "key",
			expectedQuery:  url.Values{"KeY": []string{"****"}},
		},
		{
			requestedQuery: "KeY=value",
			key:            "kEy",
			expectedQuery:  url.Values{"KeY": []string{"****"}},
		},
		{
			requestedQuery: "key=value&KEY2=value2",
			key:            "key2",
			expectedQuery:  url.Values{"key": []string{"value"}, "KEY2": []string{"****"}},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			query, err := url.ParseQuery(tc.requestedQuery)
			assert.NoError(t, err)

			q := caseInsensitiveFilterParam(query, tc.key)

			assert.Equal(t, tc.expectedQuery, q)
		})
	}
}
