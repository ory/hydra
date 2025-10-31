// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package integration_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/compose"
	"github.com/ory/hydra/v2/fosite/integration/clients"
)

type introspectJWTBearerTokenSuite struct {
	suite.Suite

	clientJWT          *clients.JWTBearer
	clientIntrospect   *clients.Introspect
	clientTokenPayload *clients.JWTBearerPayload
	appTokenPayload    *clients.JWTBearerPayload

	authorizationHeader string
	scopes              []string
	audience            []string
}

func (s *introspectJWTBearerTokenSuite) SetupTest() {
	s.scopes = []string{"fosite"}
	s.audience = []string{tokenURL, "https://example.com"}

	s.clientTokenPayload = &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:   firstJWTBearerIssuer,
			Subject:  firstJWTBearerSubject,
			Audience: s.audience,
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	s.appTokenPayload = &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:   secondJWTBearerIssuer,
			Subject:  secondJWTBearerSubject,
			Audience: s.audience,
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
}

func (s *introspectJWTBearerTokenSuite) TestSuccessResponseWithMultipleScopesToken() {
	ctx := context.Background()

	scopes := []string{"fosite", "docker"}
	token, err := s.getJWTClient().GetToken(ctx, s.clientTokenPayload, scopes)
	require.NoError(s.T(), err)

	response, err := s.clientIntrospect.IntrospectToken(
		ctx,
		clients.IntrospectForm{
			Token:  token.AccessToken,
			Scopes: nil,
		},
		map[string]string{"Authorization": s.authorizationHeader},
	)

	s.assertSuccessResponse(s.T(), response, err, firstJWTBearerSubject)
	assert.Equal(s.T(), strings.Split(response.Scope, " "), scopes)
}

func (s *introspectJWTBearerTokenSuite) TestUnActiveResponseWithInvalidScopes() {
	ctx := context.Background()

	token, err := s.getJWTClient().GetToken(ctx, s.clientTokenPayload, s.scopes)
	require.NoError(s.T(), err)

	response, err := s.clientIntrospect.IntrospectToken(
		ctx,
		clients.IntrospectForm{
			Token:  token.AccessToken,
			Scopes: []string{"invalid"},
		},
		map[string]string{"Authorization": s.authorizationHeader},
	)

	require.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.False(s.T(), response.Active)
}

func (s *introspectJWTBearerTokenSuite) TestSuccessResponseWithoutScopesForIntrospection() {
	ctx := context.Background()

	token, err := s.getJWTClient().GetToken(ctx, s.clientTokenPayload, s.scopes)
	require.NoError(s.T(), err)

	response, err := s.clientIntrospect.IntrospectToken(
		ctx,
		clients.IntrospectForm{
			Token:  token.AccessToken,
			Scopes: nil,
		},
		map[string]string{"Authorization": s.authorizationHeader},
	)

	s.assertSuccessResponse(s.T(), response, err, firstJWTBearerSubject)
}

func (s *introspectJWTBearerTokenSuite) TestSuccessResponseWithoutScopes() {
	ctx := context.Background()

	token, err := s.getJWTClient().GetToken(ctx, s.clientTokenPayload, nil)
	require.NoError(s.T(), err)

	response, err := s.clientIntrospect.IntrospectToken(
		ctx,
		clients.IntrospectForm{
			Token:  token.AccessToken,
			Scopes: nil,
		},
		map[string]string{"Authorization": s.authorizationHeader},
	)

	s.assertSuccessResponse(s.T(), response, err, firstJWTBearerSubject)
}

func (s *introspectJWTBearerTokenSuite) TestSubjectHasAccessToScopeButNotInited() {
	ctx := context.Background()

	token, err := s.getJWTClient().GetToken(ctx, s.clientTokenPayload, nil)
	require.NoError(s.T(), err)

	response, err := s.clientIntrospect.IntrospectToken(
		ctx,
		clients.IntrospectForm{
			Token:  token.AccessToken,
			Scopes: s.scopes,
		},
		map[string]string{"Authorization": s.authorizationHeader},
	)

	require.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.False(s.T(), response.Active)
}

func (s *introspectJWTBearerTokenSuite) TestTheSameTokenInRequestAndHeader() {
	ctx := context.Background()
	token, err := s.getJWTClient().GetToken(ctx, s.clientTokenPayload, s.scopes)
	require.NoError(s.T(), err)

	response, err := s.clientIntrospect.IntrospectToken(
		ctx,
		clients.IntrospectForm{
			Token:  token.AccessToken,
			Scopes: nil,
		},
		map[string]string{"Authorization": "bearer " + token.AccessToken},
	)

	s.assertUnauthorizedResponse(s.T(), response, err)
}

func (s *introspectJWTBearerTokenSuite) TestUnauthorizedResponseForRequestWithoutAuthorization() {
	ctx := context.Background()
	token, err := s.getJWTClient().GetToken(ctx, s.clientTokenPayload, s.scopes)
	require.NoError(s.T(), err)

	response, err := s.clientIntrospect.IntrospectToken(
		ctx,
		clients.IntrospectForm{
			Token:  token.AccessToken,
			Scopes: nil,
		},
		nil,
	)

	s.assertUnauthorizedResponse(s.T(), response, err)
}

func (s *introspectJWTBearerTokenSuite) getJWTClient() *clients.JWTBearer {
	client := *s.clientJWT

	return &client
}

func (s *introspectJWTBearerTokenSuite) assertSuccessResponse(
	t *testing.T,
	response *clients.IntrospectResponse,
	err error,
	subject string,
) {
	assert.Nil(t, err)
	assert.NotNil(t, response)

	assert.True(t, response.Active)
	assert.Equal(t, response.Subject, subject)
	assert.NotEmpty(t, response.ExpiresAt)
	assert.NotEmpty(t, response.IssuedAt)
	assert.Equal(t, response.Audience, s.audience)

	tokenDuration := time.Unix(response.ExpiresAt, 0).Sub(time.Unix(response.IssuedAt, 0))
	assert.Less(t, int64(tokenDuration), int64(time.Hour+time.Minute))
	assert.Greater(t, int64(tokenDuration), int64(time.Hour-time.Minute))
}

func (s *introspectJWTBearerTokenSuite) assertUnauthorizedResponse(
	t *testing.T,
	response *clients.IntrospectResponse,
	err error,
) {
	assert.Nil(t, response)
	assert.NotNil(t, err)

	retrieveError, ok := err.(*clients.RequestError)
	assert.True(t, ok)
	assert.Equal(t, retrieveError.Response.StatusCode, http.StatusUnauthorized)
}

func TestIntrospectJWTBearerTokenSuite(t *testing.T) {
	provider := compose.Compose(
		&fosite.Config{
			GrantTypeJWTBearerCanSkipClientAuth:  true,
			GrantTypeJWTBearerIDOptional:         true,
			GrantTypeJWTBearerIssuedDateOptional: true,
			AccessTokenLifespan:                  time.Hour,
			TokenURL:                             tokenURL,
		},
		fositeStore,
		jwtStrategy,
		compose.OAuth2ClientCredentialsGrantFactory,
		compose.RFC7523AssertionGrantFactory,
		compose.OAuth2TokenIntrospectionFactory,
	)
	testServer := mockServer(t, provider, &fosite.DefaultSession{})
	defer testServer.Close()

	client := newJWTBearerAppClient(testServer)
	if err := client.SetPrivateKey(secondKeyID, secondPrivateKey); err != nil {
		assert.Nil(t, err)
	}

	token, err := client.GetToken(context.Background(), &clients.JWTBearerPayload{
		Claims: &jwt.Claims{
			Issuer:   secondJWTBearerIssuer,
			Subject:  secondJWTBearerSubject,
			Audience: []string{tokenURL},
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}, []string{"fosite"})
	if err != nil {
		assert.Nil(t, err)
	}

	if err := client.SetPrivateKey(firstKeyID, firstPrivateKey); err != nil {
		assert.Nil(t, err)
	}

	suite.Run(t, &introspectJWTBearerTokenSuite{
		clientJWT:           client,
		clientIntrospect:    clients.NewIntrospectClient(testServer.URL + "/introspect"),
		authorizationHeader: "bearer " + token.AccessToken,
	})
}
