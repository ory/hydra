// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ory/hydra/v2/internal/testhelpers"

	"github.com/stretchr/testify/assert"

	"github.com/ory/x/contextx"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/healthx"
)

func TestPublicHealthHandler(t *testing.T) {
	ctx := context.Background()

	doCORSRequest := func(t *testing.T, endpoint string) *http.Response {
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		require.NoError(t, err)
		req.Host = "example.com"
		req.Header.Add("Origin", "https://example.com")
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		return resp
	}

	expectCORSHeaders := func(t *testing.T, resp *http.Response) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "Origin", resp.Header.Get("Vary"))
		assert.Equal(t, "https://example.com", resp.Header.Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", resp.Header.Get("Access-Control-Allow-Credentials"))
	}

	expectNoCORSHeaders := func(t *testing.T, resp *http.Response) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotEqual(t, "Origin", resp.Header.Get("Vary"))
		assert.Equal(t, "", resp.Header.Get("Access-Control-Allow-Origin"))
	}

	for _, tc := range []struct {
		name           string
		config         map[string]interface{}
		verifyResponse func(t *testing.T, resp *http.Response)
	}{
		{
			name: "with CORS enabled",
			config: map[string]interface{}{
				"cors.allowed_origins":   []string{"https://example.com"},
				"cors.enabled":           true,
				"cors.allowed_methods":   []string{"GET"},
				"cors.allow_credentials": true,
			},
			verifyResponse: expectCORSHeaders,
		},
		{
			name: "with CORS disabled",
			config: map[string]interface{}{
				"cors.enabled": false,
			},
			verifyResponse: expectNoCORSHeaders,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			conf := testhelpers.NewConfigurationWithDefaults()
			for k, v := range tc.config {
				conf.MustSet(ctx, config.PublicInterface.Key(k), v)
			}

			reg := testhelpers.NewRegistryMemory(t, conf, &contextx.Default{})

			public := x.NewRouterPublic()
			reg.RegisterRoutes(ctx, x.NewRouterAdmin(conf.AdminURL), public)

			ts := httptest.NewServer(public)

			tc.verifyResponse(t, doCORSRequest(t, ts.URL+healthx.AliveCheckPath))
			tc.verifyResponse(t, doCORSRequest(t, ts.URL+healthx.ReadyCheckPath))
		})
	}
}
