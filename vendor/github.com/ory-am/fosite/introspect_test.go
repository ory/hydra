package fosite_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/ory-am/fosite"
	"github.com/ory-am/fosite/compose"
	"github.com/ory-am/fosite/internal"
	"github.com/ory-am/fosite/storage"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestIntrospect(t *testing.T) {
	ctrl := gomock.NewController(t)
	validator := internal.NewMockTokenIntrospector(ctrl)
	defer ctrl.Finish()

	f := compose.ComposeAllEnabled(new(compose.Config), storage.NewMemoryStore(), []byte{}, nil).(*Fosite)
	httpreq := &http.Request{
		Header: http.Header{
			"Authorization": []string{"bearer some-token"},
		},
		Form: url.Values{},
	}

	for k, c := range []struct {
		description string
		scopes      []string
		setup       func()
		expectErr   error
	}{
		{
			description: "should fail",
			scopes:      []string{},
			setup: func() {
			},
			expectErr: ErrRequestUnauthorized,
		},
		{
			description: "should fail",
			scopes:      []string{"foo"},
			setup: func() {
				f.TokenIntrospectionHandlers = TokenIntrospectionHandlers{validator}
				validator.EXPECT().IntrospectToken(nil, "some-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrUnknownRequest)
			},
			expectErr: ErrRequestUnauthorized,
		},
		{
			description: "should fail",
			scopes:      []string{"foo"},
			setup: func() {
				validator.EXPECT().IntrospectToken(nil, "some-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrInvalidClient)
			},
			expectErr: ErrInvalidClient,
		},
		{
			description: "should pass",
			setup: func() {
				validator.EXPECT().IntrospectToken(nil, "some-token", gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx context.Context, _ string, _ TokenType, accessRequest AccessRequester, _ []string) {
					accessRequest.(*AccessRequest).GrantedScopes = []string{"bar"}
				}).Return(nil)
			},
		},
		{
			description: "should pass",
			scopes:      []string{"bar"},
			setup: func() {
				validator.EXPECT().IntrospectToken(nil, "some-token", gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx context.Context, _ string, _ TokenType, accessRequest AccessRequester, _ []string) {
					accessRequest.(*AccessRequest).GrantedScopes = []string{"bar"}
				}).Return(nil)
			},
		},
	} {
		c.setup()
		_, err := f.IntrospectToken(nil, AccessTokenFromRequest(httpreq), AccessToken, nil, c.scopes...)
		assert.True(t, errors.Cause(err) == c.expectErr, "(%d) %s\n%s\n%s", k, c.description, err, c.expectErr)
		t.Logf("Passed test case %d", k)
	}
}
