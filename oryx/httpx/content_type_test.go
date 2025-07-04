// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasContentType(t *testing.T) {
	assert.True(t, HasContentType(&http.Request{Header: map[string][]string{}}, "application/octet-stream"))
	assert.False(t, HasContentType(&http.Request{Header: map[string][]string{}}, "not-application/octet-stream"))
	assert.True(t, HasContentType(&http.Request{Header: map[string][]string{"Content-Type": {"application/octet-stream"}}}, "application/octet-stream"))

	// Invalid conent types
	assert.False(t, HasContentType(&http.Request{Header: map[string][]string{"Content-Type": {"application/octet-stream, not-application/application"}}}, "not-application/application"))
	assert.False(t, HasContentType(&http.Request{Header: map[string][]string{"Content-Type": {"application/octet-stream,not-application/application"}}}, "not-application/application"))
	assert.False(t, HasContentType(&http.Request{Header: map[string][]string{"Content-Type": {"application/octet-stream, application/not-application"}}}, "not-application/not-octet-stream"))
	assert.False(t, HasContentType(&http.Request{Header: map[string][]string{"Content-Type": {"a"}}}, "not-application/not-octet-stream"))
}
