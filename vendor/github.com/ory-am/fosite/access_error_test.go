package fosite_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/ory-am/fosite"
	. "github.com/ory-am/fosite/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteAccessError(t *testing.T) {
	f := &Fosite{}
	header := http.Header{}
	ctrl := gomock.NewController(t)
	rw := NewMockResponseWriter(ctrl)
	defer ctrl.Finish()

	rw.EXPECT().Header().AnyTimes().Return(header)
	rw.EXPECT().WriteHeader(http.StatusBadRequest)
	rw.EXPECT().Write(gomock.Any())

	f.WriteAccessError(rw, nil, ErrInvalidRequest)
}

func TestWriteAccessError_RFC6749(t *testing.T) {
	// https://tools.ietf.org/html/rfc6749#section-5.2

	f := &Fosite{}

	for k, c := range []struct {
		err  error
		code string
	}{
		{ErrInvalidRequest, "invalid_request"},
		{ErrInvalidClient, "invalid_client"},
		{ErrInvalidGrant, "invalid_grant"},
		{ErrInvalidScope, "invalid_scope"},
		{ErrUnauthorizedClient, "unauthorized_client"},
		{ErrUnsupportedGrantType, "unsupported_grant_type"},
	} {
		rw := httptest.NewRecorder()
		f.WriteAccessError(rw, nil, c.err)

		var params struct {
			Error       string `json:"error"`             // specified by RFC, required
			Description string `json:"error_description"` // specified by RFC, optional
		}

		require.NotNil(t, rw.Body, "(%d) %s: nil body", k, c.code)
		err := json.NewDecoder(rw.Body).Decode(&params)
		require.NoError(t, err, "(%d) %s", k, c.code)

		assert.Equal(t, c.code, params.Error, "(%d) %s: error", k, c.code)
		assert.Equal(t, c.err.Error(), params.Description, "(%d) %s: description", k, c.code)
	}
}
