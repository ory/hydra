// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/configx"
	"github.com/ory/x/healthx"
	"github.com/ory/x/prometheusx"
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
				"serve.public.cors.allowed_origins":   []string{"https://example.com"},
				"serve.public.cors.enabled":           true,
				"serve.public.cors.allowed_methods":   []string{"GET"},
				"serve.public.cors.allow_credentials": true,
			},
			verifyResponse: expectCORSHeaders,
		},
		{
			name: "with CORS disabled",
			config: map[string]interface{}{
				"serve.public.cors.enabled": false,
			},
			verifyResponse: expectNoCORSHeaders,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValues(tc.config)))

			public := x.NewRouterPublic(prometheusx.NewMetricsManager("", "", "", ""))
			reg.RegisterPublicRoutes(ctx, public)

			ts := httptest.NewServer(public)

			tc.verifyResponse(t, doCORSRequest(t, ts.URL+healthx.AliveCheckPath))
			tc.verifyResponse(t, doCORSRequest(t, ts.URL+healthx.ReadyCheckPath))
		})
	}
}
