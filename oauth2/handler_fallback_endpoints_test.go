/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

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
