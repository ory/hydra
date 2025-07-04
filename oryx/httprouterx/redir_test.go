// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httprouterx_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	x "github.com/ory/x/httprouterx"
	"github.com/ory/x/urlx"
)

func TestRedirectToPublicAdminRoute(t *testing.T) {
	var ts *httptest.Server
	router := x.NewRouterAdminWithPrefix("/admin", func(ctx context.Context) *url.URL {
		return urlx.ParseOrPanic(ts.URL)
	})
	ts = httptest.NewServer(router)
	t.Cleanup(ts.Close)

	router.POST("/privileged", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		body, _ := io.ReadAll(r.Body)
		w.Write(body)
	})

	router.POST("/read", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		body, _ := io.ReadAll(r.Body)
		w.Write(body)
	})

	for _, tc := range []struct {
		source string
		dest   string
	}{
		{
			source: ts.URL + "/admin/privileged?foo=bar",
			dest:   ts.URL + "/admin/privileged?foo=bar",
		},
		{
			source: ts.URL + "/privileged?foo=bar",
			dest:   ts.URL + "/admin/privileged?foo=bar",
		},
	} {
		t.Run(fmt.Sprintf("source=%s", tc.source), func(t *testing.T) {
			id := uuid.Must(uuid.NewV4()).String()
			res, err := ts.Client().Post(tc.source, "", strings.NewReader(id))
			require.NoError(t, err)
			assert.EqualValues(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, tc.dest, res.Request.URL.String())
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, id, string(body))
		})
	}
}
