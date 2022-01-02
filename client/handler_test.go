package client_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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

	newServer := func(t *testing.T, dynamicEnabled bool) (*httptest.Server, *http.Client) {
		router := httprouter.New()
		h.SetRoutes(&x.RouterAdmin{Router: router}, &x.RouterPublic{Router: router}, dynamicEnabled)
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

	fetchWithAuth := func(t *testing.T, method, url, username, password string, body io.Reader) (string, *http.Response) {
		r, err := http.NewRequest(method, url, body)
		require.NoError(t, err)
		r.SetBasicAuth(username, password)
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

	createClient := func(t *testing.T, c *client.Client, ts *httptest.Server, path string) {
		body, res := makeJSON(t, ts, "POST", path, c)
		require.Equal(t, http.StatusCreated, res.StatusCode, body)
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
						req, err := http.NewRequest(method, ts.URL+client.DynClientsHandlerPath+"?client_id="+expected.OutfacingID, nil)
						require.NoError(t, err)

						res, err := hc.Do(req)
						require.NoError(t, err)
						defer res.Body.Close()

						body, err := io.ReadAll(res.Body)
						require.NoError(t, err)

						snapshotx.SnapshotTExcept(t, newResponseSnapshot(string(body), res), nil)
					})

					t.Run("without incorrect auth", func(t *testing.T) {
						body, res := fetchWithAuth(t, method, ts.URL+client.DynClientsHandlerPath+"?client_id="+expected.OutfacingID, expected.OutfacingID, "incorrect", nil)
						assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
						snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
					})

					t.Run("with a different client auth", func(t *testing.T) {
						body, res := fetchWithAuth(t, method, ts.URL+client.DynClientsHandlerPath+"?client_id="+expected.OutfacingID, secondClient.OutfacingID, secondClient.Secret, nil)
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
						Secret:       "averylongsecret",
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
						Secret:       "averylongsecret",
						RedirectURIs: []string{"http://localhost:3000/cb"},
						Metadata:     []byte(`{"foo":"bar"}`),
					},
					path:       client.DynClientsHandlerPath,
					statusCode: http.StatusBadRequest,
				},
				{
					d: "short secret fails",
					payload: &client.Client{
						OutfacingID:  "create-client-4",
						Secret:       "short",
						RedirectURIs: []string{"http://localhost:3000/cb"},
					},
					path:       client.DynClientsHandlerPath,
					statusCode: http.StatusBadRequest,
				},
			} {
				t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
					body, res := makeJSON(t, ts, "POST", tc.path, tc.payload)
					require.Equal(t, tc.statusCode, res.StatusCode, body)
					snapshotx.SnapshotTExcept(t, json.RawMessage(body), []string{"updated_at", "created_at"})
				})
			}
		})

		t.Run("case=fetching non-existing client", func(t *testing.T) {
			for _, path := range []string{
				client.DynClientsHandlerPath + "?client_id=foo",
				client.ClientsHandlerPath + "/foo",
			} {
				t.Run("path="+path, func(t *testing.T) {
					body, res := fetchWithAuth(t, "GET", ts.URL+path, expected.OutfacingID, expected.Secret, nil)
					snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
				})
			}
		})

		t.Run("case=updating non-existing client", func(t *testing.T) {
			for _, path := range []string{
				client.DynClientsHandlerPath + "?client_id=foo",
				client.ClientsHandlerPath + "/foo",
			} {
				t.Run("path="+path, func(t *testing.T) {
					body, res := fetchWithAuth(t, "PUT", ts.URL+path, expected.OutfacingID, expected.Secret, bytes.NewBufferString("{}"))
					snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
				})
			}
		})

		t.Run("case=delete non-existing client", func(t *testing.T) {
			for _, path := range []string{
				client.DynClientsHandlerPath + "?client_id=foo",
				client.ClientsHandlerPath + "/foo",
			} {
				t.Run("path="+path, func(t *testing.T) {
					body, res := fetchWithAuth(t, "DELETE", ts.URL+path, expected.OutfacingID, expected.Secret, nil)
					snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
				})
			}
		})

		t.Run("case=patching non-existing client", func(t *testing.T) {
			body, res := fetchWithAuth(t, "PATCH", ts.URL+client.ClientsHandlerPath+"/foo", expected.OutfacingID, expected.Secret, nil)
			snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), nil)
		})

		t.Run("case=fetching existing client", func(t *testing.T) {
			t.Run("endpoint=admin", func(t *testing.T) {
				body, res := fetch(t, ts.URL+client.ClientsHandlerPath+"/"+expected.OutfacingID)
				assert.Equal(t, http.StatusOK, res.StatusCode)
				snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), []string{"body.created_at", "body.updated_at"})
			})

			t.Run("endpoint=selfservice", func(t *testing.T) {
				body, res := fetchWithAuth(t, "GET", ts.URL+client.DynClientsHandlerPath+"?client_id="+expected.OutfacingID, expected.OutfacingID, expected.Secret, nil)
				assert.Equal(t, http.StatusOK, res.StatusCode)
				snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), []string{"body.created_at", "body.updated_at"})
			})
		})

		t.Run("case=updating existing client fails with metadat on self service", func(t *testing.T) {
			expected := &client.Client{
				OutfacingID:             "update-existing-client-selfservice-metadata",
				Secret:                  "averylongsecret",
				RedirectURIs:            []string{"http://localhost:3000/cb"},
				TokenEndpointAuthMethod: "client_secret_basic",
			}
			createClient(t, expected, ts, client.ClientsHandlerPath)

			// Possible to update the secret
			expected.Metadata = []byte(`{"foo":"bar"}`)
			payload, err := json.Marshal(expected)
			require.NoError(t, err)

			body, res := fetchWithAuth(t, "PUT", ts.URL+client.DynClientsHandlerPath+"?client_id="+expected.OutfacingID, expected.OutfacingID, expected.Secret, bytes.NewReader(payload))
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

			t.Run("endpoint=selfservice", func(t *testing.T) {
				expected := &client.Client{
					OutfacingID:             "update-existing-client-selfservice",
					Secret:                  "averylongsecret",
					RedirectURIs:            []string{"http://localhost:3000/cb"},
					TokenEndpointAuthMethod: "client_secret_basic",
				}
				createClient(t, expected, ts, client.ClientsHandlerPath)

				// Possible to update the secret
				expected.Secret = "anothersecret"
				expected.RedirectURIs = append(expected.RedirectURIs, "https://foobar.com")
				payload, err := json.Marshal(expected)
				require.NoError(t, err)

				body, res := fetchWithAuth(t, "PUT", ts.URL+client.DynClientsHandlerPath+"?client_id="+expected.OutfacingID, expected.OutfacingID, "averylongsecret", bytes.NewReader(payload))
				assert.Equal(t, http.StatusOK, res.StatusCode)
				snapshotx.SnapshotTExcept(t, newResponseSnapshot(body, res), []string{"body.created_at", "body.updated_at"})
			})
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
				body, res := makeJSON(t, ts, "POST", client.ClientsHandlerPath, expected)
				require.Equal(t, http.StatusCreated, res.StatusCode, body)

				_, res = fetchWithAuth(t, "DELETE", ts.URL+client.DynClientsHandlerPath+"?client_id="+expected.OutfacingID, expected.OutfacingID, "averylongsecret", nil)
				assert.Equal(t, http.StatusNoContent, res.StatusCode)
			})
		})

		t.Run("case=fetch with different self-service auth methods", func(t *testing.T) {
			for k, tc := range []struct {
				c  *client.Client
				r  func(t *testing.T, r *http.Request, c *client.Client)
				es int
			}{
				{
					c: &client.Client{
						OutfacingID:             "get-client-auth-1",
						Secret:                  "averylongsecret",
						RedirectURIs:            []string{"http://localhost:3000/cb"},
						TokenEndpointAuthMethod: "client_secret_basic",
					},
					r: func(t *testing.T, r *http.Request, c *client.Client) {
						r.SetBasicAuth(c.OutfacingID, c.Secret)
					},
					es: http.StatusOK,
				},
				{
					c: &client.Client{
						OutfacingID:             "get-client-auth-2",
						RedirectURIs:            []string{"http://localhost:3000/cb"},
						TokenEndpointAuthMethod: "none",
					},
					r: func(t *testing.T, r *http.Request, c *client.Client) {
						r.SetBasicAuth(c.OutfacingID, "")
					},
					es: http.StatusUnauthorized,
				},
				{
					c: &client.Client{
						OutfacingID:             "get-client-auth-3",
						RedirectURIs:            []string{"http://localhost:3000/cb"},
						TokenEndpointAuthMethod: "none",
					},
					r: func(t *testing.T, r *http.Request, c *client.Client) {
						r.SetBasicAuth(c.OutfacingID, "random")
					},
					es: http.StatusUnauthorized,
				},
				{
					c: &client.Client{
						OutfacingID:             "get-client-auth-4",
						Secret:                  "averylongsecret",
						RedirectURIs:            []string{"http://localhost:3000/cb"},
						TokenEndpointAuthMethod: "client_secret_post",
					},
					r: func(t *testing.T, r *http.Request, c *client.Client) {
						q := r.URL.Query()
						q.Set("client_secret", c.Secret)
						r.URL.RawQuery = q.Encode()
					},
					es: http.StatusOK,
				},
			} {
				t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
					createClient(t, tc.c, ts, client.ClientsHandlerPath)
					req, err := http.NewRequest("GET", ts.URL+client.DynClientsHandlerPath+"?client_id="+tc.c.OutfacingID, nil)
					require.NoError(t, err)
					tc.r(t, req, tc.c)

					res, err := ts.Client().Do(req)
					assert.Equal(t, tc.es, res.StatusCode)
					require.NoError(t, err)

					body, err := io.ReadAll(res.Body)
					require.NoError(t, err)

					snapshotx.SnapshotTExcept(t, newResponseSnapshot(string(body), res), []string{"body.created_at", "body.updated_at"})
				})
			}
		})
	})
}
