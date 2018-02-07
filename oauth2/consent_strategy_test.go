// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package oauth2

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"github.com/ory/fosite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsentStrategy(t *testing.T) {
	t.Run("suite=validate", func(t *testing.T) {

		strategy := &DefaultConsentStrategy{ConsentManager: NewConsentRequestMemoryManager()}

		require.NoError(t, strategy.ConsentManager.PersistConsentRequest(&ConsentRequest{
			ID:      "not_granted",
			Consent: ConsentRequestRejected,
		}))
		require.NoError(t, strategy.ConsentManager.PersistConsentRequest(&ConsentRequest{
			ID:        "granted",
			Consent:   ConsentRequestAccepted,
			ClientID:  "client_id",
			Subject:   "peter",
			CSRF:      "csrf_token",
			ExpiresAt: time.Now().Add(time.Hour),
		}))
		require.NoError(t, strategy.ConsentManager.PersistConsentRequest(&ConsentRequest{
			ID:        "granted_csrf_cookie",
			Consent:   ConsentRequestAccepted,
			ClientID:  "client_id",
			Subject:   "peter",
			CSRF:      "csrf_token",
			ExpiresAt: time.Now().Add(time.Hour),
		}))
		require.NoError(t, strategy.ConsentManager.PersistConsentRequest(&ConsentRequest{
			ID:        "granted_csrf_request",
			Consent:   ConsentRequestAccepted,
			ClientID:  "client_id",
			Subject:   "peter",
			CSRF:      "csrf_token",
			ExpiresAt: time.Now().Add(time.Hour),
		}))
		require.NoError(t, strategy.ConsentManager.PersistConsentRequest(&ConsentRequest{
			ID:        "granted_expired",
			Consent:   ConsentRequestAccepted,
			Subject:   "peter",
			ClientID:  "client_id",
			ExpiresAt: time.Now().Add(-time.Hour),
			CSRF:      "csrf_token",
		}))

		for _, tc := range []struct {
			req       *fosite.AuthorizeRequest
			session   string
			cookie    *sessions.Session
			expectErr bool
			assert    func(*testing.T, *Session)
			d         string
		}{
			{
				d:         "invalid session",
				session:   "not_granted",
				expectErr: true,
				cookie:    &sessions.Session{Values: map[interface{}]interface{}{CookieCSRFKey: "csrf_token"}},
			},
			{
				d:         "session expired",
				session:   "granted_expired",
				expectErr: true,
				req:       &fosite.AuthorizeRequest{Request: fosite.Request{Client: &fosite.DefaultClient{ID: "client_id"}, Form: url.Values{"consent_csrf": {"csrf_token"}}}},
				cookie:    &sessions.Session{Values: map[interface{}]interface{}{CookieCSRFKey: "csrf_token"}},
			},
			{
				d:         "granted",
				session:   "granted",
				expectErr: false,
				req:       &fosite.AuthorizeRequest{Request: fosite.Request{Client: &fosite.DefaultClient{ID: "client_id"}, Form: url.Values{"consent_csrf": {"csrf_token"}}}},
				cookie:    &sessions.Session{Values: map[interface{}]interface{}{CookieCSRFKey: "csrf_token"}},
			},
			{
				d:         "client mismatch",
				session:   "granted",
				expectErr: true,
				req:       &fosite.AuthorizeRequest{Request: fosite.Request{Client: &fosite.DefaultClient{ID: "mismatch_client"}, Form: url.Values{"consent_csrf": {"csrf_token"}}}},
				cookie:    &sessions.Session{Values: map[interface{}]interface{}{CookieCSRFKey: "csrf_token"}},
			},
			{
				d:         "consent request was not initiated by this user agent",
				session:   "granted_csrf_cookie",
				expectErr: true,
				req:       &fosite.AuthorizeRequest{Request: fosite.Request{Client: &fosite.DefaultClient{ID: "client_id"}, Form: url.Values{"consent_csrf": {"csrf_token"}}}},
				cookie:    &sessions.Session{Values: map[interface{}]interface{}{CookieCSRFKey: "very_different_csrf_token"}},
				assert: func(t *testing.T, session *Session) {
					cr, err := strategy.ConsentManager.GetConsentRequest("granted_csrf_cookie")
					require.NoError(t, err)
					assert.False(t, cr.IsConsentGranted())
				},
			},
			{
				d:         "authorize url contains csrf token that does not match token set in consent request",
				session:   "granted_csrf_request",
				expectErr: true,
				req:       &fosite.AuthorizeRequest{Request: fosite.Request{Client: &fosite.DefaultClient{ID: "client_id"}, Form: url.Values{"consent_csrf": {"very_different_csrf_token"}}}},
				cookie:    &sessions.Session{Values: map[interface{}]interface{}{CookieCSRFKey: "csrf_token"}},
				assert: func(t *testing.T, session *Session) {
					cr, err := strategy.ConsentManager.GetConsentRequest("granted_csrf_request")
					require.NoError(t, err)
					assert.False(t, cr.IsConsentGranted())
				},
			},
		} {
			t.Run(fmt.Sprintf("case=%s", tc.d), func(t *testing.T) {
				res, err := strategy.ValidateConsentRequest(tc.req, tc.session, tc.cookie)
				if tc.expectErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					if tc.assert != nil {
						tc.assert(t, res)
					}
				}
			})
		}
	})
}
