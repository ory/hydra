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

package client_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/compose"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/ory/ladon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestClient(prefix string) hydra.OAuth2Client {
	return hydra.OAuth2Client{
		Id:            "1234",
		ClientName:    prefix + "name",
		ClientSecret:  prefix + "secret",
		ClientUri:     prefix + "uri",
		Contacts:      []string{prefix + "peter", prefix + "pan"},
		GrantTypes:    []string{prefix + "client_credentials", prefix + "authorize_code"},
		LogoUri:       prefix + "logo",
		Owner:         prefix + "an-owner",
		PolicyUri:     prefix + "policy-uri",
		Scope:         prefix + "foo bar baz",
		TosUri:        prefix + "tos-uri",
		ResponseTypes: []string{prefix + "id_token", prefix + "code"},
		RedirectUris:  []string{prefix + "redirect-url", prefix + "redirect-uri"},
	}
}

func TestClientSDK(t *testing.T) {
	manager := client.NewMemoryManager(nil)

	localWarden, httpClient := compose.NewMockFirewall("foo", "alice", fosite.Arguments{client.Scope}, &ladon.DefaultPolicy{
		ID:        "1",
		Subjects:  []string{"alice"},
		Resources: []string{"rn:hydra:clients<.*>"},
		Actions:   []string{"create", "get", "delete", "update"},
		Effect:    ladon.AllowAccess,
	})

	handler := &client.Handler{
		Manager: manager,
		H:       herodot.NewJSONWriter(nil),
		W:       localWarden,
	}

	router := httprouter.New()
	handler.SetRoutes(router)
	server := httptest.NewServer(router)
	c := hydra.NewOAuth2ApiWithBasePath(server.URL)
	c.Configuration.Transport = httpClient.Transport

	t.Run("foo", func(t *testing.T) {
		createClient := createTestClient("")

		result, _, err := c.CreateOAuth2Client(createClient)
		require.NoError(t, err)
		assert.EqualValues(t, createClient, *result)

		compareClient := createClient
		compareClient.ClientSecret = ""
		result, _, err = c.GetOAuth2Client(createClient.Id)
		assert.EqualValues(t, compareClient, *result)

		results, _, err := c.ListOAuth2Clients()
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.EqualValues(t, compareClient, results[0])

		updateClient := createTestClient("foo")
		result, _, err = c.UpdateOAuth2Client(createClient.Id, updateClient)
		require.NoError(t, err)
		assert.EqualValues(t, updateClient, *result)

		compareClient = updateClient
		compareClient.ClientSecret = ""
		result, _, err = c.GetOAuth2Client(updateClient.Id)
		require.NoError(t, err)
		assert.EqualValues(t, compareClient, *result)

		_, err = c.DeleteOAuth2Client(updateClient.Id)
		require.NoError(t, err)

		_, response, err := c.GetOAuth2Client(updateClient.Id)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})
}
