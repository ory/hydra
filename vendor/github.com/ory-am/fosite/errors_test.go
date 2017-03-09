package fosite

import (
	native "errors"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrorToRFC6749(t *testing.T) {
	assert.Equal(t, UnknownErrorName, ErrorToRFC6749Error(errors.New("")).Name)
	assert.Equal(t, UnknownErrorName, ErrorToRFC6749Error(native.New("")).Name)

	assert.Equal(t, errInvalidRequestName, ErrorToRFC6749Error(errors.WithStack(ErrInvalidRequest)).Name)
	assert.Equal(t, errUnauthorizedClientName, ErrorToRFC6749Error(errors.WithStack(ErrUnauthorizedClient)).Name)
	assert.Equal(t, errAccessDeniedName, ErrorToRFC6749Error(errors.WithStack(ErrAccessDenied)).Name)
	assert.Equal(t, errUnsupportedResponseTypeName, ErrorToRFC6749Error(errors.WithStack(ErrUnsupportedResponseType)).Name)
	assert.Equal(t, errInvalidScopeName, ErrorToRFC6749Error(errors.WithStack(ErrInvalidScope)).Name)
	assert.Equal(t, errServerErrorName, ErrorToRFC6749Error(errors.WithStack(ErrServerError)).Name)
	assert.Equal(t, errTemporarilyUnavailableName, ErrorToRFC6749Error(errors.WithStack(ErrTemporarilyUnavailable)).Name)
	assert.Equal(t, errUnsupportedGrantTypeName, ErrorToRFC6749Error(errors.WithStack(ErrUnsupportedGrantType)).Name)
	assert.Equal(t, errInvalidGrantName, ErrorToRFC6749Error(errors.WithStack(ErrInvalidGrant)).Name)
	assert.Equal(t, errInvalidClientName, ErrorToRFC6749Error(errors.WithStack(ErrInvalidClient)).Name)
	assert.Equal(t, errInvalidState, ErrorToRFC6749Error(errors.WithStack(ErrInvalidState)).Name)
}
