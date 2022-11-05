// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ory/x/assertx"

	"github.com/ory/x/ioutilx"

	"github.com/ory/x/snapshotx"

	"github.com/ory/x/uuidx"

	"github.com/mohae/deepcopy"

	"github.com/ory/hydra/x"
	"github.com/ory/x/contextx"
	"github.com/ory/x/pointerx"

	"github.com/ory/hydra/driver/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/internal"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/client"
)

func createTestClient(prefix string) hydra.OAuth2Client {
	return hydra.OAuth2Client{
		ClientName:                pointerx.String(prefix + "name"),
		ClientSecret:              pointerx.String(prefix + "secret"),
		ClientUri:                 pointerx.String(prefix + "uri"),
		Contacts:                  []string{prefix + "peter", prefix + "pan"},
		GrantTypes:                []string{prefix + "client_credentials", prefix + "authorize_code"},
		LogoUri:                   pointerx.String(prefix + "logo"),
		Owner:                     pointerx.String(prefix + "an-owner"),
		PolicyUri:                 pointerx.String(prefix + "policy-uri"),
		Scope:                     pointerx.String(prefix + "foo bar baz"),
		TosUri:                    pointerx.String(prefix + "tos-uri"),
		ResponseTypes:             []string{prefix + "id_token", prefix + "code"},
		RedirectUris:              []string{"https://" + prefix + "redirect-url", "https://" + prefix + "redirect-uri"},
		ClientSecretExpiresAt:     pointerx.Int64(0),
		TokenEndpointAuthMethod:   pointerx.String("client_secret_basic"),
		UserinfoSignedResponseAlg: pointerx.String("none"),
		SubjectType:               pointerx.String("public"),
		Metadata:                  map[string]interface{}{"foo": "bar"},
		// because these values are not nullable in the SQL schema, we have to set them not nil
		AllowedCorsOrigins: []string{},
		Audience:           []string{},
		Jwks:               map[string]interface{}{},
	}
}

var defaultIgnoreFields = []string{"client_id", "registration_access_token", "registration_client_uri", "created_at", "updated_at"}

func TestClientSDK(t *testing.T) {
	ctx := context.Background()
	conf := internal.NewConfigurationWithDefaults()
	conf.MustSet(ctx, config.KeySubjectTypesSupported, []string{"public"})
	conf.MustSet(ctx, config.KeyDefaultClientScope, []string{"foo", "bar"})
	conf.MustSet(ctx, config.KeyPublicAllowDynamicRegistration, true)
	r := internal.NewRegistryMemory(t, conf, &contextx.Static{C: conf.Source(ctx)})

	routerAdmin := x.NewRouterAdmin(conf.AdminURL)
	routerPublic := x.NewRouterPublic()
	handler := client.NewHandler(r)
	handler.SetRoutes(routerAdmin, routerPublic)
	server := httptest.NewServer(routerAdmin)
	conf.MustSet(ctx, config.KeyAdminURL, server.URL)

	c := hydra.NewAPIClient(hydra.NewConfiguration())
	c.GetConfig().Servers = hydra.ServerConfigurations{{URL: server.URL}}

	t.Run("case=client default scopes are set", func(t *testing.T) {
		result, _, err := c.OAuth2Api.CreateOAuth2Client(ctx).OAuth2Client(hydra.OAuth2Client{}).Execute()
		require.NoError(t, err)
		assert.EqualValues(t, conf.DefaultClientScope(ctx), strings.Split(*result.Scope, " "))

		_, err = c.OAuth2Api.DeleteOAuth2Client(ctx, *result.ClientId).Execute()
		require.NoError(t, err)
	})

	t.Run("case=client is created and updated", func(t *testing.T) {
		createClient := createTestClient("")
		compareClient := createClient
		// This is not yet supported:
		// 		createClient.SecretExpiresAt = 10

		// returned client is correct on Create
		result, _, err := c.OAuth2Api.CreateOAuth2Client(ctx).OAuth2Client(createClient).Execute()
		require.NoError(t, err)
		assert.NotEmpty(t, result.UpdatedAt)
		assert.NotEmpty(t, result.CreatedAt)
		assert.NotEmpty(t, result.RegistrationAccessToken)
		assert.NotEmpty(t, result.RegistrationClientUri)
		assert.NotEmpty(t, result.ClientId)
		createClient.ClientId = result.ClientId

		assertx.EqualAsJSONExcept(t, compareClient, result, defaultIgnoreFields)
		assert.EqualValues(t, "bar", result.Metadata.(map[string]interface{})["foo"])

		// secret is not returned on GetOAuth2Client
		compareClient.ClientSecret = x.ToPointer("")
		gresult, _, err := c.OAuth2Api.GetOAuth2Client(context.Background(), *createClient.ClientId).Execute()
		require.NoError(t, err)
		assertx.EqualAsJSONExcept(t, compareClient, gresult, append(defaultIgnoreFields, "client_secret"))

		// get client will return The request could not be authorized
		gresult, _, err = c.OAuth2Api.GetOAuth2Client(context.Background(), "unknown").Execute()
		require.Error(t, err)
		assert.Empty(t, gresult)
		assert.True(t, strings.Contains(err.Error(), "404"), err.Error())

		// listing clients returns the only added one
		results, _, err := c.OAuth2Api.ListOAuth2Clients(context.Background()).PageSize(100).Execute()
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assertx.EqualAsJSONExcept(t, compareClient, results[0], append(defaultIgnoreFields, "client_secret"))

		// SecretExpiresAt gets overwritten with 0 on Update
		compareClient.ClientSecret = createClient.ClientSecret
		uresult, _, err := c.OAuth2Api.SetOAuth2Client(context.Background(), *createClient.ClientId).OAuth2Client(createClient).Execute()
		require.NoError(t, err)
		assertx.EqualAsJSONExcept(t, compareClient, uresult, append(defaultIgnoreFields, "client_secret"))

		// create another client
		updateClient := createTestClient("foo")
		uresult, _, err = c.OAuth2Api.SetOAuth2Client(context.Background(), *createClient.ClientId).OAuth2Client(updateClient).Execute()
		require.NoError(t, err)
		assert.NotEqual(t, updateClient.ClientId, uresult.ClientId)
		updateClient.ClientId = uresult.ClientId
		assertx.EqualAsJSONExcept(t, updateClient, uresult, append(defaultIgnoreFields, "client_secret"))

		// again, test if secret is not returned on Get
		compareClient = updateClient
		compareClient.ClientSecret = x.ToPointer("")
		gresult, _, err = c.OAuth2Api.GetOAuth2Client(context.Background(), *updateClient.ClientId).Execute()
		require.NoError(t, err)
		assertx.EqualAsJSONExcept(t, compareClient, gresult, append(defaultIgnoreFields, "client_secret"))

		// client can not be found after being deleted
		_, err = c.OAuth2Api.DeleteOAuth2Client(context.Background(), *updateClient.ClientId).Execute()
		require.NoError(t, err)

		_, _, err = c.OAuth2Api.GetOAuth2Client(context.Background(), *updateClient.ClientId).Execute()
		require.Error(t, err)
	})

	t.Run("case=public client is transmitted without secret", func(t *testing.T) {
		result, _, err := c.OAuth2Api.CreateOAuth2Client(context.Background()).OAuth2Client(hydra.OAuth2Client{
			TokenEndpointAuthMethod: x.ToPointer("none"),
		}).Execute()
		require.NoError(t, err)

		assert.Equal(t, "", x.FromPointer[string](result.ClientSecret))

		result, _, err = c.OAuth2Api.CreateOAuth2Client(context.Background()).OAuth2Client(createTestClient("")).Execute()
		require.NoError(t, err)

		assert.Equal(t, "secret", x.FromPointer[string](result.ClientSecret))
	})

	t.Run("case=id can not be set", func(t *testing.T) {
		_, res, err := c.OAuth2Api.CreateOAuth2Client(context.Background()).OAuth2Client(hydra.OAuth2Client{ClientId: x.ToPointer(uuidx.NewV4().String())}).Execute()
		require.Error(t, err)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		snapshotx.SnapshotT(t, json.RawMessage(body))
	})

	t.Run("case=patch client legally", func(t *testing.T) {
		op := "add"
		path := "/redirect_uris/-"
		value := "http://foo.bar"

		client := createTestClient("")
		created, _, err := c.OAuth2Api.CreateOAuth2Client(context.Background()).OAuth2Client(client).Execute()
		require.NoError(t, err)
		client.ClientId = created.ClientId

		expected := deepcopy.Copy(client).(hydra.OAuth2Client)
		expected.RedirectUris = append(expected.RedirectUris, value)

		result, _, err := c.OAuth2Api.PatchOAuth2Client(context.Background(), *client.ClientId).JsonPatch([]hydra.JsonPatch{{Op: op, Path: path, Value: value}}).Execute()
		require.NoError(t, err)
		expected.CreatedAt = result.CreatedAt
		expected.UpdatedAt = result.UpdatedAt
		expected.ClientSecret = result.ClientSecret
		expected.ClientSecretExpiresAt = result.ClientSecretExpiresAt
		assertx.EqualAsJSONExcept(t, expected, result, nil)
	})

	t.Run("case=patch client illegally", func(t *testing.T) {
		op := "replace"
		path := "/id"
		value := "foo"

		client := createTestClient("")
		created, res, err := c.OAuth2Api.CreateOAuth2Client(context.Background()).OAuth2Client(client).Execute()
		require.NoError(t, err, "%s", ioutilx.MustReadAll(res.Body))
		client.ClientId = created.ClientId

		_, _, err = c.OAuth2Api.PatchOAuth2Client(context.Background(), *client.ClientId).JsonPatch([]hydra.JsonPatch{{Op: op, Path: path, Value: value}}).Execute()
		require.Error(t, err)
	})

	t.Run("case=patch should not alter secret if not requested", func(t *testing.T) {
		op := "replace"
		path := "/client_uri"
		value := "http://foo.bar"

		client := createTestClient("")
		created, _, err := c.OAuth2Api.CreateOAuth2Client(context.Background()).OAuth2Client(client).Execute()
		require.NoError(t, err)
		client.ClientId = created.ClientId

		result1, _, err := c.OAuth2Api.PatchOAuth2Client(context.Background(), *client.ClientId).JsonPatch([]hydra.JsonPatch{{Op: op, Path: path, Value: value}}).Execute()
		require.NoError(t, err)
		result2, _, err := c.OAuth2Api.PatchOAuth2Client(context.Background(), *client.ClientId).JsonPatch([]hydra.JsonPatch{{Op: op, Path: path, Value: value}}).Execute()
		require.NoError(t, err)

		// secret hashes shouldn't change between these PUT calls
		require.Equal(t, result1.ClientSecret, result2.ClientSecret)
	})
}
