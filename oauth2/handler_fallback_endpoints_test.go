// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ory/x/httprouterx"

	"github.com/ory/hydra/x"
	"github.com/ory/x/contextx"

	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/oauth2"

	"github.com/stretchr/testify/assert"
)

func TestHandlerConsent(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	conf.MustSet(context.Background(), config.KeyScopeStrategy, "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY")
	reg := internal.NewRegistryMemory(t, conf, &contextx.Default{})

	h := reg.OAuth2Handler()
	r := x.NewRouterAdmin(conf.AdminURL)
	h.SetRoutes(r, &httprouterx.RouterPublic{Router: r.Router}, func(h http.Handler) http.Handler {
		return h
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL + oauth2.DefaultConsentPath)
	assert.Nil(t, err)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	assert.NotEmpty(t, body)
}
