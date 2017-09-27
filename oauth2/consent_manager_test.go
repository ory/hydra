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
		Audience:         "audience",
		RequestedScopes:  []string{"foo", "bar"},
		GrantedScopes:    []string{"baz", "bar"},
		CSRF:             "some-csrf",
		ExpiresAt:        time.Now().Round(time.Second),
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

			assert.EqualValues(t, req, got)
		})
	}
}

func TestConsentRequestManagerUpdate(t *testing.T) {
	req := &ConsentRequest{
		ID:               "id-2",
		Audience:         "audience",
		RequestedScopes:  []string{"foo", "bar"},
		GrantedScopes:    []string{"baz", "bar"},
		CSRF:             "some-csrf",
		ExpiresAt:        time.Now().Round(time.Second),
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
