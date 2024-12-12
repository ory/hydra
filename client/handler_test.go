// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ory/x/sqlxx"

	"github.com/ory/x/httprouterx"

	"github.com/tidwall/sjson"

	"github.com/gofrs/uuid"
	"github.com/tidwall/gjson"

	"github.com/ory/hydra/v2/internal/testhelpers"

	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/x/contextx"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"

	"github.com/ory/x/snapshotx"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/client"
)

type responseSnapshot struct {
	Body   json.RawMessage `json:"body"`
	Status int             `json:"status"`
}

func newResponseSnapshot(body string, res *http.Response) *responseSnapshot {
	return &responseSnapshot{
		Body:   json.RawMessage(body),
		Status: res.StatusCode,
	}
}

func getClientID(body string) string {
	return gjson.Get(body, "client_id").String()
}

func TestHandler(t *testing.T) {
	ctx := context.Background()
	reg := testhelpers.NewMockedRegistry(t, &contextx.Default{})
	h := client.NewHandler(reg)
	reg.WithContextualizer(&contextx.TestContextualizer{})

	t.Run("create client registration tokens", func(t *testing.T) {
		for k, tc := range []struct {
			c       *client.Client
			dynamic bool
		}{
			{dynamic: true, c: new(client.Client)},
			{c: new(client.Client)},
			{c: &client.Client{Secret: "01bbf13a-ae3e-44d5-b4b4-dd78137041be"}},
		} {
			t.Run(fmt.Sprintf("case=%d/dynamic=%v", k, tc.dynamic), func(t *testing.T) {
				var b bytes.Buffer
				require.NoError(t, json.NewEncoder(&b).Encode(tc.c))
				r, err := http.NewRequest("POST", "/openid/registration", &b)
				require.NoError(t, err)

				hadSecret := len(tc.c.Secret) > 0
				c, err := h.CreateClient(r, func(ctx context.Context, c *client.Client) error {
					return nil
				}, tc.dynamic)
				require.NoError(t, err)
				require.NotEqual(t, c.NID, uuid.Nil)

				except := []string{"client_id", "registration_access_token", "updated_at", "created_at", "registration_client_uri"}
				require.NotEmpty(t, c.RegistrationAccessToken)
				require.NotEqual(t, c.RegistrationAccessTokenSignature, c.RegistrationAccessToken)
				if !hadSecret {
					require.NotEmpty(t, c.Secret)
					except = append(except, "client_secret")
				}

				if tc.dynamic {
					require.NotEmpty(t, c.GetID())
					assert.Equal(t, reg.Config().PublicURL(ctx).String()+"oauth2/register/"+c.GetID(), c.RegistrationClientURI)
					except = append(except, "client_id", "client_secret", "registration_client_uri")
				}

				snapshotx.SnapshotTExcept(t, c, except)
			})
		}
	})

	t.Run("dynamic client registration protocol authentication", func(t *testing.T) {
		r, err := http.NewRequest("POST", "/openid/registration", bytes.NewBufferString("{}"))
		require.NoError(t, err)
		expected, err := h.CreateClient(r, func(ctx context.Context, c *client.Client) error {
			return nil
		}, true)
		require.NoError(t, err)

		t.Run("valid auth", func(t *testing.T) {
			actual, err := h.ValidDynamicAuth(&http.Request{Header: http.Header{"Authorization": {"Bearer " + expected.RegistrationAccessToken}}}, httprouter.Params{
				{Key: "id", Value: expected.GetID()},
			})
			require.NoError(t, err, "authentication with registration access token works")
			assert.EqualValues(t, expected.GetID(), actual.GetID())
		})

		t.Run("missing auth", func(t *testing.T) {
			_, err := h.ValidDynamicAuth(&http.Request{}, httprouter.Params{
				{Key: "id", Value: expected.GetID()},
			})
			require.Error(t, err, "authentication without registration access token fails")
		})

		t.Run("incorrect auth", func(t *testing.T) {
			_, err := h.ValidDynamicAuth(&http.Request{Header: http.Header{"Authorization": {"Bearer invalid"}}}, httprouter.Params{
				{Key: "id", Value: expected.GetID()},
			})
			require.Error(t, err, "authentication with invalid registration access token fails")
		})
	})

	newServer := func(t *testing.T, dynamicEnabled bool) (*httptest.Server, *http.Client) {
		require.NoError(t, reg.Config().Set(ctx, config.KeyPublicAllowDynamicRegistration, dynamicEnabled))
		router := httprouter.New()
		h.SetRoutes(httprouterx.NewRouterAdminWithPrefixAndRouter(router, "/admin", reg.Config().AdminURL), &httprouterx.RouterPublic{Router: router})
		ts := httptest.NewServer(router)
		t.Cleanup(ts.Close)
		reg.Config().MustSet(ctx, config.KeyAdminURL, ts.URL)
		return ts, ts.Client()
	}

	fetch := func(t *testing.T, url string) (string, *http.Response) {
		res, err := http.Get(url)
		require.NoError(t, err)
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		return string(body), res
	}

	fetchWithBearerAuth := func(t *testing.T, method, url, token string, body io.Reader) (string, *http.Response) {
		r, err := http.NewRequest(method, url, body)
		require.NoError(t, err)
		r.Header.Set("Authorization", "Bearer "+token)
		res, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		defer res.Body.Close()
		out, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		return string(out), res
	}

	makeJSON := func(t *testing.T, ts *httptest.Server, method string, path string, body interface{}) (string, *http.Response) {
		var b bytes.Buffer
		require.NoError(t, json.NewEncoder(&b).Encode(body))
		r, err := http.NewRequest(method, ts.URL+path, &b)
		require.NoError(t, err)
		r.Header.Set("Content-Type", "application/json")
		res, err := ts.Client().Do(r)
		require.NoError(t, err)
		defer res.Body.Close()
		rb, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		return string(rb), res
	}

	createClient := func(t *testing.T, c *client.Client, ts *httptest.Server, path string) string {
		body, res := makeJSON(t, ts, "POST", path, c)
		require.Equal(t, http.StatusCreated, res.StatusCode, body)
		return body
	}

	t.Run("selfservice disabled", func(t *testing.T) {
		ts, hc := newServer(t, false)

		trap := &client.Client{}
		actual := createClient(t, trap, ts, client.ClientsHandlerPath)
		actualID := getClientID(actual)

		for _, tc := range []struct {
			method string
			path   string
		}{
			{method: "GET", path: ts.URL + client.DynClientsHandlerPath + "/" + actualID},
			{method: "POST", path: ts.URL + client.DynClientsHandlerPath},
			{method: "PUT", path: ts.URL + client.DynClientsHandlerPath + "/" + actualID},
			{method: "DELETE", path: ts.URL + client.DynClientsHandlerPath + "/" + actualID},
		} {
			t.Run("method="+tc.method, func(t *testing.T) {
				req, err := http.NewRequest(tc.method, tc.path, nil)
				require.NoError(t, err)

				res, err := hc.Do(req)
				require.NoError(t, err)
				require.Equal(t, http.StatusNotFound, res.StatusCode)
			})
		}
	})

	t.Run("case=selfservice with incorrect or missing auth", func(t *testing.T) {
		ts, hc := newServer(t, true)
		expectedFirst := createClient(t, &client.Client{
			Secret:                  "averylongsecret",
			RedirectURIs:            []string{"http://localhost:3000/cb"},
			TokenEndpointAuthMethod: "client_secret_basic",
		}, ts, client.ClientsHandlerPath)
		expectedFirstID := getClientID(expectedFirst)

		// Create the second client
		expectedSecond := createClient(t, &client.Client{
			Secret:       "averylongsecret",
			RedirectURIs: []string{"http://localhost:3000/cb"},
		}, ts, client.ClientsHandlerPath)
		expectedSecondID := getClientID(expectedSecond)

		t.Run("endpoint=selfservice", func(t *testing.T) {
			for _, method := range []string{"GET", "DELETE", "PUT"} {
				t.Run("method="+method, func(t *testing.T) {
					t.Run("without auth", func(t *testing.T) {
						req, err := http.NewRequest(method, ts.URL+client.DynClientsHandlerPath+"/"+expectedFirstID, nil)
						require.NoError(t, err)

						res, err := hc.Do(req)
						require.NoError(t, err)
						defer res.Body.Close()

						body, err := io.ReadAll(res.Body)
						require.NoError(t, err)

						snapshotx.SnapshotTExcept(t, newResponseSnapshot(string(body), res), nil)
					})

					t.Run("without incorrect auth", func(t *testing.T) {
						body, res := fetchWithBearerAuth(t, method, ts.URL+client.DynClientsHandlerPath+"/"+expectedFirstID, "incorrect", nil)
						assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
						snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
					})

					t.Run("with a different client auth", func(t *testing.T) {
						body, res := fetchWithBearerAuth(t, method, ts.URL+client.DynClientsHandlerPath+"/"+expectedFirstID, expectedSecondID, nil)
						assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
						snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
					})
				})
			}
		})
	})

	t.Run("common", func(t *testing.T) {
		ts, _ := newServer(t, true)
		expected := createClient(t, &client.Client{
			Secret:                  "averylongsecret",
			RedirectURIs:            []string{"http://localhost:3000/cb"},
			TokenEndpointAuthMethod: "client_secret_basic",
		}, ts, client.ClientsHandlerPath)

		t.Run("case=create clients", func(t *testing.T) {
			for k, tc := range []struct {
				d          string
				payload    *client.Client
				path       string
				statusCode int
			}{
				{
					d: "basic dynamic client registration",
					payload: &client.Client{
						RedirectURIs: []string{"http://localhost:3000/cb"},
					},
					path:       client.DynClientsHandlerPath,
					statusCode: http.StatusCreated,
				},
				{
					d: "basic admin registration",
					payload: &client.Client{
						Secret:       "averylongsecret",
						RedirectURIs: []string{"http://localhost:3000/cb"},
						Metadata:     []byte(`{"foo":"bar"}`),
					},
					path:       client.ClientsHandlerPath,
					statusCode: http.StatusCreated,
				},
				{
					d: "metadata fails for dynamic client registration",
					payload: &client.Client{
						RedirectURIs: []string{"http://localhost:3000/cb"},
						Metadata:     []byte(`{"foo":"bar"}`),
					},
					path:       client.DynClientsHandlerPath,
					statusCode: http.StatusBadRequest,
				},
				{
					d: "short secret fails for admin",
					payload: &client.Client{
						Secret:       "short",
						RedirectURIs: []string{"http://localhost:3000/cb"},
					},
					path:       client.ClientsHandlerPath,
					statusCode: http.StatusBadRequest,
				},
				{
					d: "non-uuid works",
					payload: &client.Client{
						ID:           "not-a-uuid",
						Secret:       "averylongsecret",
						RedirectURIs: []string{"http://localhost:3000/cb"},
					},
					path:       client.ClientsHandlerPath,
					statusCode: http.StatusCreated,
				},
				{
					d: "setting client id as uuid works",
					payload: &client.Client{
						ID:           "98941dac-f963-4468-8a23-9483b1e04e3c",
						Secret:       "not too short",
						RedirectURIs: []string{"http://localhost:3000/cb"},
					},
					path:       client.ClientsHandlerPath,
					statusCode: http.StatusCreated,
				},
				{
					d: "setting access token strategy fails",
					payload: &client.Client{
						RedirectURIs:        []string{"http://localhost:3000/cb"},
						AccessTokenStrategy: "jwt",
					},
					path:       client.DynClientsHandlerPath,
					statusCode: http.StatusBadRequest,
				},
				{
					d: "setting skip_consent fails for dynamic registration",
					payload: &client.Client{
						RedirectURIs: []string{"http://localhost:3000/cb"},
						SkipConsent:  true,
					},
					path:       client.DynClientsHandlerPath,
					statusCode: http.StatusBadRequest,
				},
				{
					d: "setting skip_consent succeeds for admin registration",
					payload: &client.Client{
						RedirectURIs: []string{"http://localhost:3000/cb"},
						Secret:       "2SKZkBf2P5g4toAXXnCrr~_sDM",
						SkipConsent:  true,
					},
					path:       client.ClientsHandlerPath,
					statusCode: http.StatusCreated,
				},
				{
					d: "setting skip_logout_consent fails for dynamic registration",
					payload: &client.Client{
						RedirectURIs:      []string{"http://localhost:3000/cb"},
						SkipLogoutConsent: sqlxx.NullBool{Bool: true, Valid: true},
					},
					path:       client.DynClientsHandlerPath,
					statusCode: http.StatusBadRequest,
				},
				{
					d: "setting skip_logout_consent succeeds for admin registration",
					payload: &client.Client{
						RedirectURIs:      []string{"http://localhost:3000/cb"},
						SkipLogoutConsent: sqlxx.NullBool{Bool: true, Valid: true},
						Secret:            "2SKZkBf2P5g4toAXXnCrr~_sDM",
					},
					path:       client.ClientsHandlerPath,
					statusCode: http.StatusCreated,
				},
				{
					d: "basic dynamic client registration",
					payload: &client.Client{
						ID:           "ead800c5-a316-4d0c-bf00-d25666ba72cf",
						Secret:       "averylongsecret",
						RedirectURIs: []string{"http://localhost:3000/cb"},
					},
					path:       client.DynClientsHandlerPath,
					statusCode: http.StatusBadRequest,
				},
				{
					d: "empty ID succeeds",
					payload: &client.Client{
						Secret:       "averylongsecret",
						RedirectURIs: []string{"http://localhost:3000/cb"},
					},
					path:       client.ClientsHandlerPath,
					statusCode: http.StatusCreated,
				},
			} {
				t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
					body, res := makeJSON(t, ts, "POST", tc.path, tc.payload)
					require.Equal(t, tc.statusCode, res.StatusCode, body)
					exclude := []string{"updated_at", "created_at", "registration_access_token"}
					if tc.path == client.DynClientsHandlerPath {
						exclude = append(exclude, "client_id", "client_secret", "registration_client_uri")
					}
					if tc.payload.ID == "" {
						exclude = append(exclude, "client_id", "registration_client_uri")
						assert.NotEqual(t, uuid.Nil.String(), gjson.Get(body, "client_id").String(), body)
					}
					if tc.statusCode == http.StatusOK {
						for _, key := range exclude {
							assert.NotEmpty(t, gjson.Get(body, key).String(), "%s in %s", key, body)
						}
					}
					snapshotx.SnapshotT(t, json.RawMessage(body), snapshotx.ExceptPaths(exclude...))
				})
			}
		})

		t.Run("case=fetching non-existing client", func(t *testing.T) {
			for _, path := range []string{
				client.DynClientsHandlerPath + "/foo",
				client.ClientsHandlerPath + "/foo",
			} {
				t.Run("path="+path, func(t *testing.T) {
					body, res := fetchWithBearerAuth(t, "GET", ts.URL+path, gjson.Get(expected, "registration_access_token").String(), nil)
					snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
				})
			}
		})

		t.Run("case=updating non-existing client", func(t *testing.T) {
			for _, path := range []string{
				client.DynClientsHandlerPath + "/foo",
				client.ClientsHandlerPath + "/foo",
			} {
				t.Run("path="+path, func(t *testing.T) {
					body, res := fetchWithBearerAuth(t, "PUT", ts.URL+path, "invalid", bytes.NewBufferString("{}"))
					snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
				})
			}
		})

		t.Run("case=delete non-existing client", func(t *testing.T) {
			for _, path := range []string{
				client.DynClientsHandlerPath + "/foo",
				client.ClientsHandlerPath + "/foo",
			} {
				t.Run("path="+path, func(t *testing.T) {
					body, res := fetchWithBearerAuth(t, "DELETE", ts.URL+path, "invalid", nil)
					snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
				})
			}
		})

		t.Run("case=patching non-existing client", func(t *testing.T) {
			body, res := fetchWithBearerAuth(t, "PATCH", ts.URL+client.ClientsHandlerPath+"/foo", "", nil)
			snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
		})

		t.Run("case=fetching existing client", func(t *testing.T) {
			expected := createClient(t, &client.Client{
				Secret:       "rdetzfuzgihojuzgtfrdes",
				RedirectURIs: []string{"http://localhost:3000/cb"},
			}, ts, client.ClientsHandlerPath)
			id := gjson.Get(expected, "client_id").String()
			rat := gjson.Get(expected, "registration_access_token").String()

			t.Run("endpoint=admin", func(t *testing.T) {
				body, res := fetch(t, ts.URL+client.ClientsHandlerPath+"/"+id)
				assert.Equal(t, http.StatusOK, res.StatusCode)
				assert.Equal(t, id, gjson.Get(body, "client_id").String())
				snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), []string{"body.client_id", "body.created_at", "body.updated_at"})
			})

			t.Run("endpoint=selfservice", func(t *testing.T) {
				body, res := fetchWithBearerAuth(t, "GET", ts.URL+client.DynClientsHandlerPath+"/"+id, rat, nil)
				assert.Equal(t, http.StatusOK, res.StatusCode)
				assert.Equal(t, id, gjson.Get(body, "client_id").String())
				assert.False(t, gjson.Get(body, "metadata").Bool())
				snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), []string{"body.client_id", "body.created_at", "body.updated_at"})
			})
		})

		t.Run("case=updating existing client fails with metadata on self service", func(t *testing.T) {
			expected := createClient(t, &client.Client{
				Secret:                  "averylongsecret",
				RedirectURIs:            []string{"http://localhost:3000/cb"},
				TokenEndpointAuthMethod: "client_secret_basic",
			}, ts, client.ClientsHandlerPath)
			id := gjson.Get(expected, "client_id").String()

			// Possible to update the secret
			payload, err := sjson.SetRaw(expected, "metadata", `{"foo":"bar"}`)
			require.NoError(t, err)

			payload, err = sjson.Set(payload, "client_secret", "")
			require.NoError(t, err)

			body, res := fetchWithBearerAuth(t, "PUT", ts.URL+client.DynClientsHandlerPath+"/"+id, gjson.Get(expected, "registration_access_token").String(), bytes.NewBufferString(payload))
			assert.Equal(t, http.StatusBadRequest, res.StatusCode, "%s\n%s", body, payload)
			snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
		})

		t.Run("case=updating existing client", func(t *testing.T) {
			t.Run("endpoint=admin", func(t *testing.T) {
				expected := createClient(t, &client.Client{
					Secret:                  "averylongsecret",
					RedirectURIs:            []string{"http://localhost:3000/cb"},
					TokenEndpointAuthMethod: "client_secret_basic",
				}, ts, client.ClientsHandlerPath)
				expectedID := getClientID(expected)

				payload, _ := sjson.Set(expected, "redirect_uris", []string{"http://localhost:3000/cb", "https://foobar.com"})
				body, res := makeJSON(t, ts, "PUT", client.ClientsHandlerPath+"/"+expectedID, json.RawMessage(payload))
				assert.Equal(t, http.StatusOK, res.StatusCode)
				snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), []string{"body.created_at", "body.updated_at", "body.client_id", "body.registration_client_uri", "body.registration_access_token"})
			})

			t.Run("endpoint=dynamic client registration", func(t *testing.T) {
				expected := createClient(t, &client.Client{
					Secret:                  "averylongsecret",
					RedirectURIs:            []string{"http://localhost:3000/cb"},
					TokenEndpointAuthMethod: "client_secret_basic",
				}, ts, client.ClientsHandlerPath)
				expectedID := getClientID(expected)

				// Possible to update the secret
				payload, _ := sjson.Set(expected, "redirect_uris", []string{"http://localhost:3000/cb", "https://foobar.com"})
				payload, _ = sjson.Delete(payload, "client_secret")
				payload, _ = sjson.Delete(payload, "metadata")

				originalRAT := gjson.Get(expected, "registration_access_token").String()
				body, res := fetchWithBearerAuth(t, "PUT", ts.URL+client.DynClientsHandlerPath+"/"+expectedID, originalRAT, bytes.NewBufferString(payload))
				assert.Equal(t, http.StatusOK, res.StatusCode, "%s\n%s", body, payload)
				newToken := gjson.Get(body, "registration_access_token").String()
				assert.NotEmpty(t, newToken)
				require.NotEqual(t, originalRAT, newToken, "the new token should be different from the old token")
				snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), []string{"body.created_at", "body.updated_at", "body.registration_access_token", "body.client_id", "body.registration_client_uri"})

				_, res = fetchWithBearerAuth(t, "GET", ts.URL+client.DynClientsHandlerPath+"/"+expectedID, originalRAT, bytes.NewBufferString(payload))
				assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
				body, res = fetchWithBearerAuth(t, "GET", ts.URL+client.DynClientsHandlerPath+"/"+expectedID, newToken, bytes.NewBufferString(payload))
				assert.Equal(t, http.StatusOK, res.StatusCode)
				assert.Empty(t, gjson.Get(body, "registration_access_token").String())
			})

			t.Run("endpoint=dynamic client registration does not allow changing the secret", func(t *testing.T) {
				expected := createClient(t, &client.Client{
					RedirectURIs:            []string{"http://localhost:3000/cb"},
					TokenEndpointAuthMethod: "client_secret_basic",
				}, ts, client.ClientsHandlerPath)
				expectedID := getClientID(expected)

				// Possible to update the secret
				payload, _ := sjson.Set(expected, "redirect_uris", []string{"http://localhost:3000/cb", "https://foobar.com"})
				payload, _ = sjson.Set(payload, "secret", "")

				originalRAT := gjson.Get(expected, "registration_access_token").String()
				body, res := fetchWithBearerAuth(t, "PUT", ts.URL+client.DynClientsHandlerPath+"/"+expectedID, originalRAT, bytes.NewBufferString(payload))
				assert.Equal(t, http.StatusForbidden, res.StatusCode)
				snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
			})
		})

		t.Run("case=creating a client dynamically does not allow setting the secret", func(t *testing.T) {
			body, res := makeJSON(t, ts, "POST", client.DynClientsHandlerPath, &client.Client{
				TokenEndpointAuthMethod: "client_secret_basic",
				Secret:                  "foobarbaz",
			})
			require.Equal(t, http.StatusBadRequest, res.StatusCode, body)
			snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
		})

		t.Run("case=update the lifespans of an OAuth2 client", func(t *testing.T) {
			expected := &client.Client{
				Name:                    "update-existing-client-lifespans",
				Secret:                  "averylongsecret",
				RedirectURIs:            []string{"http://localhost:3000/cb"},
				TokenEndpointAuthMethod: "client_secret_basic",
			}
			body, res := makeJSON(t, ts, "POST", client.ClientsHandlerPath, expected)
			require.Equal(t, http.StatusCreated, res.StatusCode, body)

			body, res = makeJSON(t, ts, "PUT", client.ClientsHandlerPath+"/"+gjson.Get(body, "client_id").String()+"/lifespans", testhelpers.TestLifespans)
			require.Equal(t, http.StatusOK, res.StatusCode, body)
			snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), []string{"body.client_id", "body.created_at", "body.updated_at"})
		})

		t.Run("case=delete existing client", func(t *testing.T) {
			t.Run("endpoint=admin", func(t *testing.T) {
				expected := createClient(t, &client.Client{
					RedirectURIs:            []string{"http://localhost:3000/cb"},
					TokenEndpointAuthMethod: "client_secret_basic",
				}, ts, client.ClientsHandlerPath)
				expectedID := getClientID(expected)

				_, res := makeJSON(t, ts, "DELETE", client.ClientsHandlerPath+"/"+expectedID, nil)
				assert.Equal(t, http.StatusNoContent, res.StatusCode)
			})

			t.Run("endpoint=selfservice", func(t *testing.T) {
				expected := createClient(t, &client.Client{
					Secret:                  "averylongsecret",
					RedirectURIs:            []string{"http://localhost:3000/cb"},
					TokenEndpointAuthMethod: "client_secret_basic",
				}, ts, client.ClientsHandlerPath)
				expectedID := getClientID(expected)

				originalRAT := gjson.Get(expected, "registration_access_token").String()
				_, res := fetchWithBearerAuth(t, "DELETE", ts.URL+client.DynClientsHandlerPath+"/"+expectedID, originalRAT, nil)
				assert.Equal(t, http.StatusNoContent, res.StatusCode)
			})
		})
	})
}
