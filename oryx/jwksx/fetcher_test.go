// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwksx

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	keys = `{
  "keys": [
    {
      "use": "sig",
      "kty": "oct",
      "kid": "7d5f5ad0674ec2f2960b1a34f33370a0f71471fa0e3ef0c0a692977d276dafe8",
      "alg": "HS256",
      "k": "Y2hhbmdlbWVjaGFuZ2VtZWNoYW5nZW1lY2hhbmdlbWU"
    }
  ]
}`
	secret = "changemechangemechangemechangeme"
)

func TestFetcher(t *testing.T) {
	var called int
	var h http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		called++
		w.Write([]byte(keys))
	}
	ts := httptest.NewServer(h)
	defer ts.Close()

	f := NewFetcher(ts.URL)

	k, err := f.GetKey("7d5f5ad0674ec2f2960b1a34f33370a0f71471fa0e3ef0c0a692977d276dafe8")
	require.NoError(t, err)
	assert.EqualValues(t, secret, fmt.Sprintf("%s", k.Key))
	assert.Equal(t, 1, called)

	k, err = f.GetKey("7d5f5ad0674ec2f2960b1a34f33370a0f71471fa0e3ef0c0a692977d276dafe8")
	require.NoError(t, err)
	assert.EqualValues(t, secret, fmt.Sprintf("%s", k.Key))
	assert.Equal(t, 1, called)

	_, err = f.GetKey("does-not-exist")
	require.Error(t, err)
	assert.Equal(t, 2, called)
}
