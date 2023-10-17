// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/ory/fosite/compose"
	hydra "github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/fositex"
	"github.com/ory/hydra/v2/internal"
	"github.com/ory/hydra/v2/internal/kratos"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/x/contextx"
)

func TestResourceOwnerPasswordGrant(t *testing.T) {
	ctx := context.Background()
	fakeKratos := kratos.NewFake()
	reg := internal.NewMockedRegistry(t, &contextx.Default{})
	reg.WithKratos(fakeKratos)
	reg.WithExtraFositeFactories([]fositex.Factory{compose.OAuth2ResourceOwnerPasswordCredentialsFactory})
	_, adminTS := testhelpers.NewOAuth2Server(ctx, t, reg)

	secret := uuid.New().String()
	client := &hydra.Client{
		Secret:     secret,
		GrantTypes: []string{"password"},
	}
	require.NoError(t, reg.ClientManager().CreateClient(ctx, client))

	oauth2Config := &oauth2.Config{
		ClientID:     client.GetID(),
		ClientSecret: secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   reg.Config().OAuth2AuthURL(ctx).String(),
			TokenURL:  reg.Config().OAuth2TokenURL(ctx).String(),
			AuthStyle: oauth2.AuthStyleInHeader,
		},
	}

	t.Run("case=get ROP grant token with valid username and password", func(t *testing.T) {
		token, err := oauth2Config.PasswordCredentialsToken(ctx, kratos.FakeUsername, kratos.FakePassword)
		require.NoError(t, err)
		require.NotEmpty(t, token.AccessToken)
		i := testhelpers.IntrospectToken(t, oauth2Config, token.AccessToken, adminTS)
		assert.True(t, i.Get("active").Bool(), "%s", i)
		assert.EqualValues(t, oauth2Config.ClientID, i.Get("client_id").String(), "%s", i)
	})

	t.Run("case=access denied for invalid password", func(t *testing.T) {
		_, err := oauth2Config.PasswordCredentialsToken(ctx, kratos.FakeUsername, "invalid")
		retrieveError := new(oauth2.RetrieveError)
		require.Error(t, err)
		require.ErrorAs(t, err, &retrieveError)
		assert.Contains(t, retrieveError.ErrorDescription, "Unable to authenticate the provided username and password credentials")
	})
}
