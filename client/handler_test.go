package client_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tidwall/gjson"

	"github.com/ory/hydra/driver/config"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/x"
	"github.com/ory/x/snapshotx"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/internal"
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

func TestHandler(t *testing.T) {
	reg := internal.NewMockedRegistry(t)
	h := client.NewHandler(reg)

	t.Run("create client registration tokens", func(t *testing.T) {
		for k, tc := range []struct {
			c       *client.Client
			dynamic bool
		}{
			{c: &client.Client{OutfacingID: "create-client-0"}},
			{dynamic: true, c: new(client.Client)},
			{c: &client.Client{OutfacingID: "create-client-1"}},
			{c: &client.Client{Secret: "create-client-2"}},
			{c: &client.Client{OutfacingID: "create-client-3"}, dynamic: true},
		} {
			t.Run(fmt.Sprintf("case=%d/dynamic=%v", k, tc.dynamic), func(t *testing.T) {
				var b bytes.Buffer
				require.NoError(t, json.NewEncoder(&b).Encode(tc.c))
				r, err := http.NewRequest("POST", "/openid/registration", &b)
				require.NoError(t, err)

				hadSecret := len(tc.c.Secret) > 0
				c, err := h.CreateClient(r, func(c *client.Client) error {
					return nil
				}, tc.dynamic)
				require.NoError(t, err)

				except := []string{"registration_access_token", "updated_at", "created_at"}
				require.NotEmpty(t, c.RegistrationAccessToken)
				require.NotEqual(t, c.RegistrationAccessTokenSignature, c.RegistrationAccessToken)
				if !hadSecret {
					require.NotEmpty(t, c.Secret)
					except = append(except, "client_secret")
				}

				if tc.dynamic {
					require.NotEmpty(t, c.OutfacingID)
					assert.Equal(t, reg.Config().PublicURL().String()+"oauth2/register/"+c.OutfacingID, c.RegistrationClientURI)
					except = append(except, "client_id", "client_secret", "registration_client_uri")
				}

				snapshotx.SnapshotTExcept(t, c, except)
			})
		}
	})

	t.Run("dynamic client registration protocol authentication", func(t *testing.T) {
		r, err := http.NewRequest("POST", "/openid/registration", bytes.NewBufferString("{}"))
		require.NoError(t, err)
		expected, err := h.CreateClient(r, func(c *client.Client) error {
			return nil
		}, true)
		require.NoError(t, err)

		t.Run("valid auth", func(t *testing.T) {
			actual, err := h.ValidDynamicAuth(&http.Request{Header: http.Header{"Authorization": {"Bearer " + expected.RegistrationAccessToken}}}, httprouter.Params{
				httprouter.Param{Key: "id", Value: expected.OutfacingID},
			})
			require.NoError(t, err, "authentication with registration access token works")
			assert.EqualValues(t, expected.GetID(), actual.GetID())
		})

		t.Run("missing auth", func(t *testing.T) {
			_, err := h.ValidDynamicAuth(&http.Request{}, httprouter.Params{
				httprouter.Param{Key: "id", Value: expected.OutfacingID},
			})
			require.Error(t, err, "authentication without registration access token fails")
		})

		t.Run("incorrect auth", func(t *testing.T) {
			_, err := h.ValidDynamicAuth(&http.Request{Header: http.Header{"Authorization": {"Bearer invalid"}}}, httprouter.Params{
				httprouter.Param{Key: "id", Value: expected.OutfacingID},
			})
			require.Error(t, err, "authentication with invalid registration access token fails")
		})
	})

	newServer := func(t *testing.T, dynamicEnabled bool) (*httptest.Server, *http.Client) {
		require.NoError(t, reg.Config().Set(config.KeyPublicAllowDynamicRegistration, dynamicEnabled))
		router := httprouter.New()
		h.SetRoutes(&x.RouterAdmin{Router: router}, &x.RouterPublic{Router: router})
		ts := httptest.NewServer(router)
		t.Cleanup(ts.Close)
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

		for _, method := range []string{"GET", "POST", "PUT", "DELETE"} {
			t.Run("method="+method, func(t *testing.T) {
				req, err := http.NewRequest(method, ts.URL+client.DynClientsHandlerPath, nil)
				require.NoError(t, err)

				res, err := hc.Do(req)
				require.NoError(t, err)
				require.Equal(t, http.StatusNotFound, res.StatusCode)
			})
		}
	})

	t.Run("case=selfservice with incorrect or missing auth", func(t *testing.T) {
		ts, hc := newServer(t, true)
		expected := &client.Client{
			OutfacingID:             "incorrect-missing-client",
			Secret:                  "averylongsecret",
			RedirectURIs:            []string{"http://localhost:3000/cb"},
			TokenEndpointAuthMethod: "client_secret_basic",
		}
		createClient(t, expected, ts, client.ClientsHandlerPath)

		// Create the second client
		secondClient := &client.Client{
			OutfacingID:  "second-existing-client",
			Secret:       "averylongsecret",
			RedirectURIs: []string{"http://localhost:3000/cb"},
		}
		createClient(t, secondClient, ts, client.ClientsHandlerPath)

		t.Run("endpoint=selfservice", func(t *testing.T) {
			for _, method := range []string{"GET", "DELETE", "PUT"} {
				t.Run("method="+method, func(t *testing.T) {
					t.Run("without auth", func(t *testing.T) {
						req, err := http.NewRequest(method, ts.URL+client.DynClientsHandlerPath+"/"+expected.OutfacingID, nil)
						require.NoError(t, err)

						res, err := hc.Do(req)
						require.NoError(t, err)
						defer res.Body.Close()

						body, err := io.ReadAll(res.Body)
						require.NoError(t, err)

						snapshotx.SnapshotTExcept(t, newResponseSnapshot(string(body), res), nil)
					})

					t.Run("without incorrect auth", func(t *testing.T) {
						body, res := fetchWithBearerAuth(t, method, ts.URL+client.DynClientsHandlerPath+"/"+expected.OutfacingID, "incorrect", nil)
						assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
						snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
					})

					t.Run("with a different client auth", func(t *testing.T) {
						body, res := fetchWithBearerAuth(t, method, ts.URL+client.DynClientsHandlerPath+"/"+expected.OutfacingID, secondClient.RegistrationAccessToken, nil)
						assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
						snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
					})
				})
			}
		})
	})

	t.Run("common", func(t *testing.T) {
		ts, _ := newServer(t, true)
		expected := &client.Client{
			OutfacingID:             "existing-client",
			Secret:                  "averylongsecret",
			RedirectURIs:            []string{"http://localhost:3000/cb"},
			TokenEndpointAuthMethod: "client_secret_basic",
		}
		createClient(t, expected, ts, client.ClientsHandlerPath)

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
						OutfacingID:  "create-client-1",
						RedirectURIs: []string{"http://localhost:3000/cb"},
					},
					path:       client.DynClientsHandlerPath,
					statusCode: http.StatusCreated,
				},
				{
					d: "basic admin registration",
					payload: &client.Client{
						OutfacingID:  "create-client-2",
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
						OutfacingID:  "create-client-3",
						RedirectURIs: []string{"http://localhost:3000/cb"},
						Metadata:     []byte(`{"foo":"bar"}`),
					},
					path:       client.DynClientsHandlerPath,
					statusCode: http.StatusBadRequest,
				},
				{
					d: "short secret fails for admin",
					payload: &client.Client{
						OutfacingID:  "create-client-4",
						Secret:       "short",
						RedirectURIs: []string{"http://localhost:3000/cb"},
					},
					path:       client.ClientsHandlerPath,
					statusCode: http.StatusBadRequest,
				},
				{
					d: "basic dynamic client registration",
					payload: &client.Client{
						OutfacingID:  "create-client-5",
						Secret:       "averylongsecret",
						RedirectURIs: []string{"http://localhost:3000/cb"},
					},
					path:       client.DynClientsHandlerPath,
					statusCode: http.StatusForbidden,
				},
			} {
				t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
					body, res := makeJSON(t, ts, "POST", tc.path, tc.payload)
					require.Equal(t, tc.statusCode, res.StatusCode, body)
					exclude := []string{"updated_at", "created_at", "registration_access_token"}
					if tc.path == client.DynClientsHandlerPath {
						exclude = append(exclude, "client_id", "client_secret", "registration_client_uri")
					}
					if tc.statusCode == http.StatusOK {
						for _, key := range exclude {
							assert.NotEmpty(t, gjson.Get(body, key).String(), "%s in %s", key, body)
						}
					}
					snapshotx.SnapshotTExcept(t, json.RawMessage(body), exclude)
				})
			}
		})

		t.Run("case=fetching non-existing client", func(t *testing.T) {
			for _, path := range []string{
				client.DynClientsHandlerPath + "/foo",
				client.ClientsHandlerPath + "/foo",
			} {
				t.Run("path="+path, func(t *testing.T) {
					body, res := fetchWithBearerAuth(t, "GET", ts.URL+path, expected.RegistrationAccessToken, nil)
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
				OutfacingID:  "existing-client-fetch",
				Secret:       "rdetzfuzgihojuzgtfrdes",
				RedirectURIs: []string{"http://localhost:3000/cb"},
			}, ts, client.ClientsHandlerPath)
			id := gjson.Get(expected, "client_id").String()
			rat := gjson.Get(expected, "registration_access_token").String()

			t.Run("endpoint=admin", func(t *testing.T) {
				body, res := fetch(t, ts.URL+client.ClientsHandlerPath+"/"+id)
				assert.Equal(t, http.StatusOK, res.StatusCode)
				assert.Equal(t, id, gjson.Get(body, "client_id").String())
				snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), []string{"body.created_at", "body.updated_at"})
			})

			t.Run("endpoint=selfservice", func(t *testing.T) {
				body, res := fetchWithBearerAuth(t, "GET", ts.URL+client.DynClientsHandlerPath+"/"+id, rat, nil)
				assert.Equal(t, http.StatusOK, res.StatusCode)
				assert.Equal(t, id, gjson.Get(body, "client_id").String())
				assert.False(t, gjson.Get(body, "metadata").Bool())
				snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), []string{"body.created_at", "body.updated_at"})
			})
		})

		t.Run("case=updating existing client fails with metadata on self service", func(t *testing.T) {
			expected := &client.Client{
				OutfacingID:             "update-existing-client-selfservice-metadata",
				Secret:                  "averylongsecret",
				RedirectURIs:            []string{"http://localhost:3000/cb"},
				TokenEndpointAuthMethod: "client_secret_basic",
			}
			body := createClient(t, expected, ts, client.ClientsHandlerPath)

			// Possible to update the secret
			expected.Metadata = []byte(`{"foo":"bar"}`)
			expected.Secret = ""
			payload, err := json.Marshal(expected)
			require.NoError(t, err)

			body, res := fetchWithBearerAuth(t, "PUT", ts.URL+client.DynClientsHandlerPath+"/"+expected.OutfacingID, gjson.Get(body, "registration_access_token").String(), bytes.NewReader(payload))
			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
			snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
		})

		t.Run("case=updating existing client", func(t *testing.T) {
			t.Run("endpoint=admin", func(t *testing.T) {
				expected := &client.Client{
					OutfacingID:             "update-existing-client-admin",
					Secret:                  "averylongsecret",
					RedirectURIs:            []string{"http://localhost:3000/cb"},
					TokenEndpointAuthMethod: "client_secret_basic",
				}
				createClient(t, expected, ts, client.ClientsHandlerPath)

				expected.RedirectURIs = append(expected.RedirectURIs, "https://foobar.com")
				body, res := makeJSON(t, ts, "PUT", client.ClientsHandlerPath+"/"+expected.OutfacingID, expected)
				assert.Equal(t, http.StatusOK, res.StatusCode)
				snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), []string{"body.created_at", "body.updated_at"})
			})

			t.Run("endpoint=dynamic client registration", func(t *testing.T) {
				expected := &client.Client{
					OutfacingID:             "update-existing-client-selfservice",
					Secret:                  "averylongsecret",
					RedirectURIs:            []string{"http://localhost:3000/cb"},
					TokenEndpointAuthMethod: "client_secret_basic",
				}
				actual := createClient(t, expected, ts, client.ClientsHandlerPath)

				// Possible to update the secret
				expected.RedirectURIs = append(expected.RedirectURIs, "https://foobar.com")
				expected.Secret = ""
				payload, err := json.Marshal(expected)
				require.NoError(t, err)

				originalRAT := gjson.Get(actual, "registration_access_token").String()
				body, res := fetchWithBearerAuth(t, "PUT", ts.URL+client.DynClientsHandlerPath+"/"+expected.OutfacingID, originalRAT, bytes.NewReader(payload))
				assert.Equal(t, http.StatusOK, res.StatusCode)
				newToken := gjson.Get(body, "registration_access_token").String()
				assert.NotEmpty(t, newToken)
				require.NotEqual(t, originalRAT, newToken, "the new token should be different from the old token")
				snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), []string{"body.created_at", "body.updated_at", "body.registration_access_token"})

				_, res = fetchWithBearerAuth(t, "GET", ts.URL+client.DynClientsHandlerPath+"/"+expected.OutfacingID, originalRAT, bytes.NewReader(payload))
				assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
				body, res = fetchWithBearerAuth(t, "GET", ts.URL+client.DynClientsHandlerPath+"/"+expected.OutfacingID, newToken, bytes.NewReader(payload))
				assert.Equal(t, http.StatusOK, res.StatusCode)
				assert.Empty(t, gjson.Get(body, "registration_access_token").String())
			})

			t.Run("endpoint=dynamic client registration does not allow changing the secret", func(t *testing.T) {
				expected := &client.Client{
					OutfacingID:             "update-existing-client-no-secret-change",
					RedirectURIs:            []string{"http://localhost:3000/cb"},
					TokenEndpointAuthMethod: "client_secret_basic",
				}
				actual := createClient(t, expected, ts, client.ClientsHandlerPath)

				// Possible to update the secret
				expected.Secret = "anothersecret"
				expected.RedirectURIs = append(expected.RedirectURIs, "https://foobar.com")
				payload, err := json.Marshal(expected)
				require.NoError(t, err)

				originalRAT := gjson.Get(actual, "registration_access_token").String()
				body, res := fetchWithBearerAuth(t, "PUT", ts.URL+client.DynClientsHandlerPath+"/"+expected.OutfacingID, originalRAT, bytes.NewReader(payload))
				assert.Equal(t, http.StatusForbidden, res.StatusCode)
				snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
			})
		})

		t.Run("case=creating a client dynamically does not allow setting the secret", func(t *testing.T) {
			body, res := makeJSON(t, ts, "POST", client.DynClientsHandlerPath, &client.Client{
				TokenEndpointAuthMethod: "client_secret_basic",
				Secret:                  "foobarbaz",
			})
			require.Equal(t, http.StatusForbidden, res.StatusCode, body)
			snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
		})

		t.Run("case=delete existing client", func(t *testing.T) {
			t.Run("endpoint=admin", func(t *testing.T) {
				expected := &client.Client{
					OutfacingID:             "delete-existing-client-admin",
					Secret:                  "averylongsecret",
					RedirectURIs:            []string{"http://localhost:3000/cb"},
					TokenEndpointAuthMethod: "client_secret_basic",
				}
				body, res := makeJSON(t, ts, "POST", client.ClientsHandlerPath, expected)
				require.Equal(t, http.StatusCreated, res.StatusCode, body)

				_, res = makeJSON(t, ts, "DELETE", client.ClientsHandlerPath+"/"+expected.OutfacingID, nil)
				assert.Equal(t, http.StatusNoContent, res.StatusCode)
			})

			t.Run("endpoint=selfservice", func(t *testing.T) {
				expected := &client.Client{
					OutfacingID:             "delete-existing-client-selfservice",
					Secret:                  "averylongsecret",
					RedirectURIs:            []string{"http://localhost:3000/cb"},
					TokenEndpointAuthMethod: "client_secret_basic",
				}
				actual, res := makeJSON(t, ts, "POST", client.ClientsHandlerPath, expected)
				require.Equal(t, http.StatusCreated, res.StatusCode, actual)

				originalRAT := gjson.Get(actual, "registration_access_token").String()
				_, res = fetchWithBearerAuth(t, "DELETE", ts.URL+client.DynClientsHandlerPath+"/"+expected.OutfacingID, originalRAT, nil)
				assert.Equal(t, http.StatusNoContent, res.StatusCode)
			})
		})
	})
}
