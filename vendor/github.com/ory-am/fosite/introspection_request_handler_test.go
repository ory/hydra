package fosite_test

import (
	"github.com/golang/mock/gomock"
	"github.com/ory-am/fosite"
	. "github.com/ory-am/fosite"
	"github.com/ory-am/fosite/compose"
	"github.com/ory-am/fosite/internal"
	"github.com/ory-am/fosite/storage"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestIntrospectionResponse(t *testing.T) {
	r := &fosite.IntrospectionResponse{
		AccessRequester: fosite.NewAccessRequest(nil),
		Active:          true,
	}

	assert.Equal(t, r.AccessRequester, r.GetAccessRequester())
	assert.Equal(t, r.Active, r.IsActive())
}

func TestNewIntrospectionRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	validator := internal.NewMockTokenIntrospector(ctrl)
	defer ctrl.Finish()

	f := compose.ComposeAllEnabled(new(compose.Config), storage.NewMemoryStore(), []byte{}, nil).(*Fosite)
	httpreq := &http.Request{
		Method: "POST",
		Header: http.Header{},
		Form:   url.Values{},
	}
	newErr := errors.New("")

	for k, c := range []struct {
		description string
		setup       func()
		expectErr   error
		isActive    bool
	}{
		{
			description: "should fail",
			setup: func() {
			},
			expectErr: ErrInvalidRequest,
		},
		{
			description: "should fail",
			setup: func() {
				f.TokenIntrospectionHandlers = TokenIntrospectionHandlers{validator}
				httpreq = &http.Request{
					Method: "POST",
					Header: http.Header{
						"Authorization": []string{"bearer some-token"},
					},
					PostForm: url.Values{
						"token": []string{"introspect-token"},
					},
				}
				validator.EXPECT().IntrospectToken(nil, "some-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				validator.EXPECT().IntrospectToken(nil, "introspect-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(newErr)
			},
			isActive:  false,
			expectErr: newErr,
		},
		{
			description: "should pass",
			setup: func() {
				f.TokenIntrospectionHandlers = TokenIntrospectionHandlers{validator}
				httpreq = &http.Request{
					Method: "POST",
					Header: http.Header{
						"Authorization": []string{"bearer some-token"},
					},
					PostForm: url.Values{
						"token": []string{"introspect-token"},
					},
				}
				validator.EXPECT().IntrospectToken(nil, "some-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				validator.EXPECT().IntrospectToken(nil, "introspect-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			isActive: true,
		},
	} {
		c.setup()
		res, err := f.NewIntrospectionRequest(nil, httpreq, &DefaultSession{})
		assert.True(t, errors.Cause(err) == c.expectErr, "(%d) %s\n%s\n%s", k, c.description, err, c.expectErr)
		if res != nil {
			assert.Equal(t, c.isActive, res.IsActive())
		}
		t.Logf("Passed test case %d", k)
	}
}
