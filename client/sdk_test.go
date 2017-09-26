package client_test

import (
	"net/http/httptest"
	"github.com/ory/fosite"
	"github.com/ory/ladon"
	"github.com/ory/herodot"
	"github.com/julienschmidt/httprouter"
	"testing"
	"github.com/ory/hydra/compose"
	"github.com/ory/hydra/client"
	hydra "github.com/ory/hydra/sdk/go/swagger"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"github.com/stretchr/testify/assert"
	"net/http"
)

func createTestClient(prefix string) hydra.OauthClient {
	return hydra.OauthClient{
		Id:            "1234",
		ClientName:    prefix + "name",
		ClientSecret:   prefix + "secret",
		ClientUri:      prefix + "uri",
		Contacts:      []string{ prefix + "peter",  prefix + "pan"},
		GrantTypes:    []string{ prefix + "client_credentials",  prefix + "authorize_code"},
		LogoUri:        prefix + "logo",
		Owner:         prefix +  "an-owner",
		PolicyUri:     prefix +  "policy-uri",
		Scope:         prefix +  "foo bar baz",
		TosUri:        prefix +  "tos-uri",
		ResponseTypes: []string{ prefix + "id_token", prefix +  "code"},
		RedirectUris:  []string{ prefix + "redirect-url", prefix +  "redirect-uri"},
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
	c := hydra.NewClientsApiWithBasePath(server.URL)
	token, err := httpClient.Transport.(*oauth2.Transport).Source.Token()
	require.NoError(t, err)
	c.Configuration.AccessToken = token.AccessToken

	t.Run("foo", func(t *testing.T) {
		createClient := createTestClient("")

		result, _, err := c.CreateOAuthClient(createClient)
		require.NoError(t, err)
		assert.EqualValues(t, createClient, *result)

		compareClient := createClient
		compareClient.ClientSecret = ""
		result, _ , err = c.GetOAuthClient(createClient.Id)
		assert.EqualValues(t, compareClient, *result)

		//results, _, err := c.ListOAuthClients()

		updateClient :=  createTestClient("foo")
		result, _ , err = c.UpdateOAuthClient(createClient.Id, updateClient)
		require.NoError(t, err)
		assert.EqualValues(t, updateClient, *result)

		compareClient = updateClient
		compareClient.ClientSecret = ""
		result, _ , err = c.GetOAuthClient(updateClient.Id)
		require.NoError(t, err)
		assert.EqualValues(t, compareClient, *result)

		_, err = c.DeleteOAuthClient(updateClient.Id)
		require.NoError(t, err)

		_, response, err := c.GetOAuthClient(updateClient.Id)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})
}
