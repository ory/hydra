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

package oauth2

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ory/fosite"
)

func TestHandlerConsent(t *testing.T) {
	h := &Handler{
		L:             logrus.New(),
		ScopeStrategy: fosite.HierarchicScopeStrategy,
	}
	r := httprouter.New()
	h.SetRoutes(r, r, func(h http.Handler) http.Handler {
		return h
	})
	ts := httptest.NewServer(r)

	res, err := http.Get(ts.URL + DefaultConsentPath)
	assert.Nil(t, err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	assert.NotEmpty(t, body)
}
