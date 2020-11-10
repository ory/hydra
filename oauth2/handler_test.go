/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package oauth2_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/hydra/internal/httpclient/models"

	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/x"

	"github.com/ory/viper"

	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/internal"
	"github.com/ory/x/urlx"

	jwt2 "github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
	"github.com/ory/hydra/client"
	hydra "github.com/ory/hydra/internal/httpclient/client"
	"github.com/ory/hydra/oauth2"
)

var lifespan = time.Hour
var flushRequests = []*fosite.Request{
	{
		ID:             "flush-1",
		RequestedAt:    time.Now().Round(time.Second),
		Client:         &client.Client{OutfacingID: "foobar"},
		RequestedScope: fosite.Arguments{"fa", "ba"},
		GrantedScope:   fosite.Arguments{"fa", "ba"},
		Form:           url.Values{"foo": []string{"bar", "baz"}},
		Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
	},
	{
		ID:             "flush-2",
		RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
		Client:         &client.Client{OutfacingID: "foobar"},
		RequestedScope: fosite.Arguments{"fa", "ba"},
		GrantedScope:   fosite.Arguments{"fa", "ba"},
		Form:           url.Values{"foo": []string{"bar", "baz"}},
		Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
	},
	{
		ID:             "flush-3",
		RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
		Client:         &client.Client{OutfacingID: "foobar"},
		RequestedScope: fosite.Arguments{"fa", "ba"},
		GrantedScope:   fosite.Arguments{"fa", "ba"},
		Form:           url.Values{"foo": []string{"bar", "baz"}},
		Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
	},
}

func TestHandlerDeleteHandler(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	viper.Set(configuration.ViperKeyIssuerURL, "http://hydra.localhost")
	reg := internal.NewRegistryMemory(t, conf)

	cm := reg.ClientManager()
	store := reg.OAuth2Storage()

	h := oauth2.NewHandler(reg, conf)

	deleteRequest := &fosite.Request{
		ID:             "del-1",
		RequestedAt:    time.Now().Round(time.Second),
		Client:         &client.Client{OutfacingID: "foobar"},
		RequestedScope: fosite.Arguments{"fa", "ba"},
		GrantedScope:   fosite.Arguments{"fa", "ba"},
		Form:           url.Values{"foo": []string{"bar", "baz"}},
		Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
	}
	require.NoError(t, cm.CreateClient(context.Background(), deleteRequest.Client.(*client.Client)))
	require.NoError(t, store.CreateAccessTokenSession(context.Background(), deleteRequest.ID, deleteRequest))

	r := x.NewRouterAdmin()
	h.SetRoutes(r, r.RouterPublic(), func(h http.Handler) http.Handler {
		return h
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	c := hydra.NewHTTPClientWithConfig(nil, &hydra.TransportConfig{Schemes: []string{"http"}, Host: urlx.ParseOrPanic(ts.URL).Host})
	_, err := c.Admin.DeleteOAuth2Token(admin.NewDeleteOAuth2TokenParams().WithClientID("foobar"))
	require.NoError(t, err)

	ds := new(oauth2.Session)
	ctx := context.Background()
	_, err = store.GetAccessTokenSession(ctx, "del-1", ds)
	require.Error(t, err, "not_found")
}

func TestHandlerFlushHandler(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	viper.Set(configuration.ViperKeyScopeStrategy, "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY")
	viper.Set(configuration.ViperKeyIssuerURL, "http://hydra.localhost")
	reg := internal.NewRegistryMemory(t, conf)

	cl := reg.ClientManager()
	store := reg.OAuth2Storage()

	h := oauth2.NewHandler(reg, conf)
	for _, r := range flushRequests {
		_ = cl.CreateClient(context.Background(), r.Client.(*client.Client))
		require.NoError(t, store.CreateAccessTokenSession(context.Background(), r.ID, r))
	}

	r := x.NewRouterAdmin()
	h.SetRoutes(r, r.RouterPublic(), func(h http.Handler) http.Handler {
		return h
	})
	ts := httptest.NewServer(r)
	defer ts.Close()
	c := hydra.NewHTTPClientWithConfig(nil, &hydra.TransportConfig{Schemes: []string{"http"}, Host: urlx.ParseOrPanic(ts.URL).Host})

	ds := new(oauth2.Session)
	ctx := context.Background()

	_, err := c.Admin.FlushInactiveOAuth2Tokens(admin.NewFlushInactiveOAuth2TokensParams().WithBody(&models.FlushInactiveOAuth2TokensRequest{NotAfter: strfmt.DateTime(time.Now().Add(-time.Hour * 24))}))
	require.NoError(t, err)

	_, err = store.GetAccessTokenSession(ctx, "flush-1", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-2", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-3", ds)
	require.NoError(t, err)

	_, err = c.Admin.FlushInactiveOAuth2Tokens(admin.NewFlushInactiveOAuth2TokensParams().WithBody(&models.FlushInactiveOAuth2TokensRequest{NotAfter: strfmt.DateTime(time.Now().Add(-(lifespan + time.Hour/2)))}))
	require.NoError(t, err)

	_, err = store.GetAccessTokenSession(ctx, "flush-1", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-2", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-3", ds)
	require.Error(t, err)

	_, err = c.Admin.FlushInactiveOAuth2Tokens(admin.NewFlushInactiveOAuth2TokensParams().WithBody(&models.FlushInactiveOAuth2TokensRequest{NotAfter: strfmt.DateTime(time.Now())}))
	require.NoError(t, err)

	_, err = store.GetAccessTokenSession(ctx, "flush-1", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-2", ds)
	require.Error(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-3", ds)
	require.Error(t, err)
}

func TestUserinfo(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	viper.Set(configuration.ViperKeyScopeStrategy, "")
	viper.Set(configuration.ViperKeyAuthCodeLifespan, lifespan)
	viper.Set(configuration.ViperKeyIssuerURL, "http://hydra.localhost")
	reg := internal.NewRegistryMemory(t, conf)
	internal.MustEnsureRegistryKeys(reg, x.OpenIDConnectKeyName)

	ctrl := gomock.NewController(t)
	op := NewMockOAuth2Provider(ctrl)
	defer ctrl.Finish()
	reg.WithOAuth2Provider(op)

	h := reg.OAuth2Handler()

	router := x.NewRouterAdmin()
	h.SetRoutes(router, router.RouterPublic(), func(h http.Handler) http.Handler {
		return h
	})
	ts := httptest.NewServer(router)
	defer ts.Close()

	for k, tc := range []struct {
		setup                func(t *testing.T)
		checkForSuccess      func(t *testing.T, body []byte)
		checkForUnauthorized func(t *testing.T, body []byte, header http.Header)
		expectStatusCode     int
	}{
		{
			setup: func(t *testing.T) {
				op.EXPECT().IntrospectToken(gomock.Any(), gomock.Eq("access-token"), gomock.Eq(fosite.AccessToken), gomock.Any()).Return(fosite.AccessToken, nil, errors.New("asdf"))
			},
			expectStatusCode: http.StatusInternalServerError,
		},
		{
			setup: func(t *testing.T) {
				op.EXPECT().
					IntrospectToken(gomock.Any(), gomock.Eq("access-token"), gomock.Eq(fosite.AccessToken), gomock.Any()).
					Return(fosite.RefreshToken, nil, nil)
			},
			checkForUnauthorized: func(t *testing.T, body []byte, headers http.Header) {
				assert.True(t, headers.Get("WWW-Authenticate") != "", "%s", headers)
			},
			expectStatusCode: http.StatusUnauthorized,
		},
		{
			setup: func(t *testing.T) {
				op.EXPECT().
					IntrospectToken(gomock.Any(), gomock.Eq("access-token"), gomock.Eq(fosite.AccessToken), gomock.Any()).
					Return(fosite.AccessToken, nil, fosite.ErrRequestUnauthorized)
			},
			checkForUnauthorized: func(t *testing.T, body []byte, headers http.Header) {
				assert.True(t, headers.Get("WWW-Authenticate") != "", "%s", headers)
			},
			expectStatusCode: http.StatusUnauthorized,
		},
		{
			setup: func(t *testing.T) {
				op.EXPECT().
					IntrospectToken(gomock.Any(), gomock.Eq("access-token"), gomock.Eq(fosite.AccessToken), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ fosite.TokenType, session fosite.Session, _ ...string) (fosite.TokenType, fosite.AccessRequester, error) {
						session = &oauth2.Session{
							DefaultSession: &openid.DefaultSession{
								Claims: &jwt.IDTokenClaims{
									Subject: "alice",
								},
								Headers: new(jwt.Headers),
								Subject: "alice",
							},
							Extra: map[string]interface{}{},
						}

						return fosite.AccessToken, &fosite.AccessRequest{
							Request: fosite.Request{
								Client:  &client.Client{},
								Session: session,
							},
						}, nil
					})
			},
			expectStatusCode: http.StatusOK,
			checkForSuccess: func(t *testing.T, body []byte) {
				assert.True(t, strings.Contains(string(body), `"sub":"alice"`), "%s", body)
			},
		},
		{
			setup: func(t *testing.T) {
				op.EXPECT().
					IntrospectToken(gomock.Any(), gomock.Eq("access-token"), gomock.Eq(fosite.AccessToken), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ fosite.TokenType, session fosite.Session, _ ...string) (fosite.TokenType, fosite.AccessRequester, error) {
						session = &oauth2.Session{
							DefaultSession: &openid.DefaultSession{
								Claims: &jwt.IDTokenClaims{
									Subject: "another-alice",
								},
								Headers: new(jwt.Headers),
								Subject: "alice",
							},
							Extra: map[string]interface{}{},
						}

						return fosite.AccessToken, &fosite.AccessRequest{
							Request: fosite.Request{
								Client:  &client.Client{},
								Session: session,
							},
						}, nil
					})
			},
			expectStatusCode: http.StatusOK,
			checkForSuccess: func(t *testing.T, body []byte) {
				assert.False(t, strings.Contains(string(body), `"sub":"alice"`), "%s", body)
				assert.True(t, strings.Contains(string(body), `"sub":"another-alice"`), "%s", body)
			},
		},
		{
			setup: func(t *testing.T) {
				op.EXPECT().
					IntrospectToken(gomock.Any(), gomock.Eq("access-token"), gomock.Eq(fosite.AccessToken), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ fosite.TokenType, session fosite.Session, _ ...string) (fosite.TokenType, fosite.AccessRequester, error) {
						session = &oauth2.Session{
							DefaultSession: &openid.DefaultSession{
								Claims: &jwt.IDTokenClaims{
									Subject: "alice",
								},
								Headers: new(jwt.Headers),
								Subject: "alice",
							},
							Extra: map[string]interface{}{},
						}

						return fosite.AccessToken, &fosite.AccessRequest{
							Request: fosite.Request{
								Client: &client.Client{
									UserinfoSignedResponseAlg: "none",
								},
								Session: session,
							},
						}, nil
					})
			},
			expectStatusCode: http.StatusOK,
			checkForSuccess: func(t *testing.T, body []byte) {
				assert.True(t, strings.Contains(string(body), `"sub":"alice"`), "%s", body)
			},
		},
		{
			setup: func(t *testing.T) {
				op.EXPECT().
					IntrospectToken(gomock.Any(), gomock.Eq("access-token"), gomock.Eq(fosite.AccessToken), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ fosite.TokenType, session fosite.Session, _ ...string) (fosite.TokenType, fosite.AccessRequester, error) {
						session = &oauth2.Session{
							DefaultSession: &openid.DefaultSession{
								Claims: &jwt.IDTokenClaims{
									Subject: "alice",
								},
								Headers: new(jwt.Headers),
								Subject: "alice",
							},
							Extra: map[string]interface{}{},
						}

						return fosite.AccessToken, &fosite.AccessRequest{
							Request: fosite.Request{
								Client: &client.Client{
									UserinfoSignedResponseAlg: "asdfasdf",
								},
								Session: session,
							},
						}, nil
					})
			},
			expectStatusCode: http.StatusInternalServerError,
		},
		{
			setup: func(t *testing.T) {
				op.EXPECT().
					IntrospectToken(gomock.Any(), gomock.Eq("access-token"), gomock.Eq(fosite.AccessToken), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ fosite.TokenType, session fosite.Session, _ ...string) (fosite.TokenType, fosite.AccessRequester, error) {
						session = &oauth2.Session{
							DefaultSession: &openid.DefaultSession{
								Claims: &jwt.IDTokenClaims{
									Subject: "alice",
								},
								Headers: new(jwt.Headers),
								Subject: "alice",
							},
							Extra: map[string]interface{}{},
						}

						return fosite.AccessToken, &fosite.AccessRequest{
							Request: fosite.Request{
								Client: &client.Client{
									OutfacingID:               "foobar-client",
									UserinfoSignedResponseAlg: "RS256",
								},
								Session: session,
							},
						}, nil
					})
			},
			expectStatusCode: http.StatusOK,
			checkForSuccess: func(t *testing.T, body []byte) {
				claims, err := jwt2.Parse(string(body), func(token *jwt2.Token) (interface{}, error) {
					keys, err := reg.KeyManager().GetKeySet(context.Background(), x.OpenIDConnectKeyName)
					require.NoError(t, err)
					t.Logf("%+v", keys)
					key, err := jwk.FindKeyByPrefix(keys, "public")
					return jwk.MustRSAPublic(key), nil
				})
				require.NoError(t, err)
				assert.EqualValues(t, "alice", claims.Claims.(jwt2.MapClaims)["sub"])
				assert.EqualValues(t, []interface{}{"foobar-client"}, claims.Claims.(jwt2.MapClaims)["aud"], "%#v", claims.Claims)
				assert.NotEmpty(t, claims.Claims.(jwt2.MapClaims)["jti"])
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			tc.setup(t)

			req, err := http.NewRequest("GET", ts.URL+"/userinfo", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer access-token")
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.EqualValues(t, tc.expectStatusCode, resp.StatusCode)
			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			if tc.expectStatusCode == http.StatusOK {
				tc.checkForSuccess(t, body)
			} else if tc.expectStatusCode == http.StatusUnauthorized {
				tc.checkForUnauthorized(t, body, resp.Header)
			}
		})
	}
}

func TestHandlerWellKnown(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	viper.Set(configuration.ViperKeyScopeStrategy, "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY")
	viper.Set(configuration.ViperKeyIssuerURL, "http://hydra.localhost")
	viper.Set(configuration.ViperKeySubjectTypesSupported, []string{"pairwise", "public"})
	viper.Set(configuration.ViperKeyOIDCDiscoverySupportedClaims, []string{"sub"})
	viper.Set(configuration.ViperKeyOAuth2ClientRegistrationURL, "http://client-register/registration")
	viper.Set(configuration.ViperKeyOIDCDiscoveryUserinfoEndpoint, "/userinfo")
	reg := internal.NewRegistryMemory(t, conf)

	h := oauth2.NewHandler(reg, conf)

	r := x.NewRouterAdmin()
	h.SetRoutes(r, r.RouterPublic(), func(h http.Handler) http.Handler {
		return h
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/.well-known/openid-configuration")
	require.NoError(t, err)
	defer res.Body.Close()

	trueConfig := oauth2.WellKnown{
		Issuer:                                 strings.TrimRight(conf.IssuerURL().String(), "/") + "/",
		AuthURL:                                conf.OAuth2AuthURL().String(),
		TokenURL:                               conf.OAuth2TokenURL().String(),
		JWKsURI:                                conf.JWKSURL().String(),
		RevocationEndpoint:                     urlx.AppendPaths(conf.IssuerURL(), oauth2.RevocationPath).String(),
		RegistrationEndpoint:                   conf.OAuth2ClientRegistrationURL().String(),
		SubjectTypes:                           []string{"pairwise", "public"},
		ResponseTypes:                          []string{"code", "code id_token", "id_token", "token id_token", "token", "token id_token code"},
		ClaimsSupported:                        conf.OIDCDiscoverySupportedClaims(),
		ScopesSupported:                        conf.OIDCDiscoverySupportedScope(),
		UserinfoEndpoint:                       conf.OIDCDiscoveryUserinfoEndpoint(),
		TokenEndpointAuthMethodsSupported:      []string{"client_secret_post", "client_secret_basic", "private_key_jwt", "none"},
		GrantTypesSupported:                    []string{"authorization_code", "implicit", "client_credentials", "refresh_token"},
		ResponseModesSupported:                 []string{"query", "fragment"},
		IDTokenSigningAlgValuesSupported:       []string{"RS256"},
		UserinfoSigningAlgValuesSupported:      []string{"none", "RS256"},
		RequestParameterSupported:              true,
		RequestURIParameterSupported:           true,
		RequireRequestURIRegistration:          true,
		BackChannelLogoutSupported:             true,
		BackChannelLogoutSessionSupported:      true,
		FrontChannelLogoutSupported:            true,
		FrontChannelLogoutSessionSupported:     true,
		EndSessionEndpoint:                     urlx.AppendPaths(conf.IssuerURL(), oauth2.LogoutPath).String(),
		RequestObjectSigningAlgValuesSupported: []string{"RS256", "none"},
	}
	var wellKnownResp oauth2.WellKnown
	err = json.NewDecoder(res.Body).Decode(&wellKnownResp)
	require.NoError(t, err, "problem decoding wellknown json response: %+v", err)
	assert.EqualValues(t, trueConfig, wellKnownResp)
}
