// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package serverx

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test404Handler(t *testing.T) {
	router := httprouter.New()
	router.NotFound = DefaultNotFoundHandler
	ts := httptest.NewServer(router)
	t.Cleanup(ts.Close)

	for k, tc := range []struct {
		accept              string
		expectedBody        string
		expectedContentType string
	}{
		{
			accept:              "",
			expectedBody:        string(page404HTML),
			expectedContentType: "text/html; charset=utf-8",
		},
		{
			accept:              "text/html",
			expectedBody:        string(page404HTML),
			expectedContentType: "text/html; charset=utf-8",
		},
		{
			accept:              "text/*",
			expectedBody:        string(page404HTML),
			expectedContentType: "text/html; charset=utf-8",
		},
		{
			accept:              "application/json",
			expectedBody:        string(page404JSON),
			expectedContentType: "application/json; charset=utf-8",
		},
		{
			accept:              "text/plain",
			expectedBody:        `Error 404 - The requested route does not exist. Make sure you are using the right path, domain, and port.`,
			expectedContentType: "text/plain; charset=utf-8",
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			req, err := http.NewRequest("GET", ts.URL+"/404", nil)
			require.NoError(t, err)
			req.Header.Set("Accept", tc.accept)
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
			assert.Equal(t, tc.expectedContentType, resp.Header.Get("Content-Type"))
			body := make([]byte, len(tc.expectedBody))
			_, err = io.ReadFull(resp.Body, body)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedBody, string(body))
		})
	}
}
