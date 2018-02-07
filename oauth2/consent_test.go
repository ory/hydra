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
	"sort"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSort(t *testing.T) {
	now := time.Now()
	items := []ConsentRequest{
		{RequestedAt: now.Add(-time.Hour * 3)},
		{RequestedAt: now.Add(-time.Hour * 2)},
		{RequestedAt: now.Add(-time.Hour * 1)},
	}

	sort.Sort(byTime(items))

	assert.Equal(t, now.Add(-time.Hour*1), items[0].RequestedAt)
}

func TestWhiteList(t *testing.T) {
	assert.True(t, containsWhiteListedOnly([]string{"a", "b"}, []string{"a", "b"}))
	assert.True(t, containsWhiteListedOnly([]string{}, []string{"a", "b"}))
	assert.False(t, containsWhiteListedOnly([]string{"a", "b", "c"}, []string{"a", "b"}))
	assert.False(t, containsWhiteListedOnly([]string{"a", "c", "b"}, []string{"a", "b"}))
}

func TestValidateConsentRequest(t *testing.T) {
	manager := NewConsentRequestMemoryManager()

	strategy := &DefaultConsentStrategy{
		Issuer:                   "tests",
		DefaultIDTokenLifespan:   time.Hour,
		DefaultChallengeLifespan: time.Hour,
		ConsentManager:           manager,
	}

	defaultClient := &client.Client{ID: "some-client"}

	for k, tc := range []struct {
		d          string
		req        fosite.AuthorizeRequest
		expectsErr bool
		accept     bool
		assert     func(*testing.T, *Session)
	}{
		{
			d: "should pass because authorization is granted",
			req: fosite.AuthorizeRequest{
				Request: fosite.Request{
					Client: defaultClient,
					Scopes: fosite.Arguments{"scope-a", "scope-b"},
					Form: url.Values{
						"prompt":     []string{"consent"},
						"display":    []string{"popup"},
						"ui_locales": []string{"de"},
						"login_hint": []string{"hint"},
						"acr_values": []string{"0 1"},
					},
				},
			},
			accept: true,
		},
		{
			d: "should fail because authorization is denied",
			req: fosite.AuthorizeRequest{
				Request: fosite.Request{
					Client: defaultClient,
					Scopes: fosite.Arguments{"scope-a", "scope-b"},
					Form: url.Values{
						"prompt":     []string{"consent"},
						"display":    []string{"popup"},
						"ui_locales": []string{"de"},
						"login_hint": []string{"hint"},
						"acr_values": []string{"0 1"},
					},
				},
			},
			expectsErr: true,
			accept:     false,
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
			session := sessions.NewSession(nil, "foo")
			id, err := strategy.CreateConsentRequest(&tc.req, "redirect", session)
			assert.NoError(t, err)

			if tc.accept {
				manager.AcceptConsentRequest(id, &AcceptConsentRequestPayload{Subject: "foo"})
			} else {
				manager.RejectConsentRequest(id, new(RejectConsentRequestPayload))
			}

			claims, err := strategy.ValidateConsentRequest(&tc.req, id, session)
			if tc.expectsErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tc.assert != nil {
					tc.assert(t, claims)
				}
			}
		})
	}
}

func TestCreateConsentRequest(t *testing.T) {
	manager := NewConsentRequestMemoryManager()

	strategy := &DefaultConsentStrategy{
		Issuer:                   "tests",
		DefaultIDTokenLifespan:   time.Hour,
		DefaultChallengeLifespan: time.Hour,
		ConsentManager:           manager,
	}

	defaultClient := &client.Client{ID: "some-client"}
	req := &fosite.AuthorizeRequest{
		Request: fosite.Request{
			Client: defaultClient,
			Scopes: fosite.Arguments{"scope-a", "scope-b"},
			Form: url.Values{
				"prompt":     []string{"consent"},
				"display":    []string{"popup"},
				"ui_locales": []string{"de"},
				"login_hint": []string{"hint"},
				"acr_values": []string{"0 1"},
			},
		},
	}

	session := sessions.NewSession(nil, "foo")
	id, err := strategy.CreateConsentRequest(req, "redirect", session)
	assert.NoError(t, err)

	assert.NotEmpty(t, id)
	assert.NotEmpty(t, session.Values[cookieCSRFKey])

	fromStore := manager.requests[id]
	assert.Equal(t, id, fromStore.ID)
	assert.Equal(t, req.Request.Scopes, fosite.Arguments(fromStore.RequestedScopes))
	assert.Equal(t, defaultClient.ID, fromStore.ClientID)
	assert.Equal(t, defaultClient, fromStore.Client)
	assert.Empty(t, fromStore.GrantedScopes)
	assert.Equal(t, &ConsentRequestOpenIDConnectContext{
		Prompt:    "consent",
		Display:   "popup",
		UILocales: "de",
		LoginHint: "hint",
		ACRValues: []string{"0", "1"},
	}, fromStore.OpenIDConnectContext)
}

func TestHandleConsentRequest(t *testing.T) {
	manager := NewConsentRequestMemoryManager()

	strategy := &DefaultConsentStrategy{
		Issuer:                   "tests",
		DefaultIDTokenLifespan:   time.Hour,
		DefaultChallengeLifespan: time.Hour,
		ConsentManager:           manager,
	}

	defaultClient := &client.Client{ID: "some-client"}

	var newRequest = func(u url.Values) *fosite.AuthorizeRequest {
		return &fosite.AuthorizeRequest{
			Request: fosite.Request{
				Client: defaultClient,
				Form:   u,
			},
		}
	}

	anHourAgo := time.Now().UTC().Round(time.Second).Add(-time.Hour)

	for k, tc := range []struct {
		d            string
		req          fosite.AuthorizeRequester
		cookie       sessions.Session
		expectErr    error
		requiresAuth bool
		assert       func(*testing.T, *Session, fosite.AuthorizeRequester)
		prepare      func(*testing.T)
	}{
		{
			d:         "should fail if prompt=none and no session was set",
			req:       newRequest(url.Values{"prompt": []string{"none"}}),
			expectErr: fosite.ErrLoginRequired,
		},
		{
			d:         "should fail because unknown prompt was used",
			req:       newRequest(url.Values{"prompt": []string{"foo"}}),
			expectErr: fosite.ErrInvalidRequest,
		},
		{
			d:         "should fail because unknown prompt was used",
			req:       newRequest(url.Values{"prompt": []string{"login foo"}}),
			expectErr: fosite.ErrInvalidRequest,
		},
		{
			d:         "should fail because prompt=none can not be used together with another prompt",
			req:       newRequest(url.Values{"prompt": []string{"none login"}}),
			expectErr: fosite.ErrInvalidRequest,
		},
		{
			d:            "should require authentication when prompt is login",
			req:          newRequest(url.Values{"prompt": []string{"login"}}),
			requiresAuth: true,
		},
		{
			d:            "should require authentication when prompt is consent",
			req:          newRequest(url.Values{"prompt": []string{"consent"}}),
			requiresAuth: true,
		},
		{
			d: "should require authentication when client is public",
			req: &fosite.AuthorizeRequest{
				Request: fosite.Request{
					Client: &client.Client{Public: true},
					Form:   url.Values{},
				},
			},
			requiresAuth: true,
		},
		{
			d:            "should require authentication when prompt is select_account",
			req:          newRequest(url.Values{"prompt": []string{"select_account"}}),
			requiresAuth: true,
		},
		{
			d:            "should require authentication when prompt is any of login consent select_account",
			req:          newRequest(url.Values{"prompt": []string{"login consent select_account"}}),
			requiresAuth: true,
		},
		{
			d:            "should require authentication when no session is given",
			req:          newRequest(url.Values{}),
			requiresAuth: true,
		},
		{
			d:            "should require authentication session contains invalid user value",
			req:          newRequest(url.Values{}),
			requiresAuth: true,
			cookie:       sessions.Session{Values: map[interface{}]interface{}{sessionUserKey: 1234}},
		},
		{
			d:            "should require authentication session contains valid user but not an auth time",
			req:          newRequest(url.Values{}),
			requiresAuth: true,
			cookie:       sessions.Session{Values: map[interface{}]interface{}{sessionUserKey: "some-user"}},
		},
		{
			d:            "should require authentication session contains valid user but not a valid auth time",
			req:          newRequest(url.Values{}),
			requiresAuth: true,
			cookie:       sessions.Session{Values: map[interface{}]interface{}{sessionUserKey: "some-user", sessionAuthTimeKey: 123}},
		},
		{
			d:            "should require authentication because maxAge was reached",
			req:          newRequest(url.Values{"maxAge": []string{"10"}}),
			requiresAuth: true,
			cookie:       sessions.Session{Values: map[interface{}]interface{}{sessionUserKey: "some-user", sessionAuthTimeKey: time.Now().UTC().Add(-time.Hour).Format(time.RFC3339)}},
		},
		{
			d:         "should fail because prompt is none but maxAge was reached",
			req:       newRequest(url.Values{"maxAge": []string{"10"}, "prompt": []string{"none"}}),
			expectErr: fosite.ErrLoginRequired,
			cookie:    sessions.Session{Values: map[interface{}]interface{}{sessionUserKey: "some-user", sessionAuthTimeKey: time.Now().UTC().Add(-time.Hour).Format(time.RFC3339)}},
		},
		{
			d:            "should require authentication because the request has not been granted previously",
			req:          newRequest(url.Values{}),
			requiresAuth: true,
			cookie:       sessions.Session{Values: map[interface{}]interface{}{sessionUserKey: "some-user", sessionAuthTimeKey: time.Now().UTC().Format(time.RFC3339)}},
		},
		{
			d:      "should create a session based on previous authorization",
			req:    newRequest(url.Values{}),
			cookie: sessions.Session{Values: map[interface{}]interface{}{sessionUserKey: "some-user", sessionAuthTimeKey: time.Now().UTC().Format(time.RFC3339)}},
			prepare: func(t *testing.T) {
				assert.NoError(t, manager.PersistConsentRequest(&ConsentRequest{
					ID:          "req-1",
					ClientID:    "some-client",
					Client:      defaultClient,
					Subject:     "some-user",
					RequestedAt: time.Now(),
					Consent:     ConsentRequestAccepted,
				}))
			},
		},
		{
			d:            "should require authentication because the request has not been granted previously",
			req:          newRequest(url.Values{}),
			requiresAuth: true,
			cookie:       sessions.Session{Values: map[interface{}]interface{}{sessionUserKey: "another-user", sessionAuthTimeKey: time.Now().UTC().Format(time.RFC3339)}},
		},
		{
			d: "should create a session based on previous authorization",
			req: &fosite.AuthorizeRequest{
				Request: fosite.Request{Client: defaultClient, Form: url.Values{}, Scopes: []string{"foo", "bar"}},
			},
			cookie: sessions.Session{Values: map[interface{}]interface{}{sessionUserKey: "some-user", sessionAuthTimeKey: time.Now().UTC().Format(time.RFC3339)}},
			prepare: func(t *testing.T) {
				assert.NoError(t, manager.PersistConsentRequest(&ConsentRequest{
					ID:            "req-2",
					ClientID:      "some-client",
					Client:        defaultClient,
					Subject:       "some-user",
					GrantedScopes: []string{"foo", "bar"},
					RequestedAt:   time.Now(),
					Consent:       ConsentRequestAccepted,
				}))
			},
		},
		{
			d: "should create a session based on previous authorization",
			req: &fosite.AuthorizeRequest{
				Request: fosite.Request{Client: defaultClient, Form: url.Values{}, Scopes: []string{"foo"}},
			},
			cookie: sessions.Session{Values: map[interface{}]interface{}{sessionUserKey: "some-user", sessionAuthTimeKey: time.Now().UTC().Format(time.RFC3339)}},
			prepare: func(t *testing.T) {
				assert.NoError(t, manager.PersistConsentRequest(&ConsentRequest{
					ID:               "req-2",
					ClientID:         "some-client",
					Client:           defaultClient,
					Subject:          "some-user",
					GrantedScopes:    []string{"foo", "bar"},
					RequestedAt:      anHourAgo,
					IDTokenExtra:     map[string]interface{}{"id": "foo"},
					AccessTokenExtra: map[string]interface{}{"token": "foo"},
					Consent:          ConsentRequestAccepted,
				}))
			},
			assert: func(t *testing.T, session *Session, req fosite.AuthorizeRequester) {
				assert.EqualValues(t, req.GetGrantedScopes(), []string{"foo"})
				assert.Equal(t, "0", session.DefaultSession.Claims.AuthenticationContextClassReference)
				assert.Equal(t, "some-user", session.DefaultSession.Subject)
				assert.Equal(t, anHourAgo, session.DefaultSession.Claims.AuthTime)
				assert.Equal(t, "some-client", session.DefaultSession.Claims.Audience)
				assert.EqualValues(t, map[string]interface{}{"id": "foo"}, session.DefaultSession.Claims.Extra)
				assert.EqualValues(t, map[string]interface{}{"id": "foo"}, session.DefaultSession.Claims.Extra)
				assert.EqualValues(t, map[string]interface{}{"token": "foo"}, session.Extra)
			},
		},
		{
			d: "should fail creating a session as requested scopes have not been granted in the past",
			req: &fosite.AuthorizeRequest{
				Request: fosite.Request{Client: defaultClient, Form: url.Values{}, Scopes: []string{"foo", "bar", "baz"}},
			},
			requiresAuth: true,
			cookie:       sessions.Session{Values: map[interface{}]interface{}{sessionUserKey: "some-user", sessionAuthTimeKey: time.Now().UTC().Format(time.RFC3339)}},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
			if tc.prepare != nil {
				tc.prepare(t)
			}

			session, err := strategy.HandleConsentRequest(
				tc.req,
				&tc.cookie,
			)

			if tc.expectErr != nil {
				assert.EqualError(t, err, tc.expectErr.Error())
			} else if tc.requiresAuth {
				assert.EqualError(t, err, ErrRequiresAuthentication.Error())
			} else {
				require.NoError(t, err)
				if tc.assert != nil {
					tc.assert(t, session, tc.req)
				}
			}
		})
	}
}
