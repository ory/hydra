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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/compose"
	. "github.com/ory/hydra/oauth2"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/ory/ladon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsentSDK(t *testing.T) {
	req := &ConsentRequest{
		ID:               "id-3",
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

	memm := NewConsentRequestMemoryManager()
	var localWarden, httpClient = compose.NewMockFirewall("foo", "app-client", fosite.Arguments{ConsentScope}, &ladon.DefaultPolicy{
		ID:        "1",
		Subjects:  []string{"app-client"},
		Resources: []string{"rn:hydra:oauth2:consent:requests:<.*>"},
		Actions:   []string{"get", "accept", "reject"},
		Effect:    ladon.AllowAccess,
	})

	require.NoError(t, memm.PersistConsentRequest(req))
	h := &ConsentSessionHandler{M: memm, W: localWarden, H: herodot.NewJSONWriter(nil)}

	r := httprouter.New()
	h.SetRoutes(r)
	server := httptest.NewServer(r)

	client := hydra.NewOAuth2ApiWithBasePath(server.URL)
	client.Configuration.Transport = httpClient.Transport

	got, _, err := client.GetOAuth2ConsentRequest(req.ID)
	require.NoError(t, err)
	assert.EqualValues(t, req.ID, got.Id)
	assert.EqualValues(t, req.ClientID, got.ClientId)
	assert.EqualValues(t, req.RequestedScopes, got.RequestedScopes)
	assert.EqualValues(t, req.RedirectURL, got.RedirectUrl)

	accept := hydra.ConsentRequestAcceptance{
		Subject:          "some-subject",
		GrantScopes:      []string{"scope1", "scope2"},
		AccessTokenExtra: map[string]interface{}{"at": "bar"},
		IdTokenExtra:     map[string]interface{}{"id": "bar"},
	}

	response, err := client.AcceptOAuth2ConsentRequest(req.ID, accept)
	require.NoError(t, err)
	assert.EqualValues(t, http.StatusNoContent, response.StatusCode)

	gotMem, err := memm.GetConsentRequest(req.ID)
	require.NoError(t, err)
	assert.Equal(t, accept.Subject, gotMem.Subject)
	assert.Equal(t, accept.GrantScopes, gotMem.GrantedScopes)
	assert.Equal(t, accept.AccessTokenExtra, gotMem.AccessTokenExtra)
	assert.Equal(t, accept.IdTokenExtra, gotMem.IDTokenExtra)
	assert.True(t, gotMem.IsConsentGranted())

	response, err = client.RejectOAuth2ConsentRequest(req.ID, hydra.ConsentRequestRejection{Reason: "MyReason"})
	require.NoError(t, err)
	assert.EqualValues(t, http.StatusNoContent, response.StatusCode)

	gotMem, err = memm.GetConsentRequest(req.ID)
	require.NoError(t, err)
	assert.Equal(t, "MyReason", gotMem.DenyReason)
	assert.False(t, gotMem.IsConsentGranted())
}
