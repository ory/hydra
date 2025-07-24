// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/configx"
	"github.com/ory/x/httprouterx"
)

func TestHandlerConsent(t *testing.T) {
	t.Parallel()

	reg := testhelpers.NewRegistryMemory(t, configx.WithValue(config.KeyScopeStrategy, "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY"))

	h := reg.OAuth2Handler()
	r := x.NewRouterAdmin(reg.Config().AdminURL)
	h.SetPublicRoutes(&httprouterx.RouterPublic{Router: r.Router}, func(h http.Handler) http.Handler { return h })
	h.SetAdminRoutes(r)
	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL + oauth2.DefaultConsentPath)
	assert.Nil(t, err)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	assert.NotEmpty(t, body)
}
