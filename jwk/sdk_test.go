package jwk_test

import (
	"net/http/httptest"
	"testing"

	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/compose"
	. "github.com/ory/hydra/jwk"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/ory/ladon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWKSDK(t *testing.T) {
	localWarden, httpClient := compose.NewMockFirewall(
		"tests",
		"alice",
		fosite.Arguments{
			"hydra.keys.create",
			"hydra.keys.get",
			"hydra.keys.delete",
			"hydra.keys.update",
		}, &ladon.DefaultPolicy{
			ID:        "1",
			Subjects:  []string{"alice"},
			Resources: []string{"rn:hydra:keys:<.*>"},
			Actions:   []string{"create", "get", "delete", "update"},
			Effect:    ladon.AllowAccess,
		},
	)

	manager := new(MemoryManager)

	router := httprouter.New()
	h := Handler{
		Manager: manager,
		W:       localWarden,
		H:       herodot.NewJSONWriter(nil),
	}
	h.SetRoutes(router)
	server := httptest.NewServer(router)

	client := hydra.NewJsonWebKeyApiWithBasePath(server.URL)
	client.Configuration.Transport = httpClient.Transport

	t.Run("JSON Web Key", func(t *testing.T) {
		t.Run("CreateJwkSetKey", func(t *testing.T) {
			// Create a key called set-foo
			resultKeys, _, err := client.CreateJsonWebKeySet("set-foo", hydra.JsonWebKeySetGeneratorRequest{
				Alg: "HS256",
				Kid: "key-bar",
			})
			require.NoError(t, err)
			assert.Len(t, resultKeys.Keys, 1)
			assert.Equal(t, "key-bar", resultKeys.Keys[0].Kid)
			assert.Equal(t, "HS256", resultKeys.Keys[0].Alg)
		})

		resultKeys, _, err := client.GetJsonWebKey("key-bar", "set-foo")
		t.Run("GetJwkSetKey after create", func(t *testing.T) {
			require.NoError(t, err)
			assert.Len(t, resultKeys.Keys, 1)
			assert.Equal(t, "key-bar", resultKeys.Keys[0].Kid)
			assert.Equal(t, "HS256", resultKeys.Keys[0].Alg)
		})

		t.Run("UpdateJwkSetKey", func(t *testing.T) {
			resultKeys.Keys[0].Alg = "RS256"
			resultKey, _, err := client.UpdateJsonWebKey("key-bar", "set-foo", resultKeys.Keys[0])
			require.NoError(t, err)
			assert.Equal(t, "key-bar", resultKey.Kid)
			assert.Equal(t, "RS256", resultKey.Alg)
		})

		t.Run("DeleteJwkSetKey after delete", func(t *testing.T) {
			response, err := client.DeleteJsonWebKey("key-bar", "set-foo")
			require.NoError(t, err)
			assert.Equal(t, http.StatusNoContent, response.StatusCode)
		})

		t.Run("GetJwkSetKey after delete", func(t *testing.T) {
			_, response, err := client.GetJsonWebKey("key-bar", "set-foo")
			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, response.StatusCode)
		})

	})

	t.Run("JWK Set", func(t *testing.T) {
		t.Run("CreateJwkSetKey", func(t *testing.T) {
			resultKeys, _, err := client.CreateJsonWebKeySet("set-foo2", hydra.JsonWebKeySetGeneratorRequest{
				Alg: "HS256",
				Kid: "key-bar",
			})
			require.NoError(t, err)

			assert.Len(t, resultKeys.Keys, 1)
			assert.Equal(t, "key-bar", resultKeys.Keys[0].Kid)
			assert.Equal(t, "HS256", resultKeys.Keys[0].Alg)
		})

		resultKeys, _, err := client.GetJsonWebKeySet("set-foo2")
		t.Run("GetJwkSet after create", func(t *testing.T) {
			require.NoError(t, err)
			assert.Len(t, resultKeys.Keys, 1)
			assert.Equal(t, "key-bar", resultKeys.Keys[0].Kid)
			assert.Equal(t, "HS256", resultKeys.Keys[0].Alg)
		})

		t.Run("UpdateJwkSet", func(t *testing.T) {
			resultKeys.Keys[0].Alg = "RS256"
			resultKeys, _, err = client.UpdateJsonWebKeySet("set-foo2", *resultKeys)
			require.NoError(t, err)
			assert.Len(t, resultKeys.Keys, 1)
			assert.Equal(t, "key-bar", resultKeys.Keys[0].Kid)
			assert.Equal(t, "RS256", resultKeys.Keys[0].Alg)
		})

		t.Run("DeleteJwkSet", func(t *testing.T) {
			response, err := client.DeleteJsonWebKeySet("set-foo2")
			require.NoError(t, err)
			assert.Equal(t, http.StatusNoContent, response.StatusCode)
		})

		t.Run("GetJwkSet after delete", func(t *testing.T) {
			_, response, err := client.GetJsonWebKeySet("set-foo2")
			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, response.StatusCode)
		})

		t.Run("GetJwkSetKey after delete", func(t *testing.T) {
			_, response, err := client.GetJsonWebKey("key-bar", "set-foo2")
			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, response.StatusCode)
		})
	})
}
