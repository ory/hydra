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

package oauth2_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ory/hydra/integration"
	. "github.com/ory/hydra/oauth2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var consentManagers = map[string]ConsentRequestManager{
	"memory": NewConsentRequestMemoryManager(),
}

func connectToMySQLConsent() {
	var db = integration.ConnectToMySQL()
	s := NewConsentRequestSQLManager(db)

	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	consentManagers["mysql"] = s
}

func connectToPGConsent() {
	var db = integration.ConnectToPostgres()
	s := NewConsentRequestSQLManager(db)

	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	consentManagers["postgres"] = s
}

func TestConsentRequestManagerReadWrite(t *testing.T) {
	req := &ConsentRequest{
		ID:               "id-1",
		ClientID:         "client-id",
		RequestedScopes:  []string{"foo", "bar"},
		GrantedScopes:    []string{"baz", "bar"},
		CSRF:             "some-csrf",
		ExpiresAt:        time.Now().Round(time.Minute),
		Consent:          ConsentRequestAccepted,
		DenyReason:       "some reason",
		AccessTokenExtra: map[string]interface{}{"atfoo": "bar", "atbaz": "bar"},
		IDTokenExtra:     map[string]interface{}{"idfoo": "bar", "idbaz": "bar"},
		RedirectURL:      "https://redirect-me/foo",
		Subject:          "Peter",
	}

	for k, m := range consentManagers {
		t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
			_, err := m.GetConsentRequest("1234")
			assert.Error(t, err)

			require.NoError(t, m.PersistConsentRequest(req))

			got, err := m.GetConsentRequest(req.ID)
			require.NoError(t, err)

			require.Equal(t, req.ExpiresAt.Unix(), got.ExpiresAt.Unix())
			got.ExpiresAt = req.ExpiresAt
			assert.EqualValues(t, req, got)
		})
	}
}

func TestConsentRequestManagerUpdate(t *testing.T) {
	req := &ConsentRequest{
		ID:               "id-2",
		ClientID:         "client-id",
		RequestedScopes:  []string{"foo", "bar"},
		GrantedScopes:    []string{"baz", "bar"},
		CSRF:             "some-csrf",
		ExpiresAt:        time.Now().Round(time.Minute),
		Consent:          ConsentRequestRejected,
		DenyReason:       "some reason",
		AccessTokenExtra: map[string]interface{}{"atfoo": "bar", "atbaz": "bar"},
		IDTokenExtra:     map[string]interface{}{"idfoo": "bar", "idbaz": "bar"},
		RedirectURL:      "https://redirect-me/foo",
		Subject:          "Peter",
	}

	for k, m := range consentManagers {
		t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
			require.NoError(t, m.PersistConsentRequest(req))

			got, err := m.GetConsentRequest(req.ID)
			require.NoError(t, err)
			assert.False(t, got.IsConsentGranted())
			require.Equal(t, req.ExpiresAt.Unix(), got.ExpiresAt.Unix())
			got.ExpiresAt = req.ExpiresAt
			assert.EqualValues(t, req, got)

			require.NoError(t, m.AcceptConsentRequest(req.ID, new(AcceptConsentRequestPayload)))
			got, err = m.GetConsentRequest(req.ID)
			require.NoError(t, err)
			assert.True(t, got.IsConsentGranted())

			require.NoError(t, m.RejectConsentRequest(req.ID, new(RejectConsentRequestPayload)))
			got, err = m.GetConsentRequest(req.ID)
			require.NoError(t, err)
			assert.False(t, got.IsConsentGranted())
		})
	}
}

func TestConsentRequestManagerPreviousConsent(t *testing.T) {
	reqs := []ConsentRequest{
		{
			ID:            "id-1",
			ClientID:      "client-1",
			GrantedScopes: []string{"baz", "bar"},
			RequestedAt:   time.Now().UTC().Round(time.Millisecond),
			Subject:       "peter",
			Consent:       ConsentRequestAccepted,
		},
		{
			ID:            "id-2",
			ClientID:      "client-1",
			GrantedScopes: []string{"baz", "bar"},
			RequestedAt:   time.Now().Add(-time.Hour).UTC().Round(time.Millisecond),
			Subject:       "peter",
			Consent:       ConsentRequestAccepted,
		},
		{
			ID:            "id-3",
			ClientID:      "client-2",
			GrantedScopes: []string{"baz", "bar"},
			RequestedAt:   time.Now().UTC().Round(time.Millisecond),
			Subject:       "peter",
			Consent:       ConsentRequestRejected,
		},
		{
			ID:            "id-4",
			ClientID:      "client-2",
			GrantedScopes: []string{"baz", "bar"},
			RequestedAt:   time.Now().UTC().Round(time.Millisecond),
			Subject:       "alice",
			Consent:       ConsentRequestAccepted,
		},
		{
			ID:            "id-5",
			ClientID:      "client-3",
			GrantedScopes: []string{},
			RequestedAt:   time.Now().UTC().Round(time.Millisecond),
			Subject:       "alice",
			Consent:       ConsentRequestAccepted,
		},
	}

	for k, m := range consentManagers {
		t.Run(fmt.Sprintf("manager=%s", k), func(t *testing.T) {
			for _, r := range reqs {
				require.NoError(t, m.PersistConsentRequest(&r))
			}

			t.Run("case=can not find a matching consent", func(t *testing.T) {
				session, err := m.GetPreviouslyGrantedConsent("peter", "nonexistent-client", []string{})
				require.Error(t, err)
				assert.Nil(t, session)

				session, err = m.GetPreviouslyGrantedConsent("peter", "client-1", []string{"baz", "bar", "foo"})
				require.Error(t, err)
				assert.Nil(t, session)
			})

			t.Run("case=returns the newer of two consents", func(t *testing.T) {
				session, err := m.GetPreviouslyGrantedConsent("peter", "client-1", []string{"baz", "bar"})
				require.NoError(t, err)
				assert.Equal(t, &reqs[0], session)
			})

			t.Run("case=returns the only valid consent for client-2", func(t *testing.T) {
				session, err := m.GetPreviouslyGrantedConsent("peter", "client-2", []string{"baz", "bar"})
				require.Error(t, err)
				assert.Nil(t, session)

				session, err = m.GetPreviouslyGrantedConsent("alice", "client-2", []string{"baz", "bar"})
				require.NoError(t, err)
				assert.Equal(t, &reqs[3], session)
			})
			t.Run("case=returns the only valid consent for client-3", func(t *testing.T) {
				session, err := m.GetPreviouslyGrantedConsent("alice", "client-3", []string{"baz"})
				require.Error(t, err)
				assert.Nil(t, session)

				session, err = m.GetPreviouslyGrantedConsent("alice", "client-3", []string{})
				require.NoError(t, err)
				assert.Equal(t, &reqs[4], session)
			})
		})
	}
}
