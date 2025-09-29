// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/x/configx"
	"github.com/ory/x/httprouterx"
	"github.com/ory/x/prometheusx"
)

func TestHandlerConsent(t *testing.T) {
	t.Parallel()

	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValue(config.KeyScopeStrategy, "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY")))

	h := oauth2.NewHandler(reg)
	metrics := prometheusx.NewMetricsManagerWithPrefix("hydra", prometheusx.HTTPMetrics, config.Version, config.Commit, config.Date)
	r := httprouterx.NewRouterAdminWithPrefix(metrics)
	h.SetPublicRoutes(r.ToPublic(), func(h http.Handler) http.Handler { return h })
	h.SetAdminRoutes(r)
	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL + oauth2.DefaultConsentPath)
	assert.Nil(t, err)
	defer res.Body.Close() //nolint:errcheck

	body, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	assert.NotEmpty(t, body)
}
