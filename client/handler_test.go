/*
 * Copyright © 2015-2018 Javier Viera <javier.viera@mindcurv.com>
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
 * @author		Javier Viera <javier.viera@mindcurv.com>
 * @Copyright 	2017-2018 Javier Viera <javier.viera@mindcurv.com>
 * @license 	Apache-2.0
 */

package client_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/x/urlx"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/internal"
)

func TestValidateDynClientRegistrationAuthorizationBadReq(t *testing.T) {
	reg := internal.NewMockedRegistry(t)
	h := client.NewHandler(reg)
	c := client.Client{OutfacingID: "someid"}
	u := "https://www.something.com"
	hr := &http.Request{Header: map[string][]string{}, URL: urlx.ParseOrPanic(u), RequestURI: u}
	err := h.ValidateDynClientRegistrationAuthorization(hr, c)
	require.EqualValues(t, "The request could not be authorized", err.Error())
}

func TestValidateDynClientRegistrationAuthorizationBadBasicAuthNoBase64(t *testing.T) {
	reg := internal.NewMockedRegistry(t)
	h := client.NewHandler(reg)
	c := client.Client{OutfacingID: "someid"}
	u := "https://www.something.com"
	hr := &http.Request{Header: map[string][]string{}, URL: urlx.ParseOrPanic(u), RequestURI: u}
	hr.Header.Add("Authorization", "Basic something")
	err := h.ValidateDynClientRegistrationAuthorization(hr, c)
	require.EqualValues(t, "The request could not be authorized", err.Error())
}

func TestValidateDynClientRegistrationAuthorizationBadBasicAuthKo(t *testing.T) {
	reg := internal.NewMockedRegistry(t)
	h := client.NewHandler(reg)
	c := client.Client{OutfacingID: "client"}
	u := "https://www.something.com"
	hr := &http.Request{Header: map[string][]string{}, URL: urlx.ParseOrPanic(u), RequestURI: u}
	hr.Header.Add("Authorization", "Basic Y2xpZW50OnNlY3JldA==")
	err := h.ValidateDynClientRegistrationAuthorization(hr, c)
	require.EqualValues(t, "The request could not be authorized", err.Error())
}
