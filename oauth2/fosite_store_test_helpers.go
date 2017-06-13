package oauth2

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var defaultRequest = fosite.Request{
	RequestedAt:   time.Now().Round(time.Second),
	Client:        &client.Client{ID: "foobar"},
	Scopes:        fosite.Arguments{"fa", "ba"},
	GrantedScopes: fosite.Arguments{"fa", "ba"},
	Form:          url.Values{"foo": []string{"bar", "baz"}},
	Session:       &fosite.DefaultSession{Subject: "bar"},
}

func TestHelperCreateGetDeleteAuthorizeCodes(m pkg.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		_, err := m.GetAuthorizeCodeSession(ctx, "4321", &fosite.DefaultSession{})
		assert.NotNil(t, err)

		err = m.CreateAuthorizeCodeSession(ctx, "4321", &defaultRequest)
		require.NoError(t, err)

		res, err := m.GetAuthorizeCodeSession(ctx, "4321", &fosite.DefaultSession{})
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "Scopes", "GrantedScopes", "Form", "Session")

		err = m.DeleteAuthorizeCodeSession(ctx, "4321")
		require.NoError(t, err)

		_, err = m.GetAuthorizeCodeSession(ctx, "4321", &fosite.DefaultSession{})
		assert.NotNil(t, err)
	}
}

func TestHelperCreateGetDeleteAccessTokenSession(m pkg.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		_, err := m.GetAccessTokenSession(ctx, "4321", &fosite.DefaultSession{})
		assert.NotNil(t, err)

		err = m.CreateAccessTokenSession(ctx, "4321", &defaultRequest)
		require.NoError(t, err)

		res, err := m.GetAccessTokenSession(ctx, "4321", &fosite.DefaultSession{})
		require.NoError(t, err)
		AssertObjectKeysEqual(t, &defaultRequest, res, "Scopes", "GrantedScopes", "Form", "Session")

		err = m.DeleteAccessTokenSession(ctx, "4321")
		require.NoError(t, err)

		_, err = m.GetAccessTokenSession(ctx, "4321", &fosite.DefaultSession{})
		assert.NotNil(t, err)
	}
}
