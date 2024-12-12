// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ory/hydra/v2/internal/testhelpers"

	"github.com/ory/x/assertx"

	"github.com/ory/x/ioutilx"

	"github.com/ory/x/uuidx"

	"github.com/mohae/deepcopy"

	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/contextx"
	"github.com/ory/x/pointerx"

	"github.com/ory/hydra/v2/driver/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/client"
)

func createTestClient(prefix string) hydra.OAuth2Client {
	return hydra.OAuth2Client{
		ClientName:                pointerx.Ptr(prefix + "name"),
		ClientSecret:              pointerx.Ptr(prefix + "secret"),
		ClientUri:                 pointerx.Ptr(prefix + "uri"),
		Contacts:                  []string{prefix + "peter", prefix + "pan"},
		GrantTypes:                []string{prefix + "client_credentials", prefix + "authorize_code"},
		LogoUri:                   pointerx.Ptr(prefix + "logo"),
		Owner:                     pointerx.Ptr(prefix + "an-owner"),
		PolicyUri:                 pointerx.Ptr(prefix + "policy-uri"),
		Scope:                     pointerx.Ptr(prefix + "foo bar baz"),
		TosUri:                    pointerx.Ptr(prefix + "tos-uri"),
		ResponseTypes:             []string{prefix + "id_token", prefix + "code"},
		RedirectUris:              []string{"https://" + prefix + "redirect-url", "https://" + prefix + "redirect-uri"},
		ClientSecretExpiresAt:     pointerx.Ptr[int64](0),
		TokenEndpointAuthMethod:   pointerx.Ptr("client_secret_basic"),
		UserinfoSignedResponseAlg: pointerx.Ptr("none"),
		SubjectType:               pointerx.Ptr("public"),
		Metadata:                  map[string]interface{}{"foo": "bar"},
		// because these values are not nullable in the SQL schema, we have to set them not nil
		AllowedCorsOrigins: []string{},
		Audience:           []string{},
		Jwks:               map[string]interface{}{},
		SkipConsent:        pointerx.Ptr(false),
	}
}

var defaultIgnoreFields = []string{"client_id", "registration_access_token", "registration_client_uri", "created_at", "updated_at"}

func TestClientSDK(t *testing.T) {
	ctx := context.Background()
	conf := testhelpers.NewConfigurationWithDefaults()
	conf.MustSet(ctx, config.KeySubjectTypesSupported, []string{"public"})
	conf.MustSet(ctx, config.KeyDefaultClientScope, []string{"foo", "bar"})
	conf.MustSet(ctx, config.KeyPublicAllowDynamicRegistration, true)
	r := testhelpers.NewRegistryMemory(t, conf, &contextx.Static{C: conf.Source(ctx)})

	routerAdmin := x.NewRouterAdmin(conf.AdminURL)
	routerPublic := x.NewRouterPublic()
	handler := client.NewHandler(r)
	handler.SetRoutes(routerAdmin, routerPublic)
	server := httptest.NewServer(routerAdmin)
	conf.MustSet(ctx, config.KeyAdminURL, server.URL)

	c := hydra.NewAPIClient(hydra.NewConfiguration())
	c.GetConfig().Servers = hydra.ServerConfigurations{{URL: server.URL}}

	t.Run("case=client default scopes are set", func(t *testing.T) {
		result, _, err := c.OAuth2API.CreateOAuth2Client(ctx).OAuth2Client(hydra.OAuth2Client{}).Execute()
		require.NoError(t, err)
		assert.EqualValues(t, conf.DefaultClientScope(ctx), strings.Split(*result.Scope, " "))

		_, err = c.OAuth2API.DeleteOAuth2Client(ctx, *result.ClientId).Execute()
		require.NoError(t, err)
	})

	t.Run("case=client is created and updated", func(t *testing.T) {
		createClient := createTestClient("")
		compareClient := createClient
		// This is not yet supported:
		// 		createClient.SecretExpiresAt = 10

		// returned client is correct on Create
		result, _, err := c.OAuth2API.CreateOAuth2Client(ctx).OAuth2Client(createClient).Execute()
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
		compareClient.ClientSecret = pointerx.Ptr("")
		gresult, _, err := c.OAuth2API.GetOAuth2Client(context.Background(), *createClient.ClientId).Execute()
		require.NoError(t, err)
		assertx.EqualAsJSONExcept(t, compareClient, gresult, append(defaultIgnoreFields, "client_secret"))

		// get client will return The request could not be authorized
		gresult, _, err = c.OAuth2API.GetOAuth2Client(context.Background(), "unknown").Execute()
		require.Error(t, err)
		assert.Empty(t, gresult)
		assert.True(t, strings.Contains(err.Error(), "404"), err.Error())

		// listing clients returns the only added one
		results, _, err := c.OAuth2API.ListOAuth2Clients(context.Background()).PageSize(100).Execute()
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assertx.EqualAsJSONExcept(t, compareClient, results[0], append(defaultIgnoreFields, "client_secret"))

		// SecretExpiresAt gets overwritten with 0 on Update
		compareClient.ClientSecret = createClient.ClientSecret
		uresult, _, err := c.OAuth2API.SetOAuth2Client(context.Background(), *createClient.ClientId).OAuth2Client(createClient).Execute()
		require.NoError(t, err)
		assertx.EqualAsJSONExcept(t, compareClient, uresult, append(defaultIgnoreFields, "client_secret"))

		// create another client
		updateClient := createTestClient("foo")
		uresult, _, err = c.OAuth2API.SetOAuth2Client(context.Background(), *createClient.ClientId).OAuth2Client(updateClient).Execute()
		require.NoError(t, err)
		assert.NotEqual(t, updateClient.ClientId, uresult.ClientId)
		updateClient.ClientId = uresult.ClientId
		assertx.EqualAsJSONExcept(t, updateClient, uresult, append(defaultIgnoreFields, "client_secret"))

		// again, test if secret is not returned on Get
		compareClient = updateClient
		compareClient.ClientSecret = pointerx.Ptr("")
		gresult, _, err = c.OAuth2API.GetOAuth2Client(context.Background(), *updateClient.ClientId).Execute()
		require.NoError(t, err)
		assertx.EqualAsJSONExcept(t, compareClient, gresult, append(defaultIgnoreFields, "client_secret"))

		// client can not be found after being deleted
		_, err = c.OAuth2API.DeleteOAuth2Client(context.Background(), *updateClient.ClientId).Execute()
		require.NoError(t, err)

		_, _, err = c.OAuth2API.GetOAuth2Client(context.Background(), *updateClient.ClientId).Execute()
		require.Error(t, err)
	})

	t.Run("case=public client is transmitted without secret", func(t *testing.T) {
		result, _, err := c.OAuth2API.CreateOAuth2Client(context.Background()).OAuth2Client(hydra.OAuth2Client{
			TokenEndpointAuthMethod: pointerx.Ptr("none"),
		}).Execute()
		require.NoError(t, err)

		assert.Equal(t, "", pointerx.Deref(result.ClientSecret))

		result, _, err = c.OAuth2API.CreateOAuth2Client(context.Background()).OAuth2Client(createTestClient("")).Execute()
		require.NoError(t, err)

		assert.Equal(t, "secret", pointerx.Deref(result.ClientSecret))
	})

	t.Run("case=id can be set", func(t *testing.T) {
		id := uuidx.NewV4().String()
		result, _, err := c.OAuth2API.CreateOAuth2Client(context.Background()).OAuth2Client(hydra.OAuth2Client{ClientId: pointerx.Ptr(id)}).Execute()
		require.NoError(t, err)

		assert.Equal(t, id, pointerx.Deref(result.ClientId))
	})

	t.Run("case=patch client legally", func(t *testing.T) {
		op := "add"
		path := "/redirect_uris/-"
		value := "http://foo.bar"

		cl := createTestClient("")
		created, _, err := c.OAuth2API.CreateOAuth2Client(context.Background()).OAuth2Client(cl).Execute()
		require.NoError(t, err)
		cl.ClientId = created.ClientId

		expected := deepcopy.Copy(cl).(hydra.OAuth2Client)
		expected.RedirectUris = append(expected.RedirectUris, value)

		result, _, err := c.OAuth2API.PatchOAuth2Client(context.Background(), *cl.ClientId).JsonPatch([]hydra.JsonPatch{{Op: op, Path: path, Value: value}}).Execute()
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
		created, res, err := c.OAuth2API.CreateOAuth2Client(context.Background()).OAuth2Client(client).Execute()
		require.NoError(t, err, "%s", ioutilx.MustReadAll(res.Body))
		client.ClientId = created.ClientId

		_, _, err = c.OAuth2API.PatchOAuth2Client(context.Background(), *client.ClientId).JsonPatch([]hydra.JsonPatch{{Op: op, Path: path, Value: value}}).Execute()
		require.Error(t, err)
	})

	t.Run("case=patch should not alter secret if not requested", func(t *testing.T) {
		op := "replace"
		path := "/client_uri"
		value := "http://foo.bar"

		client := createTestClient("")
		created, _, err := c.OAuth2API.CreateOAuth2Client(context.Background()).OAuth2Client(client).Execute()
		require.NoError(t, err)
		client.ClientId = created.ClientId

		result1, _, err := c.OAuth2API.PatchOAuth2Client(context.Background(), *client.ClientId).JsonPatch([]hydra.JsonPatch{{Op: op, Path: path, Value: value}}).Execute()
		require.NoError(t, err)
		result2, _, err := c.OAuth2API.PatchOAuth2Client(context.Background(), *client.ClientId).JsonPatch([]hydra.JsonPatch{{Op: op, Path: path, Value: value}}).Execute()
		require.NoError(t, err)

		// secret hashes shouldn't change between these PUT calls
		require.Equal(t, result1.ClientSecret, result2.ClientSecret)
	})
}
