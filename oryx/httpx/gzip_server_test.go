// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import (
	"bytes"
	gzip2 "compress/gzip"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gobuffalo/httptest"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/negroni"
)

func makeRequest(t *testing.T, data string, ts *httptest.Server) {
	var buf bytes.Buffer
	gzip := gzip2.NewWriter(&buf)

	_, err := gzip.Write([]byte(data))
	require.NoError(t, err)
	require.NoError(t, gzip.Close())

	c := http.Client{}
	req, err := http.NewRequest("POST", ts.URL, &buf)
	req.Header.Set("Content-Encoding", "gzip")
	require.NoError(t, err)
	res, err := c.Do(req)
	require.NoError(t, err)
	res.Body.Close()
	assert.EqualValues(t, http.StatusNoContent, res.StatusCode)
}

func TestGZipServer(t *testing.T) {
	router := httprouter.New()
	router.POST("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var f json.RawMessage
		require.NoError(t, json.NewDecoder(r.Body).Decode(&f))
		t.Logf("%s", f)
		w.WriteHeader(http.StatusNoContent)
	})
	n := negroni.New(NewCompressionRequestReader(func(w http.ResponseWriter, r *http.Request, err error) {
		require.NoError(t, err)
	}))
	n.UseHandler(router)
	ts := httptest.NewServer(n)
	defer ts.Close()

	makeRequest(t, "true", ts)
}
