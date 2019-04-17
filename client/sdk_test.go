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
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-openapi/strfmt"

	"github.com/ory/hydra/sdk/go/hydra/client/admin"
	"github.com/ory/hydra/sdk/go/hydra/models"
	"github.com/ory/hydra/x"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/urlx"

	"github.com/spf13/viper"

	"github.com/ory/hydra/driver/configuration"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/internal"

	"github.com/ory/hydra/client"
	hydra "github.com/ory/hydra/sdk/go/hydra/client"
)

func createTestClient(prefix string) *models.Client {
	return &models.Client{
		ClientID:                  "1234",
		Name:                      prefix + "name",
		Secret:                    prefix + "secret",
		ClientURI:                 prefix + "uri",
		Contacts:                  []string{prefix + "peter", prefix + "pan"},
		GrantTypes:                []string{prefix + "client_credentials", prefix + "authorize_code"},
		LogoURI:                   prefix + "logo",
		Owner:                     prefix + "an-owner",
		PolicyURI:                 prefix + "policy-uri",
		Scope:                     prefix + "foo bar baz",
		TermsOfServiceURI:         prefix + "tos-uri",
		ResponseTypes:             []string{prefix + "id_token", prefix + "code"},
		RedirectUris:              []string{"https://" + prefix + "redirect-url", "https://" + prefix + "redirect-uri"},
		SecretExpiresAt:           0,
		TokenEndpointAuthMethod:   "client_secret_basic",
		UserinfoSignedResponseAlg: "none",
		SubjectType:               "public",
		//SectorIdentifierUri:   "https://sector.com/foo",
	}
}

func TestClientSDK(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	viper.Set(configuration.ViperKeySubjectTypesSupported, []string{"public"})
	viper.Set(configuration.ViperKeyDefaultClientScope, []string{"foo", "bar"})
	r := internal.NewRegistry(conf)

	router := x.NewRouterAdmin()
	handler := client.NewHandler(r)
	handler.SetRoutes(router)
	server := httptest.NewServer(router)

	c := hydra.NewHTTPClientWithConfig(nil, &hydra.TransportConfig{Schemes: []string{"http"}, Host: urlx.ParseOrPanic(server.URL).Host})

	t.Run("case=client default scopes are set", func(t *testing.T) {
		result, err := c.Admin.CreateOAuth2Client(admin.NewCreateOAuth2ClientParams().WithBody(&models.Client{
			ClientID: "scoped",
		}))
		require.NoError(t, err)
		assert.EqualValues(t, conf.DefaultClientScope(), strings.Split(result.Payload.Scope, " "))

		_, err = c.Admin.DeleteOAuth2Client(admin.NewDeleteOAuth2ClientParams().WithID("scoped"))
		require.NoError(t, err)
	})

	t.Run("case=client is created and updated", func(t *testing.T) {
		createClient := createTestClient("")
		compareClient := createClient
		// This is not yet supported:
		// 		createClient.SecretExpiresAt = 10

		// returned client is correct on Create
		result, err := c.Admin.CreateOAuth2Client(admin.NewCreateOAuth2ClientParams().WithBody(createClient))
		require.NoError(t, err)
		assert.NotEmpty(t, result.Payload.UpdatedAt)
		result.Payload.UpdatedAt = strfmt.DateTime{}
		assert.NotEmpty(t, result.Payload.CreatedAt)
		result.Payload.CreatedAt = strfmt.DateTime{}
		assert.EqualValues(t, compareClient, result.Payload)

		// secret is not returned on GetOAuth2Client
		compareClient.Secret = ""
		gresult, err := c.Admin.GetOAuth2Client(admin.NewGetOAuth2ClientParams().WithID(createClient.ClientID))
		require.NoError(t, err)
		assert.NotEmpty(t, gresult.Payload.UpdatedAt)
		gresult.Payload.UpdatedAt = strfmt.DateTime{}
		assert.NotEmpty(t, gresult.Payload.CreatedAt)
		gresult.Payload.CreatedAt = strfmt.DateTime{}
		assert.EqualValues(t, compareClient, gresult.Payload)

		// listing clients returns the only added one
		results, err := c.Admin.ListOAuth2Clients(admin.NewListOAuth2ClientsParams().WithLimit(pointerx.Int64(100)))
		require.NoError(t, err)
		assert.Len(t, results.Payload, 1)
		assert.NotEmpty(t, results.Payload[0].UpdatedAt)
		results.Payload[0].UpdatedAt = strfmt.DateTime{}
		assert.NotEmpty(t, results.Payload[0].CreatedAt)
		results.Payload[0].CreatedAt = strfmt.DateTime{}
		assert.EqualValues(t, compareClient, results.Payload[0])

		// SecretExpiresAt gets overwritten with 0 on Update
		compareClient.Secret = createClient.Secret
		uresult, err := c.Admin.UpdateOAuth2Client(admin.NewUpdateOAuth2ClientParams().WithID(createClient.ClientID).WithBody(createClient))
		require.NoError(t, err)
		assert.NotEmpty(t, uresult.Payload.UpdatedAt)
		uresult.Payload.UpdatedAt = strfmt.DateTime{}
		assert.EqualValues(t, compareClient, uresult.Payload)

		// create another client
		updateClient := createTestClient("foo")
		uresult, err = c.Admin.UpdateOAuth2Client(admin.NewUpdateOAuth2ClientParams().WithID(createClient.ClientID).WithBody(updateClient))
		require.NoError(t, err)
		assert.NotEmpty(t, uresult.Payload.UpdatedAt)
		uresult.Payload.UpdatedAt = strfmt.DateTime{}
		assert.EqualValues(t, updateClient, uresult.Payload)

		// again, test if secret is not returned on Get
		compareClient = updateClient
		compareClient.Secret = ""
		gresult, err = c.Admin.GetOAuth2Client(admin.NewGetOAuth2ClientParams().WithID(updateClient.ClientID))
		require.NoError(t, err)
		assert.NotEmpty(t, gresult.Payload.UpdatedAt)
		gresult.Payload.UpdatedAt = strfmt.DateTime{}
		assert.EqualValues(t, compareClient, gresult.Payload)

		// client can not be found after being deleted
		_, err = c.Admin.DeleteOAuth2Client(admin.NewDeleteOAuth2ClientParams().WithID(updateClient.ClientID))
		require.NoError(t, err)

		_, err = c.Admin.GetOAuth2Client(admin.NewGetOAuth2ClientParams().WithID(updateClient.ClientID))
		require.Error(t, err)
	})

	t.Run("case=public client is transmitted without secret", func(t *testing.T) {
		result, err := c.Admin.CreateOAuth2Client(admin.NewCreateOAuth2ClientParams().WithBody(&models.Client{
			TokenEndpointAuthMethod: "none",
		}))
		require.NoError(t, err)

		assert.Equal(t, "", result.Payload.Secret)

		result, err = c.Admin.CreateOAuth2Client(admin.NewCreateOAuth2ClientParams().WithBody(createTestClient("")))
		require.NoError(t, err)

		assert.Equal(t, "secret", result.Payload.Secret)
	})

	t.Run("case=id should be set properly", func(t *testing.T) {
		for k, tc := range []struct {
			client   *models.Client
			expectID string
		}{
			{
				client: &models.Client{},
			},
			{
				client:   &models.Client{ClientID: "set-properly-1"},
				expectID: "set-properly-1",
			},
			{
				client:   &models.Client{ClientID: "set-properly-2"},
				expectID: "set-properly-2",
			},
		} {
			t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
				result, err := c.Admin.CreateOAuth2Client(admin.NewCreateOAuth2ClientParams().WithBody(tc.client))
				require.NoError(t, err)

				assert.NotEmpty(t, result.Payload.ClientID)

				id := result.Payload.ClientID
				if tc.expectID != "" {
					assert.EqualValues(t, tc.expectID, result.Payload.ClientID)
					id = tc.expectID
				}

				gresult, err := c.Admin.GetOAuth2Client(admin.NewGetOAuth2ClientParams().WithID(id))
				require.NoError(t, err)

				assert.EqualValues(t, id, gresult.Payload.ClientID)
			})
		}
	})
}
