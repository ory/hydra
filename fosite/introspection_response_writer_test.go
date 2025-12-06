// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ory/x/errorsx"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	. "github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/internal"
)

func TestWriteIntrospectionError(t *testing.T) {
	f := &Fosite{Config: new(Config)}
	c := gomock.NewController(t)
	defer c.Finish()

	rw := internal.NewMockResponseWriter(c)
	rw.EXPECT().WriteHeader(http.StatusUnauthorized)
	rw.EXPECT().Header().AnyTimes().Return(http.Header{})
	rw.EXPECT().Write(gomock.Any())
	f.WriteIntrospectionError(context.Background(), rw, errorsx.WithStack(ErrRequestUnauthorized))

	rw.EXPECT().WriteHeader(http.StatusBadRequest)
	rw.EXPECT().Write(gomock.Any())
	f.WriteIntrospectionError(context.Background(), rw, errorsx.WithStack(ErrInvalidRequest))

	rw.EXPECT().Write([]byte("{\"active\":false}\n"))
	f.WriteIntrospectionError(context.Background(), rw, errors.New(""))

	rw.EXPECT().Write([]byte("{\"active\":false}\n"))
	f.WriteIntrospectionError(context.Background(), rw, errorsx.WithStack(ErrInactiveToken.WithWrap(ErrRequestUnauthorized)))

	f.WriteIntrospectionError(context.Background(), rw, nil)
}

func TestWriteIntrospectionResponse(t *testing.T) {
	f := new(Fosite)
	c := gomock.NewController(t)
	defer c.Finish()

	rw := internal.NewMockResponseWriter(c)
	rw.EXPECT().Write(gomock.Any()).AnyTimes()
	rw.EXPECT().Header().AnyTimes().Return(http.Header{})
	f.WriteIntrospectionResponse(context.Background(), rw, &IntrospectionResponse{
		AccessRequester: NewAccessRequest(nil),
	})
}

func TestWriteIntrospectionResponseBody(t *testing.T) {
	f := new(Fosite)
	ires := &IntrospectionResponse{}
	rw := httptest.NewRecorder()

	for _, c := range []struct {
		description string
		setup       func()
		active      bool
		hasExp      bool
		hasExtra    bool
	}{
		{
			description: "should success for not expired access token",
			setup: func() {
				ires.Active = true
				ires.TokenUse = AccessToken
				sess := &DefaultSession{}
				sess.SetExpiresAt(ires.TokenUse, time.Now().Add(time.Hour*2))
				ires.AccessRequester = NewAccessRequest(sess)
			},
			active:   true,
			hasExp:   true,
			hasExtra: false,
		},
		{
			description: "should success for expired access token",
			setup: func() {
				ires.Active = false
				ires.TokenUse = AccessToken
				sess := &DefaultSession{}
				sess.SetExpiresAt(ires.TokenUse, time.Now().Add(-time.Hour*2))
				ires.AccessRequester = NewAccessRequest(sess)
			},
			active:   false,
			hasExp:   false,
			hasExtra: false,
		},
		{
			description: "should success for ExpiresAt not set access token",
			setup: func() {
				ires.Active = true
				ires.TokenUse = AccessToken
				sess := &DefaultSession{}
				sess.SetExpiresAt(ires.TokenUse, time.Time{})
				ires.AccessRequester = NewAccessRequest(sess)
			},
			active:   true,
			hasExp:   false,
			hasExtra: false,
		},
		{
			description: "should output extra claims",
			setup: func() {
				ires.Active = true
				ires.TokenUse = AccessToken
				sess := &DefaultSession{}
				sess.GetExtraClaims()["extra"] = "foobar"
				// We try to set these, but they should be ignored.
				for _, field := range []string{"exp", "client_id", "scope", "iat", "sub", "aud", "username"} {
					sess.GetExtraClaims()[field] = "invalid"
				}
				sess.SetExpiresAt(ires.TokenUse, time.Time{})
				ires.AccessRequester = NewAccessRequest(sess)
			},
			active:   true,
			hasExp:   false,
			hasExtra: true,
		},
	} {
		t.Run(c.description, func(t *testing.T) {
			c.setup()
			f.WriteIntrospectionResponse(context.Background(), rw, ires)
			var params struct {
				Active   bool   `json:"active"`
				Exp      *int64 `json:"exp"`
				Iat      *int64 `json:"iat"`
				Extra    string `json:"extra"`
				ClientId string `json:"client_id"`
				Scope    string `json:"scope"`
				Subject  string `json:"sub"`
				Audience string `json:"aud"`
				Username string `json:"username"`
			}
			assert.Equal(t, 200, rw.Code)
			err := json.NewDecoder(rw.Body).Decode(&params)
			require.NoError(t, err)
			assert.Equal(t, c.active, params.Active)
			if c.active {
				assert.NotNil(t, params.Iat)
				if c.hasExp {
					assert.NotNil(t, params.Exp)
				} else {
					assert.Nil(t, params.Exp)
				}
				if c.hasExtra {
					assert.Equal(t, params.Extra, "foobar")
				} else {
					assert.Empty(t, params.Extra)
				}
				assert.NotEqual(t, "invalid", params.Exp)
				assert.NotEqual(t, "invalid", params.ClientId)
				assert.NotEqual(t, "invalid", params.Scope)
				assert.NotEqual(t, "invalid", params.Iat)
				assert.NotEqual(t, "invalid", params.Subject)
				assert.NotEqual(t, "invalid", params.Audience)
				assert.NotEqual(t, "invalid", params.Username)
			}
		})
	}
}
