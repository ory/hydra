// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package integration_test

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ory/hydra/v2/fosite/internal"
	"github.com/ory/hydra/v2/fosite/internal/gen"

	"github.com/go-jose/go-jose/v3"
	"github.com/gorilla/mux"
	goauth "golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/integration/clients"
	"github.com/ory/hydra/v2/fosite/storage"
	"github.com/ory/hydra/v2/fosite/token/hmac"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

const (
	firstKeyID  = "123"
	secondKeyID = "321"

	firstJWTBearerIssuer  = "first@example.com"
	secondJWTBearerIssuer = "second@example.com"

	firstJWTBearerSubject  = "first-service-client"
	secondJWTBearerSubject = "second-service-client"

	tokenURL          = "https://www.ory.sh/api"
	tokenRelativePath = "/token"

	deviceAuthRelativePath = "/device/auth"
)

var (
	firstPrivateKey, _  = rsa.GenerateKey(rand.Reader, 2048)
	secondPrivateKey, _ = rsa.GenerateKey(rand.Reader, 2048)
)

var fositeStore = &storage.MemoryStore{
	Clients: map[string]fosite.Client{
		"my-client": &fosite.DefaultClient{
			ID:            "my-client",
			Secret:        []byte(`$2a$10$IxMdI6d.LIRZPpSfEwNoeu4rY3FhDREsxFJXikcgdRRAStxUlsuEO`), // = "foobar"
			RedirectURIs:  []string{"http://localhost:3846/callback"},
			ResponseTypes: []string{"id_token", "code", "token", "token code", "id_token code", "token id_token", "token code id_token"},
			GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials", "urn:ietf:params:oauth:grant-type:device_code"},
			Scopes:        []string{"fosite", "offline", "openid"},
			Audience:      []string{tokenURL},
		},
		"custom-lifespan-client": &fosite.DefaultClientWithCustomTokenLifespans{
			DefaultClient: &fosite.DefaultClient{
				ID:             "custom-lifespan-client",
				Secret:         []byte(`$2a$10$IxMdI6d.LIRZPpSfEwNoeu4rY3FhDREsxFJXikcgdRRAStxUlsuEO`),            // = "foobar"
				RotatedSecrets: [][]byte{[]byte(`$2y$10$X51gLxUQJ.hGw1epgHTE5u0bt64xM0COU7K9iAp.OFg8p2pUd.1zC `)}, // = "foobaz",
				RedirectURIs:   []string{"http://localhost:3846/callback"},
				ResponseTypes:  []string{"id_token", "code", "token", "id_token token", "code id_token", "code token", "code id_token token"},
				GrantTypes:     []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
				Scopes:         []string{"fosite", "openid", "photos", "offline"},
			},
			TokenLifespans: &internal.TestLifespans,
		},
		"public-client": &fosite.DefaultClient{
			ID:            "public-client",
			Secret:        []byte{},
			Public:        true,
			RedirectURIs:  []string{"http://localhost:3846/callback"},
			ResponseTypes: []string{"id_token", "code", "code id_token"},
			GrantTypes:    []string{"refresh_token", "authorization_code"},
			Scopes:        []string{"fosite", "offline", "openid"},
			Audience:      []string{tokenURL},
		},
		"device-client": &fosite.DefaultClient{
			ID:         "device-client",
			Secret:     []byte(`$2a$10$IxMdI6d.LIRZPpSfEwNoeu4rY3FhDREsxFJXikcgdRRAStxUlsuEO`), // = "foobar"
			GrantTypes: []string{"urn:ietf:params:oauth:grant-type:device_code", "refresh_token"},
			Scopes:     []string{"fosite", "offline", "openid"},
			Audience:   []string{tokenURL},
			Public:     true,
		},
	},
	Users: map[string]storage.MemoryUserRelation{
		"peter": {
			Username: "peter",
			Password: "secret",
		},
	},
	IssuerPublicKeys: map[string]storage.IssuerPublicKeys{
		firstJWTBearerIssuer: createIssuerPublicKey(
			firstJWTBearerIssuer,
			firstJWTBearerSubject,
			firstKeyID,
			firstPrivateKey.Public(),
			[]string{"fosite", "gitlab", "example.com", "docker"},
		),
		secondJWTBearerIssuer: createIssuerPublicKey(
			secondJWTBearerIssuer,
			secondJWTBearerSubject,
			secondKeyID,
			secondPrivateKey.Public(),
			[]string{"fosite"},
		),
	},
	BlacklistedJTIs:        map[string]time.Time{},
	AuthorizeCodes:         map[string]storage.StoreAuthorizeCode{},
	PKCES:                  map[string]fosite.Requester{},
	AccessTokens:           map[string]fosite.Requester{},
	RefreshTokens:          map[string]storage.StoreRefreshToken{},
	IDSessions:             map[string]fosite.Requester{},
	AccessTokenRequestIDs:  map[string]string{},
	RefreshTokenRequestIDs: map[string]string{},
	PARSessions:            map[string]fosite.AuthorizeRequester{},
	DeviceAuths:            map[string]fosite.DeviceRequester{},
	DeviceCodesRequestIDs:  map[string]storage.DeviceAuthPair{},
	UserCodesRequestIDs:    map[string]string{},
}

type defaultSession struct {
	*openid.DefaultSession
}

var accessTokenLifespan = time.Hour

var authCodeLifespan = time.Minute

func createIssuerPublicKey(issuer, subject, keyID string, key crypto.PublicKey, scopes []string) storage.IssuerPublicKeys {
	return storage.IssuerPublicKeys{
		Issuer: issuer,
		KeysBySub: map[string]storage.SubjectPublicKeys{
			subject: {
				Subject: subject,
				Keys: map[string]storage.PublicKeyScopes{
					keyID: {
						Key: &jose.JSONWebKey{
							Key:       key,
							Algorithm: string(jose.RS256),
							Use:       "sig",
							KeyID:     keyID,
						},
						Scopes: scopes,
					},
				},
			},
		},
	}
}

func newOAuth2Client(ts *httptest.Server) *goauth.Config {
	return &goauth.Config{
		ClientID:     "my-client",
		ClientSecret: "foobar",
		RedirectURL:  ts.URL + "/callback",
		Scopes:       []string{"fosite"},
		Endpoint: goauth.Endpoint{
			TokenURL:  ts.URL + tokenRelativePath,
			AuthURL:   ts.URL + "/auth",
			AuthStyle: goauth.AuthStyleInHeader,
		},
	}
}

func newOAuth2AppClient(ts *httptest.Server) *clientcredentials.Config {
	return &clientcredentials.Config{
		ClientID:     "my-client",
		ClientSecret: "foobar",
		Scopes:       []string{"fosite"},
		TokenURL:     ts.URL + tokenRelativePath,
	}
}

func newJWTBearerAppClient(ts *httptest.Server) *clients.JWTBearer {
	return clients.NewJWTBearer(ts.URL + tokenRelativePath)
}

var hmacStrategy = oauth2.NewHMACSHAStrategy(
	&hmac.HMACStrategy{
		Config: &fosite.Config{
			GlobalSecret: []byte("some-super-cool-secret-that-nobody-knows"),
		},
	},
	&fosite.Config{
		AccessTokenLifespan:   accessTokenLifespan,
		AuthorizeCodeLifespan: authCodeLifespan,
	},
)

var (
	defaultRSAKey = gen.MustRSAKey()
	jwtStrategy   = &oauth2.DefaultJWTStrategy{
		Signer: &jwt.DefaultSigner{
			GetPrivateKey: func(ctx context.Context) (interface{}, error) {
				return defaultRSAKey, nil
			},
		},
		Config:   &fosite.Config{},
		Strategy: hmacStrategy,
	}
)

func mockServer(t *testing.T, f fosite.OAuth2Provider, session fosite.Session) *httptest.Server {
	router := mux.NewRouter()
	router.HandleFunc("/auth", authEndpointHandler(t, f, session))
	router.HandleFunc(tokenRelativePath, tokenEndpointHandler(t, f))
	router.HandleFunc("/callback", authCallbackHandler(t))
	router.HandleFunc("/info", tokenInfoHandler(t, f, session))
	router.HandleFunc("/introspect", tokenIntrospectionHandler(t, f, session))
	router.HandleFunc("/revoke", tokenRevocationHandler(t, f, session))
	router.HandleFunc("/par", pushedAuthorizeRequestHandler(t, f, session))
	router.HandleFunc(deviceAuthRelativePath, deviceAuthorizationEndpointHandler(t, f, session))

	ts := httptest.NewServer(router)
	return ts
}
