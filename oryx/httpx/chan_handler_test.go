// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChanHandler(t *testing.T) {
	h, c := NewChanHandler(1)
	s := httptest.NewServer(h)

	c <- func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(555)
	}
	resp, err := s.Client().Get(s.URL)
	require.NoError(t, err)
	assert.Equal(t, 555, resp.StatusCode)

	c <- func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(337)
	}
	resp, err = s.Client().Get(s.URL)
	require.NoError(t, err)
	assert.Equal(t, 337, resp.StatusCode)
}
