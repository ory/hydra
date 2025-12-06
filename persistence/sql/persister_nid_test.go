// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/oauth2/trust"
	"github.com/ory/hydra/v2/persistence"
	persistencesql "github.com/ory/hydra/v2/persistence/sql"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/assertx"
	"github.com/ory/x/contextx"
	"github.com/ory/x/networkx"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/servicelocatorx"
	"github.com/ory/x/sqlcon"
	"github.com/ory/x/sqlxx"
	"github.com/ory/x/uuidx"
)

type PersisterTestSuite struct {
	suite.Suite
	registries map[string]*driver.RegistrySQL
	t1         context.Context
	t2         context.Context
	t1NID      uuid.UUID
	t2NID      uuid.UUID
}

var _ interface {
	suite.SetupAllSuite
	suite.TearDownTestSuite
} = (*PersisterTestSuite)(nil)

func (s *PersisterTestSuite) SetupSuite() {
	withCtxer := driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.TestContextualizer{}))

	s.registries = testhelpers.ConnectDatabases(s.T(), true, withCtxer)

	s.t1NID, s.t2NID = uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4())
	s.t1 = contextx.SetNIDContext(context.Background(), s.t1NID)
	s.t2 = contextx.SetNIDContext(context.Background(), s.t2NID)

	for _, r := range s.registries {
		require.NoError(s.T(), r.Persister().Connection(context.Background()).Create(&networkx.Network{ID: s.t1NID}))
		require.NoError(s.T(), r.Persister().Connection(context.Background()).Create(&networkx.Network{ID: s.t2NID}))
	}
}

func (s *PersisterTestSuite) TearDownTest() {
	for _, r := range s.registries {
		x.DeleteHydraRows(s.T(), r.Persister().Connection(context.Background()))
	}
}

func (s *PersisterTestSuite) TestAcceptLogoutRequest() {
	lr := newLogoutRequest()

	for k, r := range s.registries {
		s.T().Run("dialect="+k, func(t *testing.T) {
			require.NoError(t, r.LogoutManager().CreateLogoutRequest(s.t1, lr))

			expected, err := r.LogoutManager().GetLogoutRequest(s.t1, lr.ID)
			require.NoError(t, err)
			require.Equal(t, false, expected.Accepted)

			lrAccepted, err := r.LogoutManager().AcceptLogoutRequest(s.t2, lr.ID)
			require.Error(t, err)
			require.Equal(t, &flow.LogoutRequest{}, lrAccepted)

			actual, err := r.LogoutManager().GetLogoutRequest(s.t1, lr.ID)
			require.NoError(t, err)
			require.Equal(t, expected, actual)
		})
	}
}

func (s *PersisterTestSuite) TestAddKeyGetKeyDeleteKey() {
	key := newKey("test-ks", "test")
	for k, r := range s.registries {
		s.T().Run("dialect="+k, func(t *testing.T) {
			ks := "key-set"
			require.NoError(t, r.KeyManager().AddKey(s.t1, ks, &key))
			actual, err := r.KeyManager().GetKey(s.t2, ks, key.KeyID)
			require.Error(t, err)
			require.Equal(t, (*jose.JSONWebKeySet)(nil), actual)
			actual, err = r.KeyManager().GetKey(s.t1, ks, key.KeyID)
			require.NoError(t, err)
			assertx.EqualAsJSON(t, &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{key}}, actual)

			require.NoError(t, r.KeyManager().DeleteKey(s.t2, ks, key.KeyID))
			_, err = r.KeyManager().GetKey(s.t1, ks, key.KeyID)
			require.NoError(t, err)
			require.NoError(t, r.KeyManager().DeleteKey(s.t1, ks, key.KeyID))
			_, err = r.KeyManager().GetKey(s.t1, ks, key.KeyID)
			require.Error(t, err)
		})
	}
}

func (s *PersisterTestSuite) TestAddKeySetGetKeySetDeleteKeySet() {
	ks := newKeySet("test-ks", "test")
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			ksID := "key-set"
			require.NoError(t, r.KeyManager().AddKeySet(s.t1, ksID, ks))
			actual, err := r.KeyManager().GetKeySet(s.t2, ksID)
			require.Error(t, err)
			require.Equal(t, (*jose.JSONWebKeySet)(nil), actual)
			actual, err = r.KeyManager().GetKeySet(s.t1, ksID)
			require.NoError(t, err)
			requireKeySetEqual(t, ks, actual)

			require.NoError(t, r.KeyManager().DeleteKeySet(s.t2, ksID))
			_, err = r.KeyManager().GetKeySet(s.t1, ksID)
			require.NoError(t, err)
			require.NoError(t, r.KeyManager().DeleteKeySet(s.t1, ksID))
			_, err = r.KeyManager().GetKeySet(s.t1, ksID)
			require.Error(t, err)
		})
	}
}

func (s *PersisterTestSuite) TestAuthenticate() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id", Secret: "secret"}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))

			actual, err := r.Persister().AuthenticateClient(s.t2, "client-id", []byte("secret"))
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().AuthenticateClient(s.t1, "client-id", []byte("secret"))
			require.NoError(t, err)
			require.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestClientAssertionJWTValid() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			jti := oauth2.NewBlacklistedJTI(uuid.Must(uuid.NewV4()).String(), time.Now().Add(24*time.Hour))
			require.NoError(t, r.Persister().SetClientAssertionJWT(s.t1, jti.JTI, jti.Expiry))

			require.NoError(t, r.Persister().ClientAssertionJWTValid(s.t2, jti.JTI))
			require.Error(t, r.Persister().ClientAssertionJWTValid(s.t1, jti.JTI))
		})
	}
}

func (s *PersisterTestSuite) TestConfirmLoginSession() {
	ls := newLoginSession()
	ls.AuthenticatedAt = sqlxx.NullTime(time.Now().UTC())
	ls.Remember = true
	ls.NID = s.t1NID

	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			// Expects the login session to be confirmed in the correct context.
			require.NoError(t, r.Persister().ConfirmLoginSession(s.t1, ls))
			actual := &flow.LoginSession{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(actual, ls.ID))

			require.True(t, time.Time(ls.AuthenticatedAt).UTC().Equal(time.Time(actual.AuthenticatedAt).UTC()))
			require.True(t, time.Time(ls.ExpiresAt).UTC().Equal(time.Time(actual.ExpiresAt).UTC()))

			exp, _ := json.Marshal(ls)
			act, _ := json.Marshal(actual)
			require.JSONEq(t, string(exp), string(act))

			// Can't find the login session in the wrong context.
			require.ErrorIs(t,
				r.Persister().ConfirmLoginSession(s.t2, ls),
				x.ErrNotFound,
			)
		})
	}
}

func (s *PersisterTestSuite) TestCreateSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			ls := newLoginSession()
			require.NoError(t, r.Persister().ConfirmLoginSession(s.t1, ls))
			actual := &flow.LoginSession{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(actual, ls.ID))
			require.Equal(t, s.t1NID, actual.NID)
			require.Equal(t, ls, actual)
		})
	}
}

func (s *PersisterTestSuite) TestCreateAccessTokenSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			c1 := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, c1))
			c2 := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t2, c2))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()

			fr.Client = &fosite.DefaultClient{ID: c1.ID}
			fr.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}
			require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1, sig, fr))
			actual := persistencesql.OAuth2RequestSQL{Table: "access"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, x.SignatureHash(sig)))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateAuthorizeCodeSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			c1 := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, c1))
			c2 := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t2, c2))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: c1.ID}
			fr.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}
			require.NoError(t, r.Persister().CreateAuthorizeCodeSession(s.t1, sig, fr))
			actual := persistencesql.OAuth2RequestSQL{Table: "code"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, sig))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateClient() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			expected := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, expected))
			actual := client.Client{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, expected.ID))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateForcedObfuscatedLoginSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			session := &consent.ForcedObfuscatedLoginSession{ClientID: cl.ID}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			require.NoError(t, r.Persister().CreateForcedObfuscatedLoginSession(s.t1, session))
			actual, err := r.Persister().GetForcedObfuscatedLoginSession(s.t1, cl.ID, "")
			require.NoError(t, err)
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateGrant() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()),
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().CreateGrant(s.t1, grant, ks.Keys[0].Public()))
			actual := persistencesql.SQLGrant{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, grant.ID))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateLogoutRequest() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			lr := flow.LogoutRequest{
				// TODO there is not FK for SessionID so we don't need it here; TODO make sure the missing FK is intentional
				ID:       uuid.Must(uuid.NewV4()).String(),
				ClientID: sql.NullString{Valid: true, String: cl.ID},
			}

			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			require.NoError(t, r.Persister().CreateLogoutRequest(s.t1, &lr))
			actual, err := r.Persister().GetLogoutRequest(s.t1, lr.ID)
			require.NoError(t, err)
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateOpenIDConnectSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			request.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}

			authorizeCode := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateOpenIDConnectSession(s.t1, authorizeCode, request))

			actual := persistencesql.OAuth2RequestSQL{Table: "oidc"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreatePKCERequestSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			request.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}

			authorizeCode := uuid.Must(uuid.NewV4()).String()

			actual := persistencesql.OAuth2RequestSQL{Table: "pkce"}
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.NoError(t, r.Persister().CreatePKCERequestSession(s.t1, authorizeCode, request))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateRefreshTokenSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			request.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}

			authorizeCode := uuid.Must(uuid.NewV4()).String()
			actual := persistencesql.OAuth2RefreshTable{}
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t1, authorizeCode, "", request))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateWithNetwork() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			expected := &client.Client{ID: "client-id"}
			store, ok := r.OAuth2Storage().(*persistencesql.Persister)
			require.True(t, ok)
			require.NoError(t, store.CreateWithNetwork(s.t1, expected))

			actual := &client.Client{}
			require.NoError(t, r.Persister().Connection(context.Background()).Where("id = ?", expected.ID).First(actual))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) DeleteAccessTokenSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: cl.ID}
			fr.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}
			require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1, sig, fr))
			require.NoError(t, r.Persister().DeleteAccessTokenSession(s.t2, sig))

			actual := persistencesql.OAuth2RequestSQL{Table: "access"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, x.SignatureHash(sig)))
			require.Equal(t, s.t1NID, actual.NID)

			require.NoError(t, r.Persister().DeleteAccessTokenSession(s.t1, sig))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, x.SignatureHash(sig)))
		})
	}
}

func (s *PersisterTestSuite) TestDeleteAccessTokens() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: cl.ID}
			fr.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}
			require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1, sig, fr))
			require.NoError(t, r.Persister().DeleteAccessTokens(s.t2, cl.ID))

			actual := persistencesql.OAuth2RequestSQL{Table: "access"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, x.SignatureHash(sig)))
			require.Equal(t, s.t1NID, actual.NID)

			require.NoError(t, r.Persister().DeleteAccessTokens(s.t1, cl.ID))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, x.SignatureHash(sig)))
		})
	}
}

func (s *PersisterTestSuite) TestDeleteClient() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			c := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, c))
			actual := client.Client{}
			require.Error(t, r.Persister().DeleteClient(s.t2, c.ID))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, c.ID))
			require.NoError(t, r.Persister().DeleteClient(s.t1, c.ID))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, c.ID))
		})
	}
}

func (s *PersisterTestSuite) TestDeleteGrant() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()),
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().CreateGrant(s.t1, grant, ks.Keys[0].Public()))

			actual := persistencesql.SQLGrant{}
			require.Error(t, r.Persister().DeleteGrant(s.t2, grant.ID))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, grant.ID))
			require.NoError(t, r.Persister().DeleteGrant(s.t1, grant.ID))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, grant.ID))
		})
	}
}

func (s *PersisterTestSuite) TestDeleteLoginSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			ls := flow.LoginSession{
				ID:                        uuid.Must(uuid.NewV4()).String(),
				Remember:                  true,
				IdentityProviderSessionID: sqlxx.NullString(uuid.Must(uuid.NewV4()).String()),
			}
			require.NoError(t, r.Persister().ConfirmLoginSession(s.t1, &ls))

			deletedLS, err := r.Persister().DeleteLoginSession(s.t2, ls.ID)
			require.ErrorIs(t, err, sqlcon.ErrNoRows)
			assert.Nil(t, deletedLS)
			_, err = r.Persister().GetRememberedLoginSession(s.t1, ls.ID)
			require.NoError(t, err)

			deletedLS, err = r.Persister().DeleteLoginSession(s.t1, ls.ID)
			require.NoError(t, err)
			assert.Equal(t, ls, *deletedLS)
			_, err = r.Persister().GetRememberedLoginSession(s.t1, ls.ID)
			require.ErrorIs(t, err, x.ErrNotFound)
		})
	}
}

func (s *PersisterTestSuite) TestDeleteOpenIDConnectSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			request.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}

			authorizeCode := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateOpenIDConnectSession(s.t1, authorizeCode, request))

			actual := persistencesql.OAuth2RequestSQL{Table: "oidc"}

			require.NoError(t, r.Persister().DeleteOpenIDConnectSession(s.t2, authorizeCode))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.NoError(t, r.Persister().DeleteOpenIDConnectSession(s.t1, authorizeCode))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
		})
	}
}

func (s *PersisterTestSuite) TestDeletePKCERequestSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			request.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}

			authorizeCode := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreatePKCERequestSession(s.t1, authorizeCode, request))

			actual := persistencesql.OAuth2RequestSQL{Table: "pkce"}

			require.NoError(t, r.Persister().DeletePKCERequestSession(s.t2, authorizeCode))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.NoError(t, r.Persister().DeletePKCERequestSession(s.t1, authorizeCode))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
		})
	}
}

func (s *PersisterTestSuite) TestDeleteRefreshTokenSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			request.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}

			signature := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t1, signature, "", request))

			actual := persistencesql.OAuth2RefreshTable{}

			require.NoError(t, r.Persister().DeleteRefreshTokenSession(s.t2, signature))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, signature))
			require.NoError(t, r.Persister().DeleteRefreshTokenSession(s.t1, signature))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, signature))
		})
	}
}

func (s *PersisterTestSuite) TestDetermineNetwork() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			store, ok := r.OAuth2Storage().(*persistencesql.Persister)
			require.True(t, ok)

			require.NoError(t, r.Persister().Connection(t.Context()).Where("id <> ? AND id <> ?", s.t1NID, s.t2NID).Delete(&networkx.Network{}))

			actual, err := store.DetermineNetwork(t.Context())
			require.NoError(t, err)
			assert.Contains(t, []uuid.UUID{s.t1NID, s.t2NID}, actual.ID)
		})
	}
}

func (s *PersisterTestSuite) TestFindGrantedAndRememberedConsentRequests() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			cl := &client.Client{ID: "client-id"}
			f := newFlow(s.t1NID, cl.ID, "sub", sqlxx.NullString(sessionID))
			persistLoginSession(s.t1, t, r.Persister(), &flow.LoginSession{ID: sessionID})
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))

			req := &flow.OAuth2ConsentRequest{
				ConsentRequestID: "consent-request-id",
				LoginChallenge:   sqlxx.NullString(f.ID),
				Skip:             false,
			}

			f.ConsentRequestID = sqlxx.NullString(req.ConsentRequestID)
			require.NoError(t, f.HandleConsentRequest(&flow.AcceptOAuth2ConsentRequest{
				Remember: true,
			}))

			f.State = flow.FlowStateConsentUsed
			require.NoError(t, r.Persister().CreateConsentSession(s.t1, f))

			actual, err := r.Persister().FindGrantedAndRememberedConsentRequest(s.t2, cl.ID, f.Subject)
			require.ErrorIs(t, err, consent.ErrNoPreviousConsentFound)
			assert.Nil(t, actual)

			actual, err = r.Persister().FindGrantedAndRememberedConsentRequest(s.t1, cl.ID, f.Subject)
			require.NoError(t, err)
			assert.EqualValues(t, req.ConsentRequestID, actual.ConsentRequestID)
		})
	}
}

func (s *PersisterTestSuite) TestFindSubjectsGrantedConsentRequests() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			cl := &client.Client{ID: "client-id"}
			f := newFlow(s.t1NID, cl.ID, "sub", sqlxx.NullString(sessionID))
			persistLoginSession(s.t1, t, r.Persister(), &flow.LoginSession{ID: sessionID})
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			require.NoError(t, r.Persister().CreateConsentSession(s.t1, f))

			_, _, err := r.Persister().FindSubjectsGrantedConsentRequests(s.t2, f.Subject)
			require.ErrorIs(t, err, consent.ErrNoPreviousConsentFound)

			actual, nextPage, err := r.Persister().FindSubjectsGrantedConsentRequests(s.t1, f.Subject)
			require.NoError(t, err)
			require.Len(t, actual, 1)
			assert.Equal(t, f.ConsentRequestID.String(), actual[0].ConsentRequestID.String())
			assert.True(t, nextPage.IsLast())
		})
	}
}

func (s *PersisterTestSuite) TestFlushInactiveAccessTokens() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.RequestedAt = time.Now().UTC().Add(-24 * time.Hour)
			fr.Client = &fosite.DefaultClient{ID: cl.ID}
			fr.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}
			require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1, sig, fr))

			actual := persistencesql.OAuth2RequestSQL{Table: "access"}

			require.NoError(t, r.Persister().FlushInactiveAccessTokens(s.t2, time.Now().Add(time.Hour), 100, 100))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, x.SignatureHash(sig)))
			require.NoError(t, r.Persister().FlushInactiveAccessTokens(s.t1, time.Now().Add(time.Hour), 100, 100))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, x.SignatureHash(sig)))
		})
	}
}

func (s *PersisterTestSuite) TestGenerateAndPersistKeySet() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			actual := &jwk.SQLData{}

			key, err := jwk.GenerateJWK("RS256", "kid", "use")
			require.NoError(t, err)
			require.NoError(t, r.KeyManager().AddKey(s.t1, "ks", pointerx.Ptr(key.Keys[0].Public())))

			err = sqlcon.HandleError(r.Persister().Connection(t.Context()).Where("sid = ? AND kid = ? AND nid = ?", "ks", "kid", s.t2NID).First(actual))
			require.ErrorIs(t, err, sqlcon.ErrNoRows)
			require.NoError(t, r.Persister().Connection(t.Context()).Where("sid = ? AND kid = ? AND nid = ?", "ks", "kid", s.t1NID).First(actual))
		})
	}
}

func (s *PersisterTestSuite) TestFlushInactiveGrants() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()),
				ExpiresAt: time.Now().Add(-24 * time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().CreateGrant(s.t1, grant, ks.Keys[0].Public()))

			actual := persistencesql.SQLGrant{}
			require.NoError(t, r.Persister().FlushInactiveGrants(s.t2, time.Now(), 100, 100))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, grant.ID))
			require.NoError(t, r.Persister().FlushInactiveGrants(s.t1, time.Now(), 100, 100))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, grant.ID))
		})
	}
}

func (s *PersisterTestSuite) TestFlushInactiveLoginConsentRequests() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			sessionID := uuidx.NewV4().String()
			cl := &client.Client{ID: uuidx.NewV4().String()}
			f := newFlow(s.t1NID, cl.ID, "sub", sqlxx.NullString(sessionID))
			f.RequestedAt = time.Now().Add(-24 * time.Hour)
			persistLoginSession(s.t1, t, r.Persister(), &flow.LoginSession{ID: sessionID})
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			require.NoError(t, r.Persister().Connection(s.t1).Create(&persistencesql.FlowWithConstantColumns{Flow: f, State: f.State}))

			actual := flow.Flow{}

			require.NoError(t, r.Persister().FlushInactiveLoginConsentRequests(s.t2, time.Now(), 100, 100))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, f.ID))
			require.NoError(t, r.Persister().FlushInactiveLoginConsentRequests(s.t1, time.Now(), 100, 100))
			require.ErrorIs(t, r.Persister().Connection(context.Background()).Find(&actual, f.ID), sql.ErrNoRows)
		})
	}
}

func (s *PersisterTestSuite) TestFlushInactiveRefreshTokens() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			request := fosite.NewRequest()
			request.RequestedAt = time.Now().Add(-240 * 365 * time.Hour)
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			request.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}
			signature := uuid.Must(uuid.NewV4()).String()

			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t1, signature, "", request))

			actual := persistencesql.OAuth2RefreshTable{}

			require.NoError(t, r.Persister().FlushInactiveRefreshTokens(s.t2, time.Now(), 100, 100))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, signature))
			require.NoError(t, r.Persister().FlushInactiveRefreshTokens(s.t1, time.Now(), 100, 100))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, signature))
		})
	}
}

func (s *PersisterTestSuite) TestGetAccessTokenSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: cl.ID}
			fr.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}
			require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1, sig, fr))

			actual, err := r.Persister().GetAccessTokenSession(s.t2, sig, &fosite.DefaultSession{})
			require.Error(t, err)
			require.Nil(t, actual)
			actual, err = r.Persister().GetAccessTokenSession(s.t1, sig, &fosite.DefaultSession{})
			require.NoError(t, err)
			require.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetAuthorizeCodeSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: cl.ID}
			fr.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}
			require.NoError(t, r.Persister().CreateAuthorizeCodeSession(s.t1, sig, fr))

			actual, err := r.Persister().GetAuthorizeCodeSession(s.t2, sig, &fosite.DefaultSession{})
			require.Error(t, err)
			require.Nil(t, actual)
			actual, err = r.Persister().GetAuthorizeCodeSession(s.t1, sig, &fosite.DefaultSession{})
			require.NoError(t, err)
			require.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetClient() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			expected := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, expected))

			actual, err := r.Persister().GetClient(s.t2, expected.ID)
			require.Error(t, err)
			require.Nil(t, actual)
			actual, err = r.Persister().GetClient(s.t1, expected.ID)
			require.NoError(t, err)
			require.Equal(t, expected.ID, actual.GetID())
		})
	}
}

func (s *PersisterTestSuite) TestGetClientAssertionJWT() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			store, ok := r.OAuth2Storage().(oauth2.AssertionJWTReader)
			require.True(t, ok)
			expected := oauth2.NewBlacklistedJTI(uuid.Must(uuid.NewV4()).String(), time.Now().Add(24*time.Hour))
			require.NoError(t, r.Persister().SetClientAssertionJWT(s.t1, expected.JTI, expected.Expiry))

			_, err := store.GetClientAssertionJWT(s.t2, expected.JTI)
			require.Error(t, err)
			_, err = store.GetClientAssertionJWT(s.t1, expected.JTI)
			require.NoError(t, err)
		})
	}
}

func (s *PersisterTestSuite) TestGetClients() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			c := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, c))

			actual, nextPage, err := r.Persister().GetClients(s.t2, client.Filter{})
			require.NoError(t, err)
			assert.Len(t, actual, 0)
			assert.True(t, nextPage.IsLast())

			actual, nextPage, err = r.Persister().GetClients(s.t1, client.Filter{})
			require.NoError(t, err)
			assert.Len(t, actual, 1)
			assert.True(t, nextPage.IsLast())
		})
	}
}

func (s *PersisterTestSuite) TestGetConcreteClient() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			expected := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, expected))

			actual, err := r.Persister().GetConcreteClient(s.t2, expected.ID)
			require.Error(t, err)
			require.Nil(t, actual)
			actual, err = r.Persister().GetConcreteClient(s.t1, expected.ID)
			require.NoError(t, err)
			require.Equal(t, expected.ID, actual.GetID())
		})
	}
}

func (s *PersisterTestSuite) TestGetConcreteGrant() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()),
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().CreateGrant(s.t1, grant, ks.Keys[0].Public()))

			actual, err := r.Persister().GetConcreteGrant(s.t2, grant.ID)
			require.Error(t, err)
			require.Equal(t, trust.Grant{}, actual)

			actual, err = r.Persister().GetConcreteGrant(s.t1, grant.ID)
			require.NoError(t, err)
			require.NotEqual(t, trust.Grant{}, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetForcedObfuscatedLoginSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			session := &consent.ForcedObfuscatedLoginSession{ClientID: cl.ID}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			require.NoError(t, r.Persister().CreateForcedObfuscatedLoginSession(s.t1, session))

			actual, err := r.Persister().GetForcedObfuscatedLoginSession(s.t2, cl.ID, "")
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetForcedObfuscatedLoginSession(s.t1, cl.ID, "")
			require.NoError(t, err)
			require.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetGrants() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()),
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().CreateGrant(s.t1, grant, ks.Keys[0].Public()))

			actual, nextPage, err := r.Persister().GetGrants(s.t2, "")
			require.NoError(t, err)
			assert.Len(t, actual, 0)
			assert.True(t, nextPage.IsLast())

			actual, nextPage, err = r.Persister().GetGrants(s.t1, "")
			require.NoError(t, err)
			assert.Len(t, actual, 1)
			assert.True(t, nextPage.IsLast())
		})
	}
}

func (s *PersisterTestSuite) TestGetLogoutRequest() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			lr := flow.LogoutRequest{
				ID:       uuid.Must(uuid.NewV4()).String(),
				ClientID: sql.NullString{Valid: true, String: cl.ID},
			}

			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			require.NoError(t, r.Persister().CreateLogoutRequest(s.t1, &lr))

			actual, err := r.Persister().GetLogoutRequest(s.t2, lr.ID)
			require.Error(t, err)
			require.Equal(t, &flow.LogoutRequest{}, actual)

			actual, err = r.Persister().GetLogoutRequest(s.t1, lr.ID)
			require.NoError(t, err)
			require.NotEqual(t, &flow.LogoutRequest{}, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetOpenIDConnectSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			request := fosite.NewRequest()
			request.SetID("request-id")
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			request.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}
			authorizeCode := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			require.NoError(t, r.Persister().CreateOpenIDConnectSession(s.t1, authorizeCode, request))

			actual, err := r.Persister().GetOpenIDConnectSession(s.t2, authorizeCode, &fosite.Request{})
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetOpenIDConnectSession(s.t1, authorizeCode, &fosite.Request{})
			require.NoError(t, err)
			require.Equal(t, request.GetID(), actual.GetID())
		})
	}
}

func (s *PersisterTestSuite) TestGetPKCERequestSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			request := fosite.NewRequest()
			request.SetID("request-id")
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			request.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}
			sig := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			require.NoError(t, r.Persister().CreatePKCERequestSession(s.t1, sig, request))

			actual, err := r.Persister().GetPKCERequestSession(s.t2, sig, &fosite.DefaultSession{})
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetPKCERequestSession(s.t1, sig, &fosite.DefaultSession{})
			require.NoError(t, err)
			require.Equal(t, request.GetID(), actual.GetID())
		})
	}
}

func (s *PersisterTestSuite) TestGetPublicKey() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()),
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().CreateGrant(s.t1, grant, ks.Keys[0].Public()))

			actual, err := r.Persister().GetPublicKey(s.t2, grant.Issuer, grant.Subject, grant.PublicKey.KeyID)
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetPublicKey(s.t1, grant.Issuer, grant.Subject, grant.PublicKey.KeyID)
			require.NoError(t, err)
			require.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetPublicKeyScopes() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()),
				Scope:     []string{"a", "b", "c"},
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().CreateGrant(s.t1, grant, ks.Keys[0].Public()))

			actual, err := r.Persister().GetPublicKeyScopes(s.t2, grant.Issuer, grant.Subject, grant.PublicKey.KeyID)
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetPublicKeyScopes(s.t1, grant.Issuer, grant.Subject, grant.PublicKey.KeyID)
			require.NoError(t, err)
			require.Equal(t, grant.Scope, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetPublicKeys() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			const issuer = "ks-id"
			ks := newKeySet(issuer, "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()),
				ExpiresAt: time.Now().UTC().Add(time.Hour),
				Issuer:    issuer,
				PublicKey: trust.PublicKey{Set: issuer, KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().CreateGrant(s.t1, grant, ks.Keys[0].Public()))

			actual, err := r.Persister().GetPublicKeys(s.t2, grant.Issuer, grant.Subject)
			require.NoError(t, err)
			require.Nil(t, actual.Keys)

			actual, err = r.Persister().GetPublicKeys(s.t1, grant.Issuer, grant.Subject)
			require.NoError(t, err)
			require.NotNil(t, actual.Keys)
		})
	}
}

func (s *PersisterTestSuite) TestGetRefreshTokenSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			request := fosite.NewRequest()
			request.SetID("request-id")
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			request.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}
			sig := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t1, sig, "", request))

			actual, err := r.Persister().GetRefreshTokenSession(s.t2, sig, &fosite.DefaultSession{})
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetRefreshTokenSession(s.t1, sig, &fosite.DefaultSession{})
			require.NoError(t, err)
			require.Equal(t, request.GetID(), actual.GetID())
		})
	}
}

func (s *PersisterTestSuite) TestGetRememberedLoginSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			ls := flow.LoginSession{
				ID:       uuid.Must(uuid.NewV4()).String(),
				NID:      s.t1NID,
				Remember: true,
			}
			require.NoError(t, r.Persister().ConfirmLoginSession(s.t1, &ls))

			actual, err := r.Persister().GetRememberedLoginSession(s.t2, ls.ID)
			require.ErrorIs(t, err, x.ErrNotFound)
			require.Nil(t, actual)

			actual, err = r.Persister().GetRememberedLoginSession(s.t1, ls.ID)
			require.NoError(t, err)
			require.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestHandleConsentRequest() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			c1 := &client.Client{ID: uuidx.NewV4().String()}
			f := newFlow(s.t1NID, c1.ID, "sub", sqlxx.NullString(sessionID))

			persistLoginSession(s.t1, t, r.Persister(), &flow.LoginSession{ID: sessionID})
			require.NoError(t, r.Persister().CreateClient(s.t1, c1))
			require.NoError(t, r.Persister().CreateClient(s.t2, c1))

			f.ConsentRequestID = "consent-request-id"

			actual, err := r.Persister().FindGrantedAndRememberedConsentRequest(s.t1, c1.ID, f.Subject)
			require.ErrorIs(t, err, consent.ErrNoPreviousConsentFound)
			assert.Nil(t, actual)

			require.NoError(t, f.HandleConsentRequest(&flow.AcceptOAuth2ConsentRequest{
				Remember: true,
			}))

			f.State = flow.FlowStateConsentUsed

			require.NoError(t, r.Persister().CreateConsentSession(s.t1, f))
			actual, err = r.Persister().FindGrantedAndRememberedConsentRequest(s.t1, c1.ID, f.Subject)
			require.NoError(t, err)
			assert.EqualValues(t, f.ConsentRequestID, actual.ConsentRequestID)
		})
	}
}

func (s *PersisterTestSuite) TestInvalidateAuthorizeCodeSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: uuidx.NewV4().String()}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			require.NoError(t, r.Persister().CreateClient(s.t2, cl))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: cl.ID}
			fr.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}
			require.NoError(t, r.Persister().CreateAuthorizeCodeSession(s.t1, sig, fr))

			require.NoError(t, r.Persister().InvalidateAuthorizeCodeSession(s.t2, sig))
			actual := persistencesql.OAuth2RequestSQL{Table: "code"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, sig))
			require.Equal(t, true, actual.Active)

			require.NoError(t, r.Persister().InvalidateAuthorizeCodeSession(s.t1, sig))
			actual = persistencesql.OAuth2RequestSQL{Table: "code"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, sig))
			require.Equal(t, false, actual.Active)
		})
	}
}

func (s *PersisterTestSuite) TestIsJWTUsed() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			jti := oauth2.NewBlacklistedJTI(uuid.Must(uuid.NewV4()).String(), time.Now().Add(24*time.Hour))
			require.NoError(t, r.Persister().SetClientAssertionJWT(s.t1, jti.JTI, jti.Expiry))

			actual, err := r.Persister().IsJWTUsed(s.t2, jti.JTI)
			require.NoError(t, err)
			require.False(t, actual)

			actual, err = r.Persister().IsJWTUsed(s.t1, jti.JTI)
			require.NoError(t, err)
			require.True(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestListUserAuthenticatedClientsWithBackChannelLogout() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			c1 := &client.Client{ID: "client-1", BackChannelLogoutURI: "not-null"}
			c2 := &client.Client{ID: "client-2", BackChannelLogoutURI: "not-null"}
			require.NoError(t, r.Persister().CreateClient(s.t1, c1))
			require.NoError(t, r.Persister().CreateClient(s.t2, c1))
			require.NoError(t, r.Persister().CreateClient(s.t2, c2))

			t1f1 := newFlow(s.t1NID, c1.ID, "sub", sqlxx.NullString(uuid.Must(uuid.NewV4()).String()))
			t1f1.ConsentRequestID = "t1f1-consent-challenge"

			t2f1 := newFlow(s.t2NID, c1.ID, "sub", t1f1.SessionID)
			t2f1.ConsentRequestID = "t2f1-consent-challenge"

			t2f2 := newFlow(s.t2NID, c2.ID, "sub", t1f1.SessionID)
			t2f2.ConsentRequestID = "t2f2-consent-challenge"

			persistLoginSession(s.t1, t, r.Persister(), &flow.LoginSession{ID: t1f1.SessionID.String()})

			require.NoError(t, r.Persister().CreateConsentSession(s.t1, t1f1))
			require.NoError(t, r.Persister().CreateConsentSession(s.t2, t2f1))
			require.NoError(t, r.Persister().CreateConsentSession(s.t2, t2f2))

			t1f1.ConsentRequestID = sqlxx.NullString(t1f1.ID)
			t2f1.ConsentRequestID = sqlxx.NullString(t2f1.ID)
			t2f2.ConsentRequestID = sqlxx.NullString(t2f2.ID)

			cs, err := r.Persister().ListUserAuthenticatedClientsWithBackChannelLogout(s.t1, "sub", t1f1.SessionID.String())
			require.NoError(t, err)
			require.Equal(t, 1, len(cs))

			cs, err = r.Persister().ListUserAuthenticatedClientsWithBackChannelLogout(s.t2, "sub", t1f1.SessionID.String())
			require.NoError(t, err)
			require.Equal(t, 2, len(cs))
		})
	}
}

func (s *PersisterTestSuite) TestListUserAuthenticatedClientsWithFrontChannelLogout() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			c1 := &client.Client{ID: "client-1", FrontChannelLogoutURI: "not-null"}
			c2 := &client.Client{ID: "client-2", FrontChannelLogoutURI: "not-null"}
			require.NoError(t, r.Persister().CreateClient(s.t1, c1))
			require.NoError(t, r.Persister().CreateClient(s.t2, c1))
			require.NoError(t, r.Persister().CreateClient(s.t2, c2))

			t1f1 := newFlow(s.t1NID, c1.ID, "sub", sqlxx.NullString(uuid.Must(uuid.NewV4()).String()))
			t1f1.ConsentRequestID = "t1f1-consent-challenge"

			t2f1 := newFlow(s.t2NID, c1.ID, "sub", t1f1.SessionID)
			t2f1.ConsentRequestID = "t2f1-consent-challenge"

			t2f2 := newFlow(s.t2NID, c2.ID, "sub", t1f1.SessionID)
			t2f2.ConsentRequestID = "t2f2-consent-challenge"

			persistLoginSession(s.t1, t, r.Persister(), &flow.LoginSession{ID: t1f1.SessionID.String()})

			require.NoError(t, r.Persister().CreateConsentSession(s.t1, t1f1))
			require.NoError(t, r.Persister().CreateConsentSession(s.t2, t2f1))
			require.NoError(t, r.Persister().CreateConsentSession(s.t2, t2f2))

			t1f1.ConsentRequestID = sqlxx.NullString(t1f1.ID)
			t2f1.ConsentRequestID = sqlxx.NullString(t2f1.ID)
			t2f2.ConsentRequestID = sqlxx.NullString(t2f2.ID)

			cs, err := r.Persister().ListUserAuthenticatedClientsWithFrontChannelLogout(s.t1, "sub", t1f1.SessionID.String())
			require.NoError(t, err)
			require.Equal(t, 1, len(cs))

			cs, err = r.Persister().ListUserAuthenticatedClientsWithFrontChannelLogout(s.t2, "sub", t1f1.SessionID.String())
			require.NoError(t, err)
			require.Equal(t, 2, len(cs))
		})
	}
}

func (s *PersisterTestSuite) TestMarkJWTUsedForTime() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			require.NoError(t, r.Persister().SetClientAssertionJWT(s.t1, "a", time.Now().Add(-24*time.Hour)))
			require.NoError(t, r.Persister().SetClientAssertionJWT(s.t2, "a", time.Now().Add(-24*time.Hour)))
			require.NoError(t, r.Persister().SetClientAssertionJWT(s.t2, "b", time.Now().Add(-24*time.Hour)))

			require.NoError(t, r.Persister().MarkJWTUsedForTime(s.t2, "a", time.Now().Add(48*time.Hour)))

			store, ok := r.OAuth2Storage().(oauth2.AssertionJWTReader)
			require.True(t, ok)

			_, err := store.GetClientAssertionJWT(s.t1, "a")
			require.NoError(t, err)
			_, err = store.GetClientAssertionJWT(s.t2, "a")
			require.NoError(t, err)
			_, err = store.GetClientAssertionJWT(s.t2, "b")
			require.Error(t, err)
		})
	}
}

func (s *PersisterTestSuite) TestQueryWithNetwork() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			require.NoError(t, r.Persister().CreateClient(s.t1, &client.Client{ID: "client-1", FrontChannelLogoutURI: "not-null"}))

			store, ok := r.Persister().(*persistencesql.Persister)
			require.True(t, ok)

			var actual []client.Client
			require.NoError(t, store.QueryWithNetwork(s.t2).All(&actual))
			require.Len(t, actual, 0)
			require.NoError(t, store.QueryWithNetwork(s.t1).All(&actual))
			require.Len(t, actual, 1)
		})
	}
}

func (s *PersisterTestSuite) TestRejectLogoutRequest() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			lr := newLogoutRequest()
			require.NoError(t, r.LogoutManager().CreateLogoutRequest(s.t1, lr))

			require.Error(t, r.LogoutManager().RejectLogoutRequest(s.t2, lr.ID))
			actual, err := r.LogoutManager().GetLogoutRequest(s.t1, lr.ID)
			require.NoError(t, err)
			require.Equal(t, lr, actual)

			require.NoError(t, r.LogoutManager().RejectLogoutRequest(s.t1, lr.ID))
			actual, err = r.LogoutManager().GetLogoutRequest(s.t1, lr.ID)
			require.Error(t, err)
			require.Equal(t, &flow.LogoutRequest{}, actual)
		})
	}
}

func (s *PersisterTestSuite) TestRevokeAccessToken() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: cl.ID}
			fr.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}
			require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1, sig, fr))
			require.NoError(t, r.Persister().RevokeAccessToken(s.t2, fr.ID))

			actual := persistencesql.OAuth2RequestSQL{Table: "access"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, x.SignatureHash(sig)))
			require.Equal(t, s.t1NID, actual.NID)

			require.NoError(t, r.Persister().RevokeAccessToken(s.t1, fr.ID))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, x.SignatureHash(sig)))
		})
	}
}

func (s *PersisterTestSuite) TestRevokeRefreshToken() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			cl := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			request.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}

			signature := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t1, signature, "", request))

			var actualt2 persistencesql.OAuth2RefreshTable
			require.NoError(t, r.Persister().RevokeRefreshToken(s.t2, request.ID))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actualt2, signature))
			require.Equal(t, true, actualt2.Active)

			require.NoError(t, r.Persister().RevokeRefreshToken(s.t1, request.ID))
			require.ErrorIs(t, r.Persister().Connection(context.Background()).Find(new(persistencesql.OAuth2RefreshTable), signature), sql.ErrNoRows)
		})
	}
}

func (s *PersisterTestSuite) TestRotateRefreshToken() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			t.Run("with access signature", func(t *testing.T) {
				clientID := uuid.Must(uuid.NewV4()).String()
				require.NoError(t, r.Persister().CreateClient(s.t1, &client.Client{ID: clientID}))
				require.NoError(t, r.Persister().CreateClient(s.t2, &client.Client{ID: clientID}))

				request := fosite.NewRequest()
				request.Client = &fosite.DefaultClient{ID: clientID}
				request.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}

				// Create token T1
				signatureT1 := uuid.Must(uuid.NewV4()).String()
				accessSignatureT1 := uuid.Must(uuid.NewV4()).String()
				require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1, accessSignatureT1, request))
				require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t1, signatureT1, accessSignatureT1, request))

				// Create token T2
				signatureT2 := uuid.Must(uuid.NewV4()).String()
				accessSignatureT2 := uuid.Must(uuid.NewV4()).String()
				require.ErrorIs(t, r.Persister().RotateRefreshToken(s.t2, request.ID, signatureT2), fosite.ErrNotFound, "Rotation fails as token is non-existent.")
				require.NoError(t, r.Persister().CreateAccessTokenSession(s.t2, accessSignatureT2, request))
				require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t2, signatureT2, accessSignatureT2, request))

				accessT2 := persistencesql.OAuth2RequestSQL{Table: "access"}
				assert.NoError(t, r.Persister().Connection(s.t2).Where("signature = ?", x.SignatureHash(accessSignatureT2)).First(&accessT2))
				require.Equal(t, true, accessT2.Active)

				accessT1 := persistencesql.OAuth2RequestSQL{Table: "access"}
				assert.NoError(t, r.Persister().Connection(s.t1).Where("signature = ?", x.SignatureHash(accessSignatureT1)).First(&accessT1))
				require.Equal(t, true, accessT2.Active)

				// Rotate token T1
				require.NoError(t, r.Persister().RotateRefreshToken(s.t1, request.ID, signatureT1))
				{
					refreshT1 := persistencesql.OAuth2RefreshTable{}
					require.NoError(t, r.Persister().Connection(s.t1).Where("signature = ?", signatureT1).First(&refreshT1))
					require.Equal(t, false, refreshT1.Active)

					accessT1 := persistencesql.OAuth2RequestSQL{Table: "access"}
					require.ErrorIs(t, r.Persister().Connection(s.t1).Where("signature = ?", x.SignatureHash(accessSignatureT1)).First(&accessT1), sql.ErrNoRows)

					refreshT2 := persistencesql.OAuth2RefreshTable{}
					require.NoError(t, r.Persister().Connection(s.t2).Where("signature = ?", signatureT2).First(&refreshT2))
					require.Equal(t, true, refreshT2.Active)

					accessT2 := persistencesql.OAuth2RequestSQL{Table: "access"}
					require.NoError(t, r.Persister().Connection(s.t2).Where("signature = ?", x.SignatureHash(accessSignatureT2)).First(&accessT2))
					require.Equal(t, true, accessT2.Active)
				}

				require.NoError(t, r.Persister().RotateRefreshToken(s.t2, request.ID, signatureT2))
				{
					refreshT2 := persistencesql.OAuth2RefreshTable{}
					require.NoError(t, r.Persister().Connection(s.t2).Where("signature = ?", signatureT2).First(&refreshT2))
					require.Equal(t, false, refreshT2.Active)

					accessT2 := persistencesql.OAuth2RequestSQL{Table: "access"}
					require.ErrorIs(t, r.Persister().Connection(s.t2).Where("signature = ?", x.SignatureHash(accessSignatureT2)).First(&accessT2), sql.ErrNoRows)
					require.Equal(t, false, accessT2.Active)
				}
			})

			t.Run("without access signature", func(t *testing.T) {
				clientID := uuid.Must(uuid.NewV4()).String()
				require.NoError(t, r.Persister().CreateClient(s.t1, &client.Client{ID: clientID}))

				request1 := fosite.NewRequest()
				request1.Client = &fosite.DefaultClient{ID: clientID}
				request1.Session = &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "sub"}}

				signature := uuid.Must(uuid.NewV4()).String()
				require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t1, signature, "", request1))

				accessSignature1 := uuid.Must(uuid.NewV4()).String()
				require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1, accessSignature1, request1))

				accessSignature2 := uuid.Must(uuid.NewV4()).String()
				require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1, accessSignature2, request1))

				require.NoError(t, r.Persister().RotateRefreshToken(s.t1, request1.ID, signature))
				{
					accessT1 := persistencesql.OAuth2RequestSQL{Table: "access"}
					require.ErrorIs(t, r.Persister().Connection(s.t1).Where("signature = ?", x.SignatureHash(accessSignature1)).First(&accessT1), sql.ErrNoRows)

					refresh := persistencesql.OAuth2RefreshTable{}
					require.NoError(t, r.Persister().Connection(s.t1).Where("signature = ?", signature).First(&refresh))
					require.Equal(t, false, refresh.Active)

					accessT2 := persistencesql.OAuth2RequestSQL{Table: "access"}
					require.ErrorIs(t, r.Persister().Connection(s.t1).Where("signature = ?", x.SignatureHash(accessSignature2)).First(&accessT2), sql.ErrNoRows)
				}
			})
		})
	}
}

func (s *PersisterTestSuite) TestRevokeSubjectClientConsentSession() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			cl := &client.Client{ID: "client-id"}
			f := newFlow(s.t1NID, cl.ID, "sub", sqlxx.NullString(sessionID))
			f.RequestedAt = time.Now().Add(-24 * time.Hour)
			persistLoginSession(s.t1, t, r.Persister(), &flow.LoginSession{ID: sessionID})
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			require.NoError(t, r.Persister().CreateConsentSession(s.t1, f))

			actual := flow.Flow{}

			require.NoError(t, r.Persister().RevokeSubjectClientConsentSession(s.t2, "sub", cl.ID), "should not error if nothing was found")
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, f.ID))
			require.NoError(t, r.Persister().RevokeSubjectClientConsentSession(s.t1, "sub", cl.ID))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, f.ID))
		})
	}
}

func (s *PersisterTestSuite) TestSetClientAssertionJWT() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			jti := oauth2.NewBlacklistedJTI(uuid.Must(uuid.NewV4()).String(), time.Now().Add(24*time.Hour))
			require.NoError(t, r.Persister().SetClientAssertionJWT(s.t1, jti.JTI, jti.Expiry))

			actual := &oauth2.BlacklistedJTI{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(actual, jti.ID))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestSetClientAssertionJWTRaw() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			store, ok := r.Persister().(*persistencesql.Persister)
			require.True(t, ok)

			jti := oauth2.NewBlacklistedJTI(uuid.Must(uuid.NewV4()).String(), time.Now().Add(24*time.Hour))
			require.NoError(t, store.SetClientAssertionJWTRaw(s.t1, jti))

			actual := &oauth2.BlacklistedJTI{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(actual, jti.ID))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestUpdateClient() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			t1c1 := &client.Client{ID: "client-id", Name: "original", Secret: "original-secret"}
			t2c1 := &client.Client{ID: "client-id", Name: "original", Secret: "original-secret"}
			require.NoError(t, r.Persister().CreateClient(s.t1, t1c1))
			require.NoError(t, r.Persister().CreateClient(s.t2, t2c1))
			t1Hash, t2Hash := t1c1.Secret, t2c1.Secret

			u1 := *t1c1
			u1.Name = "updated"
			u1.Secret = ""
			require.NoError(t, r.Persister().UpdateClient(s.t2, &u1))

			actual, err := r.Persister().GetConcreteClient(s.t1, t1c1.ID)
			require.NoError(t, err)
			require.Equal(t, "original", actual.Name)
			require.Equal(t, t1Hash, actual.Secret)

			actual, err = r.Persister().GetConcreteClient(s.t2, t1c1.ID)
			require.NoError(t, err)
			require.Equal(t, "updated", actual.Name)
			require.Equal(t, t2Hash, actual.Secret)

			u2 := *t1c1
			u2.Name = "updated"
			u2.Secret = ""
			require.NoError(t, r.Persister().UpdateClient(s.t1, &u2))

			actual, err = r.Persister().GetConcreteClient(s.t1, t1c1.ID)
			require.NoError(t, err)
			require.Equal(t, "updated", actual.Name)
			require.Equal(t, t1Hash, actual.Secret)

			u3 := *t1c1
			u3.Name = "updated"
			u3.Secret = "updated-secret"
			require.NoError(t, r.Persister().UpdateClient(s.t1, &u3))

			actual, err = r.Persister().GetConcreteClient(s.t1, t1c1.ID)
			require.NoError(t, err)
			require.Equal(t, "updated", actual.Name)
			require.NotEqual(t, t1Hash, actual.Secret)

			actual, err = r.Persister().GetConcreteClient(s.t2, t2c1.ID)
			require.NoError(t, err)
			require.Equal(t, "updated", actual.Name)
			require.Equal(t, t2Hash, actual.Secret)
		})
	}
}

func (s *PersisterTestSuite) TestUpdateKey() {
	for k, r := range s.registries {
		s.T().Run("dialect="+k, func(t *testing.T) {
			k1 := newKey("test-ks", "test")
			ks := "key-set"
			require.NoError(t, r.KeyManager().AddKey(s.t1, ks, &k1))
			actual, err := r.KeyManager().GetKey(s.t1, ks, k1.KeyID)
			require.NoError(t, err)
			assertx.EqualAsJSON(t, &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{k1}}, actual)

			k2 := newKey("test-ks", "test")
			require.NoError(t, r.KeyManager().UpdateKey(s.t2, ks, &k2))
			actual, err = r.KeyManager().GetKey(s.t1, ks, k1.KeyID)
			require.NoError(t, err)
			assertx.EqualAsJSON(t, &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{k1}}, actual)

			require.NoError(t, r.KeyManager().UpdateKey(s.t1, ks, &k2))
			actual, err = r.KeyManager().GetKey(s.t1, ks, k2.KeyID)
			require.NoError(t, err)
			require.NotEqual(t, &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{k1}}, actual)
		})
	}
}

func (s *PersisterTestSuite) TestUpdateKeySet() {
	for k, r := range s.registries {
		s.T().Run("dialect="+k, func(t *testing.T) {
			ks := "key-set"
			ks1 := newKeySet(ks, "test")
			require.NoError(t, r.KeyManager().AddKeySet(s.t1, ks, ks1))
			actual, err := r.KeyManager().GetKeySet(s.t1, ks)
			require.NoError(t, err)
			requireKeySetEqual(t, ks1, actual)

			ks2 := newKeySet(ks, "test")
			require.NoError(t, r.KeyManager().UpdateKeySet(s.t2, ks, ks2))
			actual, err = r.KeyManager().GetKeySet(s.t1, ks)
			require.NoError(t, err)
			requireKeySetEqual(t, ks1, actual)

			require.NoError(t, r.KeyManager().UpdateKeySet(s.t1, ks, ks2))
			actual, err = r.KeyManager().GetKeySet(s.t1, ks)
			require.NoError(t, err)
			requireKeySetEqual(t, ks2, actual)
		})
	}
}

func (s *PersisterTestSuite) TestUpdateWithNetwork() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			t1c1 := &client.Client{ID: "client-id", Name: "original", Secret: "original-secret"}
			t2c1 := &client.Client{ID: "client-id", Name: "original", Secret: "original-secret", Owner: "erase-me"}
			require.NoError(t, r.Persister().CreateClient(s.t1, t1c1))
			require.NoError(t, r.Persister().CreateClient(s.t2, t2c1))

			store, ok := r.Persister().(*persistencesql.Persister)
			require.True(t, ok)

			count, err := store.UpdateWithNetwork(s.t1, &client.Client{ID: "client-id", Name: "updated", Secret: "original-secret"})
			require.NoError(t, err)
			require.Equal(t, int64(1), count)
			actualt1, err := store.GetConcreteClient(s.t1, "client-id")
			require.NoError(t, err)
			require.Equal(t, "updated", actualt1.Name)
			require.Equal(t, "", actualt1.Owner)
			actualt2, err := store.GetConcreteClient(s.t2, "client-id")
			require.NoError(t, err)
			require.Equal(t, "original", actualt2.Name)
			require.Equal(t, "erase-me", actualt2.Owner)
		})
	}
}

func (s *PersisterTestSuite) TestVerifyAndInvalidateConsentRequest() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			sub := uuid.Must(uuid.NewV4()).String()
			sessionID := uuid.Must(uuid.NewV4()).String()
			persistLoginSession(s.t1, t, r.Persister(), &flow.LoginSession{ID: sessionID})
			cl := &client.Client{ID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, cl))
			f := newFlow(s.t1NID, cl.ID, sub, sqlxx.NullString(sessionID))
			f.ConsentSkip = false
			f.GrantedScope = sqlxx.StringSliceJSONFormat{}
			f.ConsentRemember = false
			crf := 86400
			f.ConsentRememberFor = &crf
			f.ConsentError = &flow.RequestDeniedError{}
			f.SessionAccessToken = map[string]interface{}{}
			f.SessionIDToken = map[string]interface{}{}
			f.State = flow.FlowStateConsentUnused

			require.NoError(t, f.InvalidateConsentRequest())

			err := r.ConsentManager().CreateConsentSession(s.t2, f)
			require.ErrorIs(t, err, sqlcon.ErrNoRows)

			err = r.ConsentManager().CreateConsentSession(s.t1, f)
			require.NoError(t, err)
		})
	}
}

func (s *PersisterTestSuite) TestVerifyAndInvalidateLogoutRequest() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			run := func(t *testing.T, lr *flow.LogoutRequest) {
				lr.Verifier = uuid.Must(uuid.NewV4()).String()
				lr.Accepted = true
				lr.Rejected = false
				require.NoError(t, r.LogoutManager().CreateLogoutRequest(s.t1, lr))

				expected, err := r.LogoutManager().GetLogoutRequest(s.t1, lr.ID)
				require.NoError(t, err)

				lrInvalidated, err := r.LogoutManager().VerifyAndInvalidateLogoutRequest(s.t2, lr.Verifier)
				require.Error(t, err)
				require.Nil(t, lrInvalidated)
				actual := &flow.LogoutRequest{}
				require.NoError(t, r.Persister().Connection(context.Background()).Find(actual, lr.ID))
				require.Equal(t, expected, actual)

				lrInvalidated, err = r.LogoutManager().VerifyAndInvalidateLogoutRequest(s.t1, lr.Verifier)
				require.NoError(t, err)
				require.NoError(t, r.Persister().Connection(context.Background()).Find(actual, lr.ID))
				require.Equal(t, lrInvalidated, actual)
				require.Equal(t, true, actual.WasHandled)
			}

			t.Run("case=legacy logout request without expiry", func(t *testing.T) {
				lr := newLogoutRequest()
				run(t, lr)
			})

			t.Run("case=logout request with expiry", func(t *testing.T) {
				lr := newLogoutRequest()
				lr.ExpiresAt = sqlxx.NullTime(time.Now().Add(time.Hour))
				run(t, lr)
			})

			t.Run("case=logout request that expired returns error", func(t *testing.T) {
				lr := newLogoutRequest()
				lr.ExpiresAt = sqlxx.NullTime(time.Now().UTC().Add(-time.Hour))
				lr.Verifier = uuid.Must(uuid.NewV4()).String()
				lr.Accepted = true
				lr.Rejected = false
				require.NoError(t, r.LogoutManager().CreateLogoutRequest(s.t1, lr))

				_, err := r.LogoutManager().VerifyAndInvalidateLogoutRequest(s.t2, lr.Verifier)
				require.ErrorIs(t, err, x.ErrNotFound)

				_, err = r.LogoutManager().VerifyAndInvalidateLogoutRequest(s.t1, lr.Verifier)
				require.ErrorIs(t, err, flow.ErrorLogoutFlowExpired)
			})
		})
	}
}

func (s *PersisterTestSuite) TestWithFallbackNetworkID() {
	for k, r := range s.registries {
		s.T().Run(k, func(t *testing.T) {
			store1, ok := r.Persister().(*persistencesql.Persister)
			require.True(t, ok)
			original := store1.NetworkID(context.Background())
			expected := uuid.Must(uuid.NewV4())
			store2 := store1.WithFallbackNetworkID(expected)

			assert.NotEqual(t, original, expected)
			assert.Equal(t, expected, store2.NetworkID(context.Background()))
		})
	}
}

func TestPersisterTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(PersisterTestSuite))
}

func newFlow(nid uuid.UUID, clientID string, subject string, sessionID sqlxx.NullString) *flow.Flow {
	return &flow.Flow{
		NID:              nid,
		ID:               uuid.Must(uuid.NewV4()).String(),
		ClientID:         clientID,
		Subject:          subject,
		State:            flow.FlowStateConsentUnused,
		ConsentRequestID: "not-null",
		ConsentCSRF:      "not-null",
		SessionID:        sessionID,
		RequestedAt:      time.Now(),
	}
}

func newLogoutRequest() *flow.LogoutRequest {
	return &flow.LogoutRequest{
		ID: uuid.Must(uuid.NewV4()).String(),
	}
}

func newKey(ksID string, use string) jose.JSONWebKey {
	ks, err := jwk.GenerateJWK(jose.RS256, ksID, use)
	if err != nil {
		panic(err)
	}
	return ks.Keys[0]
}

func newKeySet(id string, use string) *jose.JSONWebKeySet {
	return x.Must(jwk.GenerateJWK(jose.RS256, id, use))
}

func newLoginSession() *flow.LoginSession {
	return &flow.LoginSession{
		ID:              uuid.Must(uuid.NewV4()).String(),
		AuthenticatedAt: sqlxx.NullTime(time.Time{}),
		Subject:         uuid.Must(uuid.NewV4()).String(),
		Remember:        false,
	}
}

func requireKeySetEqual(t *testing.T, expected *jose.JSONWebKeySet, actual *jose.JSONWebKeySet) {
	assertx.EqualAsJSON(t, expected, actual)
}

func persistLoginSession(ctx context.Context, t *testing.T, p persistence.Persister, session *flow.LoginSession) {
	t.Helper()
	require.NoError(t, p.ConfirmLoginSession(ctx, session))
}
