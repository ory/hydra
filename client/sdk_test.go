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
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/mohae/deepcopy"

	"github.com/ory/x/pointerx"
	"github.com/ory/x/urlx"

	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/hydra/internal/httpclient/models"
	"github.com/ory/hydra/x"

	"github.com/ory/hydra/driver/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/internal"

	"github.com/ory/hydra/client"
	hydra "github.com/ory/hydra/internal/httpclient/client"
)

func createTestClient(prefix string) *models.OAuth2Client {
	return &models.OAuth2Client{
		ClientID:                  "1234",
		ClientName:                prefix + "name",
		ClientSecret:              prefix + "secret",
		ClientURI:                 prefix + "uri",
		Contacts:                  []string{prefix + "peter", prefix + "pan"},
		GrantTypes:                []string{prefix + "client_credentials", prefix + "authorize_code"},
		LogoURI:                   prefix + "logo",
		Owner:                     prefix + "an-owner",
		PolicyURI:                 prefix + "policy-uri",
		Scope:                     prefix + "foo bar baz",
		TosURI:                    prefix + "tos-uri",
		ResponseTypes:             []string{prefix + "id_token", prefix + "code"},
		RedirectUris:              []string{"https://" + prefix + "redirect-url", "https://" + prefix + "redirect-uri"},
		ClientSecretExpiresAt:     0,
		TokenEndpointAuthMethod:   "client_secret_basic",
		UserinfoSignedResponseAlg: "none",
		SubjectType:               "public",
		Metadata:                  map[string]interface{}{"foo": "bar"},
		// because these values are not nullable in the SQL schema, we have to set them not nil
		AllowedCorsOrigins: models.StringSlicePipeDelimiter{},
		Audience:           models.StringSlicePipeDelimiter{},
		Jwks:               models.JoseJSONWebKeySet(map[string]interface{}{}),
		// SectorIdentifierUri:   "https://sector.com/foo",
	}
}

func TestClientSDK(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	conf.MustSet(config.KeySubjectTypesSupported, []string{"public"})
	conf.MustSet(config.KeyDefaultClientScope, []string{"foo", "bar"})
	conf.MustSet(config.KeyPublicAllowDynamicRegistration, true)
	r := internal.NewRegistryMemory(t, conf)

	routerAdmin := x.NewRouterAdmin()
	routerPublic := x.NewRouterPublic()
	handler := client.NewHandler(r)
	handler.SetRoutes(routerAdmin, routerPublic)
	server := httptest.NewServer(routerAdmin)

	c := hydra.NewHTTPClientWithConfig(nil, &hydra.TransportConfig{Schemes: []string{"http"}, Host: urlx.ParseOrPanic(server.URL).Host})

	t.Run("case=client default scopes are set", func(t *testing.T) {
		result, err := c.Admin.CreateOAuth2Client(admin.NewCreateOAuth2ClientParams().WithBody(&models.OAuth2Client{
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
		assert.NotEmpty(t, result.Payload.RegistrationAccessToken)
		assert.NotEmpty(t, result.Payload.RegistrationClientURI)
		result.Payload.RegistrationAccessToken = ""
		result.Payload.RegistrationClientURI = ""

		assert.EqualValues(t, compareClient, result.Payload)
		assert.EqualValues(t, "bar", result.Payload.Metadata.(map[string]interface{})["foo"])

		// secret is not returned on GetOAuth2Client
		compareClient.ClientSecret = ""
		gresult, err := c.Admin.GetOAuth2Client(admin.NewGetOAuth2ClientParams().WithID(createClient.ClientID).WithContext(context.Background()))
		require.NoError(t, err)
		assert.NotEmpty(t, gresult.Payload.UpdatedAt)
		gresult.Payload.UpdatedAt = strfmt.DateTime{}
		assert.NotEmpty(t, gresult.Payload.CreatedAt)
		gresult.Payload.CreatedAt = strfmt.DateTime{}
		assert.EqualValues(t, compareClient, gresult.Payload)

		// get client will return The request could not be authorized
		gresult, err = c.Admin.GetOAuth2Client(admin.NewGetOAuth2ClientParams().WithID("unknown"))
		require.Error(t, err)
		assert.Empty(t, gresult)
		assert.True(t, strings.Contains(err.Error(), "404"), err.Error())

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
		compareClient.ClientSecret = createClient.ClientSecret
		uresult, err := c.Admin.UpdateOAuth2Client(admin.NewUpdateOAuth2ClientParams().WithID(createClient.ClientID).WithBody(createClient))
		require.NoError(t, err)
		assert.NotEmpty(t, uresult.Payload.UpdatedAt)
		uresult.Payload.UpdatedAt = strfmt.DateTime{}
		assert.NotEmpty(t, uresult.Payload.CreatedAt)
		uresult.Payload.CreatedAt = strfmt.DateTime{}
		assert.EqualValues(t, compareClient, uresult.Payload)

		// create another client
		updateClient := createTestClient("foo")
		uresult, err = c.Admin.UpdateOAuth2Client(admin.NewUpdateOAuth2ClientParams().WithID(createClient.ClientID).WithBody(updateClient))
		require.NoError(t, err)
		assert.NotEmpty(t, uresult.Payload.UpdatedAt)
		uresult.Payload.UpdatedAt = strfmt.DateTime{}
		assert.NotEmpty(t, uresult.Payload.CreatedAt)
		uresult.Payload.CreatedAt = strfmt.DateTime{}
		assert.EqualValues(t, updateClient, uresult.Payload)

		// again, test if secret is not returned on Get
		compareClient = updateClient
		compareClient.ClientSecret = ""
		gresult, err = c.Admin.GetOAuth2Client(admin.NewGetOAuth2ClientParams().WithID(updateClient.ClientID))
		require.NoError(t, err)
		assert.NotEmpty(t, gresult.Payload.UpdatedAt)
		gresult.Payload.UpdatedAt = strfmt.DateTime{}
		assert.NotEmpty(t, gresult.Payload.CreatedAt)
		gresult.Payload.CreatedAt = strfmt.DateTime{}
		assert.EqualValues(t, compareClient, gresult.Payload)

		// client can not be found after being deleted
		_, err = c.Admin.DeleteOAuth2Client(admin.NewDeleteOAuth2ClientParams().WithID(updateClient.ClientID))
		require.NoError(t, err)

		_, err = c.Admin.GetOAuth2Client(admin.NewGetOAuth2ClientParams().WithID(updateClient.ClientID))
		require.Error(t, err)
	})

	t.Run("case=public client is transmitted without secret", func(t *testing.T) {
		result, err := c.Admin.CreateOAuth2Client(admin.NewCreateOAuth2ClientParams().WithBody(&models.OAuth2Client{
			TokenEndpointAuthMethod: "none",
		}))
		require.NoError(t, err)

		assert.Equal(t, "", result.Payload.ClientSecret)

		result, err = c.Admin.CreateOAuth2Client(admin.NewCreateOAuth2ClientParams().WithBody(createTestClient("")))
		require.NoError(t, err)

		assert.Equal(t, "secret", result.Payload.ClientSecret)
	})

	t.Run("case=id should be set properly", func(t *testing.T) {
		for k, tc := range []struct {
			client   *models.OAuth2Client
			expectID string
		}{
			{
				client: &models.OAuth2Client{},
			},
			{
				client:   &models.OAuth2Client{ClientID: "set-properly-1"},
				expectID: "set-properly-1",
			},
			{
				client:   &models.OAuth2Client{ClientID: "set-properly-2"},
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
	t.Run("case=patch client legally", func(t *testing.T) {
		op := "add"
		path := "/redirect_uris/-"
		value := "http://foo.bar"

		client := createTestClient("")
		client.ClientID = "patch1_client"
		_, err := c.Admin.CreateOAuth2Client(admin.NewCreateOAuth2ClientParams().WithBody(client))
		require.NoError(t, err)

		expected := deepcopy.Copy(client).(*models.OAuth2Client)
		expected.RedirectUris = append(expected.RedirectUris, value)

		result, err := c.Admin.PatchOAuth2Client(admin.NewPatchOAuth2ClientParams().WithID(client.ClientID).WithBody(models.PatchRequest{{Op: &op, Path: &path, Value: value}}))
		require.NoError(t, err)
		expected.CreatedAt = result.Payload.CreatedAt
		expected.UpdatedAt = result.Payload.UpdatedAt
		expected.ClientSecret = result.Payload.ClientSecret
		expected.ClientSecretExpiresAt = result.Payload.ClientSecretExpiresAt
		require.Equal(t, expected, result.Payload)
	})

	t.Run("case=patch client illegally", func(t *testing.T) {
		op := "replace"
		path := "/id"
		value := "foo"

		client := createTestClient("")
		client.ClientID = "patch2_client"
		_, err := c.Admin.CreateOAuth2Client(admin.NewCreateOAuth2ClientParams().WithBody(client))
		require.NoError(t, err)

		_, err = c.Admin.PatchOAuth2Client(admin.NewPatchOAuth2ClientParams().WithID(client.ClientID).WithBody(models.PatchRequest{{Op: &op, Path: &path, Value: value}}))
		require.Error(t, err)
	})

	t.Run("case=patch should not alter secret if not requested", func(t *testing.T) {
		op := "replace"
		path := "/client_uri"
		value := "http://foo.bar"

		client := createTestClient("")
		client.ClientID = "patch3_client"
		_, err := c.Admin.CreateOAuth2Client(admin.NewCreateOAuth2ClientParams().WithBody(client))
		require.NoError(t, err)

		result1, err := c.Admin.PatchOAuth2Client(admin.NewPatchOAuth2ClientParams().WithID(client.ClientID).WithBody(models.PatchRequest{{Op: &op, Path: &path, Value: value}}))
		require.NoError(t, err)
		result2, err := c.Admin.PatchOAuth2Client(admin.NewPatchOAuth2ClientParams().WithID(client.ClientID).WithBody(models.PatchRequest{{Op: &op, Path: &path, Value: value}}))
		require.NoError(t, err)

		// secret hashes shouldn't change between these PUT calls
		require.Equal(t, result1.Payload.ClientSecret, result2.Payload.ClientSecret)
	})
}
