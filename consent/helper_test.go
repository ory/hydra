/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package consent

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/ory/hydra/internal/mock"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
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
		granted         []AcceptOAuth2ConsentRequest
		requested       []string
		expectChallenge string
	}{
		{
			granted:         []AcceptOAuth2ConsentRequest{{ID: "1", GrantedScope: []string{"foo", "bar"}}},
			requested:       []string{"foo", "bar"},
			expectChallenge: "1",
		},
		{
			granted:         []AcceptOAuth2ConsentRequest{{ID: "1", GrantedScope: []string{"foo", "bar"}}},
			requested:       []string{"foo", "bar", "baz"},
			expectChallenge: "",
		},
		{
			granted: []AcceptOAuth2ConsentRequest{
				{ID: "1", GrantedScope: []string{"foo", "bar"}},
				{ID: "2", GrantedScope: []string{"foo", "bar"}},
			},
			requested:       []string{"foo", "bar"},
			expectChallenge: "1",
		},
		{
			granted: []AcceptOAuth2ConsentRequest{
				{ID: "1", GrantedScope: []string{"foo", "bar"}},
				{ID: "2", GrantedScope: []string{"foo", "bar", "baz"}},
			},
			requested:       []string{"foo", "bar", "baz"},
			expectChallenge: "2",
		},
		{
			granted: []AcceptOAuth2ConsentRequest{
				{ID: "1", GrantedScope: []string{"foo", "bar"}},
				{ID: "2", GrantedScope: []string{"foo", "bar", "baz"}},
			},
			requested:       []string{"zab"},
			expectChallenge: "",
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			got := matchScopes(fosite.ExactScopeStrategy, tc.granted, tc.requested)
			if tc.expectChallenge == "" {
				assert.Nil(t, got)
				return
			}
			assert.Equal(t, tc.expectChallenge, got.ID)
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
					name:      legacyCsrfSessionName(name),
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
					name:      legacyCsrfSessionName(name),
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
					name:      legacyCsrfSessionName(name),
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
					name:      legacyCsrfSessionName(name),
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
					name:      legacyCsrfSessionName(name),
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

			err := validateCsrfSession(r, config, store, name, tc.csrfValue)
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
	}
	for _, tc := range []struct {
		name                     string
		secure                   bool
		domain                   string
		sameSite                 http.SameSite
		sameSiteLegacyWorkaround bool
		expectedCookies          map[string]cookie
	}{
		{
			name:                     "csrf_default",
			secure:                   true,
			sameSite:                 http.SameSiteDefaultMode,
			sameSiteLegacyWorkaround: false,
			expectedCookies: map[string]cookie{
				"csrf_default": {
					httpOnly: true,
					secure:   true,
					sameSite: 0, // see https://golang.org/doc/go1.16#net/http
				},
			},
		},
		{
			name:                     "csrf_lax_insecure",
			secure:                   false,
			sameSite:                 http.SameSiteLaxMode,
			sameSiteLegacyWorkaround: false,
			expectedCookies: map[string]cookie{
				"csrf_lax_insecure": {
					httpOnly: true,
					secure:   false,
					sameSite: http.SameSiteLaxMode,
				},
			},
		},
		{
			name:                     "csrf_none",
			secure:                   true,
			sameSite:                 http.SameSiteNoneMode,
			sameSiteLegacyWorkaround: false,
			expectedCookies: map[string]cookie{
				"csrf_none": {
					httpOnly: true,
					secure:   true,
					sameSite: http.SameSiteNoneMode,
				},
			},
		},
		{
			name:                     "csrf_none_fallback",
			secure:                   true,
			sameSite:                 http.SameSiteNoneMode,
			sameSiteLegacyWorkaround: true,
			expectedCookies: map[string]cookie{
				"csrf_none_fallback": {
					httpOnly: true,
					secure:   true,
					sameSite: http.SameSiteNoneMode,
				},
				"csrf_none_fallback_legacy": {
					httpOnly: true,
					secure:   true,
					sameSite: 0,
				},
			},
		},
		{
			name:                     "csrf_strict_fallback_ignored",
			secure:                   true,
			sameSite:                 http.SameSiteStrictMode,
			sameSiteLegacyWorkaround: true,
			expectedCookies: map[string]cookie{
				"csrf_strict_fallback_ignored": {
					httpOnly: true,
					secure:   true,
					sameSite: http.SameSiteStrictMode,
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

			err := createCsrfSession(rr, req, config, store, tc.name, "value")
			assert.NoError(t, err)

			cookies := make(map[string]cookie)
			for _, c := range rr.Result().Cookies() {
				cookies[c.Name] = cookie{
					httpOnly: c.HttpOnly,
					secure:   c.Secure,
					sameSite: c.SameSite,
					domain:   c.Domain,
				}
			}
			assert.Equal(t, tc.expectedCookies, cookies)
		})
	}
}
