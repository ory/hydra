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
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package client_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestClient(prefix string) hydra.OAuth2Client {
	return hydra.OAuth2Client{
		ClientId:              "1234",
		ClientName:            prefix + "name",
		ClientSecret:          prefix + "secret",
		ClientUri:             prefix + "uri",
		Contacts:              []string{prefix + "peter", prefix + "pan"},
		GrantTypes:            []string{prefix + "client_credentials", prefix + "authorize_code"},
		LogoUri:               prefix + "logo",
		Owner:                 prefix + "an-owner",
		PolicyUri:             prefix + "policy-uri",
		Scope:                 prefix + "foo bar baz",
		TosUri:                prefix + "tos-uri",
		ResponseTypes:         []string{prefix + "id_token", prefix + "code"},
		RedirectUris:          []string{prefix + "redirect-url", prefix + "redirect-uri"},
		ClientSecretExpiresAt: 0,
	}
}

func TestClientSDK(t *testing.T) {
	manager := client.NewMemoryManager(nil)
	handler := &client.Handler{
		Manager:             manager,
		H:                   herodot.NewJSONWriter(nil),
		DefaultClientScopes: []string{"foo", "bar"},
	}

	router := httprouter.New()
	handler.SetRoutes(router)
	server := httptest.NewServer(router)
	c := hydra.NewOAuth2ApiWithBasePath(server.URL)
	t.Run("case=client default scopes are set", func(t *testing.T) {
		result, response, err := c.CreateOAuth2Client(hydra.OAuth2Client{
			ClientId: "scoped",
		})
		require.NoError(t, err)
		require.EqualValues(t, http.StatusCreated, response.StatusCode)
		assert.EqualValues(t, handler.DefaultClientScopes, strings.Split(result.Scope, " "))

		response, err = c.DeleteOAuth2Client("scoped")
		require.NoError(t, err)
		require.EqualValues(t, http.StatusNoContent, response.StatusCode)
	})

	t.Run("case=client is created and updated", func(t *testing.T) {
		createClient := createTestClient("")
		compareClient := createClient
		createClient.ClientSecretExpiresAt = 10

		// returned client is correct on Create
		result, _, err := c.CreateOAuth2Client(createClient)
		require.NoError(t, err)
		assert.EqualValues(t, compareClient, *result)

		// secret is not returned on GetOAuth2Client
		compareClient.ClientSecret = ""
		result, _, err = c.GetOAuth2Client(createClient.ClientId)
		assert.EqualValues(t, compareClient, *result)

		// listing clients returns the only added one
		results, _, err := c.ListOAuth2Clients(100, 0)
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.EqualValues(t, compareClient, results[0])

		// SecretExpiresAt gets overwritten with 0 on Update
		compareClient.ClientSecret = createClient.ClientSecret
		result, _, err = c.UpdateOAuth2Client(createClient.ClientId, createClient)
		require.NoError(t, err)
		assert.EqualValues(t, compareClient, *result)

		// create another client
		updateClient := createTestClient("foo")
		result, _, err = c.UpdateOAuth2Client(createClient.ClientId, updateClient)
		require.NoError(t, err)
		assert.EqualValues(t, updateClient, *result)

		// again, test if secret is not returned on Get
		compareClient = updateClient
		compareClient.ClientSecret = ""
		result, _, err = c.GetOAuth2Client(updateClient.ClientId)
		require.NoError(t, err)
		assert.EqualValues(t, compareClient, *result)

		// client can not be found after being deleted
		_, err = c.DeleteOAuth2Client(updateClient.ClientId)
		require.NoError(t, err)

		_, response, err := c.GetOAuth2Client(updateClient.ClientId)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})

	t.Run("case=public client is transmitted without secret", func(t *testing.T) {
		result, _, err := c.CreateOAuth2Client(hydra.OAuth2Client{
			Public: true,
		})

		require.NoError(t, err)
		assert.Equal(t, "", result.ClientSecret)

		result, _, err = c.CreateOAuth2Client(createTestClient(""))

		require.NoError(t, err)
		assert.NotEqual(t, "", result.ClientSecret)
	})
}
