// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	. "github.com/ory/hydra/v2/fosite"
	. "github.com/ory/hydra/v2/fosite/internal"
)

func TestWriteAccessError(t *testing.T) {
	f := &Fosite{Config: new(Config)}
	header := http.Header{}
	ctrl := gomock.NewController(t)
	rw := NewMockResponseWriter(ctrl)
	t.Cleanup(ctrl.Finish)

	rw.EXPECT().Header().AnyTimes().Return(header)
	rw.EXPECT().WriteHeader(http.StatusBadRequest)
	rw.EXPECT().Write(gomock.Any())

	f.WriteAccessError(context.Background(), rw, nil, ErrInvalidRequest)
}

func TestWriteAccessError_RFC6749(t *testing.T) {
	// https://tools.ietf.org/html/rfc6749#section-5.2

	config := new(Config)
	f := &Fosite{Config: config}

	for k, c := range []struct {
		err                *RFC6749Error
		code               string
		debug              bool
		expectDebugMessage string
		includeExtraFields bool
	}{
		{ErrInvalidRequest.WithDebug("some-debug"), "invalid_request", true, "some-debug", true},
		{ErrInvalidRequest.WithDebugf("some-debug-%d", 1234), "invalid_request", true, "some-debug-1234", true},
		{ErrInvalidRequest.WithDebug("some-debug"), "invalid_request", false, "some-debug", true},
		{ErrInvalidClient.WithDebug("some-debug"), "invalid_client", false, "some-debug", true},
		{ErrInvalidGrant.WithDebug("some-debug"), "invalid_grant", false, "some-debug", true},
		{ErrInvalidScope.WithDebug("some-debug"), "invalid_scope", false, "some-debug", true},
		{ErrUnauthorizedClient.WithDebug("some-debug"), "unauthorized_client", false, "some-debug", true},
		{ErrUnsupportedGrantType.WithDebug("some-debug"), "unsupported_grant_type", false, "some-debug", true},
		{ErrUnsupportedGrantType.WithDebug("some-debug"), "unsupported_grant_type", false, "some-debug", false},
		{ErrUnsupportedGrantType.WithDebug("some-debug"), "unsupported_grant_type", true, "some-debug", false},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			config.SendDebugMessagesToClients = c.debug
			config.UseLegacyErrorFormat = c.includeExtraFields

			rw := httptest.NewRecorder()
			f.WriteAccessError(context.Background(), rw, nil, c.err)

			var params struct {
				Error       string `json:"error"`             // specified by RFC, required
				Description string `json:"error_description"` // specified by RFC, optional
				Debug       string `json:"error_debug"`
				Hint        string `json:"error_hint"`
			}

			require.NotNil(t, rw.Body)
			err := json.NewDecoder(rw.Body).Decode(&params)
			require.NoError(t, err)

			assert.Equal(t, c.code, params.Error)
			if !c.includeExtraFields {
				assert.Empty(t, params.Debug)
				assert.Empty(t, params.Hint)
				assert.Contains(t, params.Description, c.err.DescriptionField)
				assert.Contains(t, params.Description, c.err.HintField)

				if c.debug {
					assert.Contains(t, params.Description, c.err.DebugField)
				} else {
					assert.NotContains(t, params.Description, c.err.DebugField)
				}
			} else {
				assert.EqualValues(t, c.err.DescriptionField, params.Description)
				assert.EqualValues(t, c.err.HintField, params.Hint)

				if !c.debug {
					assert.Empty(t, params.Debug)
				} else {
					assert.EqualValues(t, c.err.DebugField, params.Debug)
				}
			}
		})
	}
}
