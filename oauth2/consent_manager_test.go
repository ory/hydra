package oauth2_test

import (
	"fmt"
	"log"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/compose"
	"github.com/ory/hydra/integration"
	. "github.com/ory/hydra/oauth2"
	"github.com/ory/ladon"
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

func TestConsentHttpClient(t *testing.T) {
	req := &ConsentRequest{
		ID:               "id-3",
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

	memm := NewConsentRequestMemoryManager()
	var localWarden, httpClient = compose.NewMockFirewall("foo", "app-client", fosite.Arguments{ConsentScope}, &ladon.DefaultPolicy{
		ID:        "1",
		Subjects:  []string{"app-client"},
		Resources: []string{"rn:hydra:oauth2:consent:requests:<.*>"},
		Actions:   []string{"get", "accept", "reject"},
		Effect:    ladon.AllowAccess,
	})

	require.NoError(t, memm.PersistConsentRequest(req))

	h := &ConsentSessionHandler{
		M: memm,
		W: localWarden,
		H: herodot.NewJSONWriter(nil),
	}

	r := httprouter.New()
	h.SetRoutes(r)
	ts := httptest.NewServer(r)
	u, _ := url.Parse(ts.URL + ConsentRequestPath)

	m := HTTPConsentManager{
		Client:   httpClient,
		Endpoint: u,
	}

	got, err := m.GetConsentRequest(req.ID)
	require.NoError(t, err)
	assert.EqualValues(t, req.ID, got.ID)
	assert.EqualValues(t, req.Audience, got.Audience)
	assert.EqualValues(t, req.RequestedScopes, got.RequestedScopes)
	assert.EqualValues(t, req.ExpiresAt, got.ExpiresAt)
	assert.EqualValues(t, req.RedirectURL, got.RedirectURL)
	assert.False(t, got.IsConsentGranted())

	accept := &AcceptConsentRequestPayload{
		Subject:          "some-subject",
		GrantScopes:      []string{"scope1", "scope2"},
		AccessTokenExtra: map[string]interface{}{"at": "bar"},
		IDTokenExtra:     map[string]interface{}{"id": "bar"},
	}

	require.NoError(t, m.AcceptConsentRequest(req.ID, accept))
	got, err = memm.GetConsentRequest(req.ID)
	require.NoError(t, err)
	assert.Equal(t, accept.Subject, got.Subject)
	assert.Equal(t, accept.GrantScopes, got.GrantedScopes)
	assert.Equal(t, accept.AccessTokenExtra, got.AccessTokenExtra)
	assert.Equal(t, accept.IDTokenExtra, got.IDTokenExtra)
	assert.True(t, got.IsConsentGranted())

	require.NoError(t, m.RejectConsentRequest(req.ID, &RejectConsentRequestPayload{
		Reason: "MyReason",
	}))
	got, err = memm.GetConsentRequest(req.ID)
	require.NoError(t, err)
	assert.Equal(t, "MyReason", got.DenyReason)
	assert.False(t, got.IsConsentGranted())
}
