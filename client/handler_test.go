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
	"bytes"
	"net/http"
	"net/http/httptest"
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

func TestCreateOk(t *testing.T) {
	u := "https://www.something.com"
	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()
	reg := internal.NewMockedRegistry(t)
	h := client.NewHandler(reg)
	h.Create(rr, req, nil)
	require.EqualValues(t, http.StatusCreated, rr.Result().StatusCode)
}

func TestCreateDynamicRegistrationOk(t *testing.T) {
	u := "https://www.something.com"
	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()
	reg := internal.NewMockedRegistry(t)
	h := client.NewHandler(reg)
	h.CreateDynamicRegistration(rr, req, nil)
	require.EqualValues(t, http.StatusCreated, rr.Result().StatusCode)
}

func TestUpdateKo(t *testing.T) {
	u := "https://www.something.com"
	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()
	reg := internal.NewMockedRegistry(t)
	h := client.NewHandler(reg)
	h.Update(rr, req, nil)
	require.EqualValues(t, http.StatusNotFound, rr.Result().StatusCode)
}

func TestPatchKo(t *testing.T) {
	u := "https://www.something.com"
	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()
	reg := internal.NewMockedRegistry(t)
	h := client.NewHandler(reg)
	h.Patch(rr, req, nil)
	require.EqualValues(t, http.StatusNotFound, rr.Result().StatusCode)
}

func TestUpdateDynamicRegistrationKo(t *testing.T) {
	u := "https://www.something.com"
	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()
	reg := internal.NewMockedRegistry(t)
	h := client.NewHandler(reg)
	h.UpdateDynamicRegistration(rr, req, nil)
	require.EqualValues(t, http.StatusUnauthorized, rr.Result().StatusCode)
}

func TestListOk(t *testing.T) {
	u := "https://www.something.com"
	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()
	reg := internal.NewMockedRegistry(t)
	h := client.NewHandler(reg)
	h.List(rr, req, nil)
	require.EqualValues(t, http.StatusOK, rr.Result().StatusCode)
}

func TestGetKo(t *testing.T) {
	u := "https://www.something.com"
	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()
	reg := internal.NewMockedRegistry(t)
	h := client.NewHandler(reg)
	h.Get(rr, req, nil)
	require.EqualValues(t, http.StatusUnauthorized, rr.Result().StatusCode)
}

func TestGetDynamicRegistrationKo(t *testing.T) {
	u := "https://www.something.com"
	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()
	reg := internal.NewMockedRegistry(t)
	h := client.NewHandler(reg)
	h.GetDynamicRegistration(rr, req, nil)
	require.EqualValues(t, http.StatusUnauthorized, rr.Result().StatusCode)
}

func TestDeleteKo(t *testing.T) {
	u := "https://www.something.com"
	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()
	reg := internal.NewMockedRegistry(t)
	h := client.NewHandler(reg)
	h.Delete(rr, req, nil)
	require.EqualValues(t, http.StatusNotFound, rr.Result().StatusCode)
}

func TestDeleteDynamicRegistrationKo(t *testing.T) {
	u := "https://www.something.com"
	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()
	reg := internal.NewMockedRegistry(t)
	h := client.NewHandler(reg)
	h.DeleteDynamicRegistration(rr, req, nil)
	require.EqualValues(t, http.StatusUnauthorized, rr.Result().StatusCode)
}
