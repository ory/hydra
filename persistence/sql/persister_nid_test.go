// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/ory/hydra/v2/oauth2/flowctx"
	"github.com/ory/hydra/v2/persistence"
	"github.com/ory/x/uuidx"

	"github.com/ory/x/assertx"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/square/go-jose.v2"

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/internal"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/oauth2/trust"
	persistencesql "github.com/ory/hydra/v2/persistence/sql"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/contextx"
	"github.com/ory/x/dbal"
	"github.com/ory/x/networkx"
	"github.com/ory/x/sqlxx"
)

type PersisterTestSuite struct {
	suite.Suite
	registries map[string]driver.Registry
	clean      func(*testing.T)
	t1Ctx      context.Context
	t2Ctx      context.Context
	t1NID      uuid.UUID
	t2NID      uuid.UUID
}

var _ PersisterTestSuite = PersisterTestSuite{}

func (s *PersisterTestSuite) SetupSuite() {
	s.registries = map[string]driver.Registry{
		"memory": internal.NewRegistrySQLFromURL(s.T(), dbal.NewSQLiteTestDatabase(s.T()), true, &contextx.Default{}),
	}

	if !testing.Short() {
		s.registries["postgres"], s.registries["mysql"], s.registries["cockroach"], s.clean = internal.ConnectDatabases(s.T(), true, &contextx.Default{})
	}

	s.t1NID, s.t2NID = uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4())
	s.t1Ctx = flowctx.WithDefaultValues(contextx.SetNIDContext(context.Background(), s.t1NID))
	s.t2Ctx = flowctx.WithDefaultValues(contextx.SetNIDContext(context.Background(), s.t2NID))

	for _, r := range s.registries {
		require.NoError(s.T(), r.Persister().Connection(context.Background()).Create(&networkx.Network{ID: s.t1NID}))
		require.NoError(s.T(), r.Persister().Connection(context.Background()).Create(&networkx.Network{ID: s.t2NID}))
		r.WithContextualizer(&contextx.TestContextualizer{})
	}
}

func (s *PersisterTestSuite) TearDownTest() {
	for _, r := range s.registries {
		r.WithContextualizer(&contextx.TestContextualizer{})
		x.DeleteHydraRows(s.T(), r.Persister().Connection(context.Background()))
	}
}

func (s *PersisterTestSuite) TestAcceptLogoutRequest() {
	t := s.T()
	lr := newLogoutRequest()

	for k, r := range s.registries {
		t.Run("dialect="+k, func(*testing.T) {
			require.NoError(t, r.ConsentManager().CreateLogoutRequest(s.t1Ctx, lr))

			expected, err := r.ConsentManager().GetLogoutRequest(s.t1Ctx, lr.ID)
			require.NoError(t, err)
			require.Equal(t, false, expected.Accepted)

			lrAccepted, err := r.ConsentManager().AcceptLogoutRequest(s.t2Ctx, lr.ID)
			require.Error(t, err)
			require.Equal(t, &flow.LogoutRequest{}, lrAccepted)

			actual, err := r.ConsentManager().GetLogoutRequest(s.t1Ctx, lr.ID)
			require.NoError(t, err)
			require.Equal(t, expected, actual)
		})
	}
}

func (s *PersisterTestSuite) TestAddKeyGetKeyDeleteKey() {
	t := s.T()
	key := newKey("test-ks", "test")
	for k, r := range s.registries {
		t.Run("dialect="+k, func(*testing.T) {
			ks := "key-set"
			require.NoError(t, r.Persister().AddKey(s.t1Ctx, ks, &key))
			actual, err := r.Persister().GetKey(s.t2Ctx, ks, key.KeyID)
			require.Error(t, err)
			require.Equal(t, (*jose.JSONWebKeySet)(nil), actual)
			actual, err = r.Persister().GetKey(s.t1Ctx, ks, key.KeyID)
			require.NoError(t, err)
			assertx.EqualAsJSON(t, &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{key}}, actual)

			r.Persister().DeleteKey(s.t2Ctx, ks, key.KeyID)
			_, err = r.Persister().GetKey(s.t1Ctx, ks, key.KeyID)
			require.NoError(t, err)
			r.Persister().DeleteKey(s.t1Ctx, ks, key.KeyID)
			_, err = r.Persister().GetKey(s.t1Ctx, ks, key.KeyID)
			require.Error(t, err)
		})
	}
}

func (s *PersisterTestSuite) TestAddKeySetGetKeySetDeleteKeySet() {
	t := s.T()
	ks := newKeySet("test-ks", "test")
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			ksID := "key-set"
			r.Persister().AddKeySet(s.t1Ctx, ksID, ks)
			actual, err := r.Persister().GetKeySet(s.t2Ctx, ksID)
			require.Error(t, err)
			require.Equal(t, (*jose.JSONWebKeySet)(nil), actual)
			actual, err = r.Persister().GetKeySet(s.t1Ctx, ksID)
			require.NoError(t, err)
			requireKeySetEqual(t, ks, actual)

			r.Persister().DeleteKeySet(s.t2Ctx, ksID)
			_, err = r.Persister().GetKeySet(s.t1Ctx, ksID)
			require.NoError(t, err)
			r.Persister().DeleteKeySet(s.t1Ctx, ksID)
			_, err = r.Persister().GetKeySet(s.t1Ctx, ksID)
			require.Error(t, err)
		})
	}
}

func (s *PersisterTestSuite) TestAuthenticate() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id", Secret: "secret"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))

			actual, err := r.Persister().Authenticate(s.t2Ctx, "client-id", []byte("secret"))
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().Authenticate(s.t1Ctx, "client-id", []byte("secret"))
			require.NoError(t, err)
			require.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestClientAssertionJWTValid() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			jti := oauth2.NewBlacklistedJTI(uuid.Must(uuid.NewV4()).String(), time.Now().Add(24*time.Hour))
			require.NoError(t, r.Persister().SetClientAssertionJWT(s.t1Ctx, jti.JTI, jti.Expiry))

			require.NoError(t, r.Persister().ClientAssertionJWTValid(s.t2Ctx, jti.JTI))
			require.Error(t, r.Persister().ClientAssertionJWTValid(s.t1Ctx, jti.JTI))
		})
	}
}

func (s *PersisterTestSuite) TestConfirmLoginSession() {
	t := s.T()
	ls := newLoginSession()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			require.NoError(t, r.Persister().CreateLoginSession(s.t1Ctx, ls))

			// Expects the login session to be confirmed in the correct context.
			require.NoError(t, r.Persister().ConfirmLoginSession(s.t1Ctx, nil, ls.ID, time.Now(), ls.Subject, !ls.Remember))
			actual := &flow.LoginSession{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(actual, ls.ID))
			exp, _ := json.Marshal(ls)
			act, _ := json.Marshal(actual)
			require.JSONEq(t, string(exp), string(act))

			// Can't find the login session in the wrong context.
			require.ErrorIs(t,
				r.Persister().ConfirmLoginSession(s.t2Ctx, nil, ls.ID, time.Now(), ls.Subject, !ls.Remember),
				x.ErrNotFound,
			)
		})
	}
}

func (s *PersisterTestSuite) TestCreateSession() {
	t := s.T()
	ls := newLoginSession()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			require.NoError(t, r.Persister().CreateLoginSession(s.t1Ctx, ls))
			require.Equal(t, s.t1NID, ls.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCountClients() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			count, err := r.Persister().CountClients(s.t1Ctx)
			require.NoError(t, err)
			require.Equal(t, 0, count)

			count, err = r.Persister().CountClients(s.t2Ctx)
			require.NoError(t, err)
			require.Equal(t, 0, count)

			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, newClient()))

			count, err = r.Persister().CountClients(s.t1Ctx)
			require.NoError(t, err)
			require.Equal(t, 1, count)

			count, err = r.Persister().CountClients(s.t2Ctx)
			require.NoError(t, err)
			require.Equal(t, 0, count)
		})
	}
}

func (s *PersisterTestSuite) TestCountGrants() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			count, err := r.Persister().CountGrants(s.t1Ctx)
			require.NoError(t, err)
			require.Equal(t, 0, count)

			count, err = r.Persister().CountGrants(s.t2Ctx)
			require.NoError(t, err)
			require.Equal(t, 0, count)

			keySet := uuid.Must(uuid.NewV4()).String()
			publicKey := newKey(keySet, "use")
			grant := newGrant(keySet, publicKey.KeyID)
			require.NoError(t, r.Persister().AddKey(s.t1Ctx, keySet, &publicKey))
			require.NoError(t, r.Persister().CreateGrant(s.t1Ctx, grant, publicKey))

			count, err = r.Persister().CountGrants(s.t1Ctx)
			require.NoError(t, err)
			require.Equal(t, 1, count)

			count, err = r.Persister().CountGrants(s.t2Ctx)
			require.NoError(t, err)
			require.Equal(t, 0, count)
		})
	}
}

func (s *PersisterTestSuite) TestCountSubjectsGrantedConsentRequests() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			sub := uuid.Must(uuid.NewV4()).String()
			count, err := r.Persister().CountSubjectsGrantedConsentRequests(s.t1Ctx, sub)
			require.NoError(t, err)
			require.Equal(t, 0, count)

			count, err = r.Persister().CountSubjectsGrantedConsentRequests(s.t2Ctx, sub)
			require.NoError(t, err)
			require.Equal(t, 0, count)

			sessionID := uuid.Must(uuid.NewV4()).String()
			persistLoginSession(s.t1Ctx, t, r.Persister(), &flow.LoginSession{ID: sessionID})
			client := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			f := newFlow(s.t1NID, client.LegacyClientID, sub, sqlxx.NullString(sessionID))
			f.ConsentSkip = false
			f.ConsentError = &flow.RequestDeniedError{}
			f.State = flow.FlowStateConsentUnused
			require.NoError(t, r.Persister().Connection(context.Background()).Create(f))

			count, err = r.Persister().CountSubjectsGrantedConsentRequests(s.t1Ctx, sub)
			require.NoError(t, err)
			require.Equal(t, 1, count)

			count, err = r.Persister().CountSubjectsGrantedConsentRequests(s.t2Ctx, sub)
			require.NoError(t, err)
			require.Equal(t, 0, count)
		})
	}
}

func (s *PersisterTestSuite) TestCreateAccessTokenSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			c1 := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, c1))
			c2 := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t2Ctx, c2))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()

			fr.Client = &fosite.DefaultClient{ID: c1.LegacyClientID}
			require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1Ctx, sig, fr))
			actual := persistencesql.OAuth2RequestSQL{Table: "access"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, persistencesql.SignatureHash(sig)))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateAuthorizeCodeSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			c1 := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, c1))
			c2 := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t2Ctx, c2))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: c1.LegacyClientID}
			require.NoError(t, r.Persister().CreateAuthorizeCodeSession(s.t1Ctx, sig, fr))
			actual := persistencesql.OAuth2RequestSQL{Table: "code"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, sig))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateClient() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			expected := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, expected))
			actual := client.Client{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, expected.ID))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateConsentRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			client := &client.Client{LegacyClientID: "client-id"}
			f := newFlow(s.t1NID, client.LegacyClientID, "sub", sqlxx.NullString(sessionID))
			ls := &flow.LoginSession{ID: sessionID}
			require.NoError(t, r.Persister().CreateLoginSession(s.t1Ctx, ls))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(ls))
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(f))

			req := &flow.OAuth2ConsentRequest{
				ID:             "consent-request-id",
				LoginChallenge: sqlxx.NullString(f.ID),
				Skip:           false,
				Verifier:       "verifier",
				CSRF:           "csrf",
			}
			require.NoError(t, r.Persister().CreateConsentRequest(s.t1Ctx, nil, req))

			actual := flow.Flow{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, f.ID))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateForcedObfuscatedLoginSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			session := &consent.ForcedObfuscatedLoginSession{ClientID: client.LegacyClientID}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			require.NoError(t, r.Persister().CreateForcedObfuscatedLoginSession(s.t1Ctx, session))
			actual, err := r.Persister().GetForcedObfuscatedLoginSession(s.t1Ctx, client.LegacyClientID, "")
			require.NoError(t, err)
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateGrant() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()).String(),
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().AddKeySet(s.t1Ctx, "ks-id", ks))
			require.NoError(t, r.Persister().CreateGrant(s.t1Ctx, grant, ks.Keys[0]))
			actual := trust.SQLData{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, grant.ID))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateLoginRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			lr := flow.LoginRequest{ID: "lr-id", ClientID: client.LegacyClientID, RequestedAt: time.Now()}

			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			f, err := r.ConsentManager().CreateLoginRequest(s.t1Ctx, &lr)
			require.NoError(t, err)
			require.Equal(t, s.t1NID, f.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateLoginSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			ls := flow.LoginSession{ID: uuid.Must(uuid.NewV4()).String(), Remember: true}
			require.NoError(t, r.Persister().CreateLoginSession(s.t1Ctx, &ls))
			actual, err := r.Persister().GetRememberedLoginSession(s.t1Ctx, nil, ls.ID)
			require.NoError(t, err)
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateLogoutRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			lr := flow.LogoutRequest{
				// TODO there is not FK for SessionID so we don't need it here; TODO make sure the missing FK is intentional
				ID:       uuid.Must(uuid.NewV4()).String(),
				ClientID: sql.NullString{Valid: true, String: client.LegacyClientID},
			}

			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			require.NoError(t, r.Persister().CreateLogoutRequest(s.t1Ctx, &lr))
			actual, err := r.Persister().GetLogoutRequest(s.t1Ctx, lr.ID)
			require.NoError(t, err)
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateOpenIDConnectSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}

			authorizeCode := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateOpenIDConnectSession(s.t1Ctx, authorizeCode, request))

			actual := persistencesql.OAuth2RequestSQL{Table: "oidc"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreatePKCERequestSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}

			authorizeCode := uuid.Must(uuid.NewV4()).String()

			actual := persistencesql.OAuth2RequestSQL{Table: "pkce"}
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.NoError(t, r.Persister().CreatePKCERequestSession(s.t1Ctx, authorizeCode, request))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateRefreshTokenSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}

			authorizeCode := uuid.Must(uuid.NewV4()).String()
			actual := persistencesql.OAuth2RequestSQL{Table: "refresh"}
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t1Ctx, authorizeCode, request))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateWithNetwork() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			expected := &client.Client{LegacyClientID: "client-id"}
			store, ok := r.OAuth2Storage().(*persistencesql.Persister)
			if !ok {
				t.Fatal("type assertion failed")
			}
			store.CreateWithNetwork(s.t1Ctx, expected)

			actual := &client.Client{}
			require.NoError(t, r.Persister().Connection(context.Background()).Where("id = ?", expected.LegacyClientID).First(actual))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) DeleteAccessTokenSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: client.LegacyClientID}
			require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1Ctx, sig, fr))
			require.NoError(t, r.Persister().DeleteAccessTokenSession(s.t2Ctx, sig))

			actual := persistencesql.OAuth2RequestSQL{Table: "access"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, sig))
			require.Equal(t, s.t1NID, actual.NID)

			require.NoError(t, r.Persister().DeleteAccessTokenSession(s.t1Ctx, sig))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, sig))
		})
	}
}

func (s *PersisterTestSuite) TestDeleteAccessTokens() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: client.LegacyClientID}
			require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1Ctx, sig, fr))
			require.NoError(t, r.Persister().DeleteAccessTokens(s.t2Ctx, client.LegacyClientID))

			actual := persistencesql.OAuth2RequestSQL{Table: "access"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, persistencesql.SignatureHash(sig)))
			require.Equal(t, s.t1NID, actual.NID)

			require.NoError(t, r.Persister().DeleteAccessTokens(s.t1Ctx, client.LegacyClientID))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, persistencesql.SignatureHash(sig)))
		})
	}
}

func (s *PersisterTestSuite) TestDeleteClient() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			c := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, c))
			actual := client.Client{}
			require.Error(t, r.Persister().DeleteClient(s.t2Ctx, c.LegacyClientID))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, c.ID))
			require.NoError(t, r.Persister().DeleteClient(s.t1Ctx, c.LegacyClientID))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, c.ID))
		})
	}
}

func (s *PersisterTestSuite) TestDeleteGrant() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()).String(),
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().AddKeySet(s.t1Ctx, "ks-id", ks))
			require.NoError(t, r.Persister().CreateGrant(s.t1Ctx, grant, ks.Keys[0]))

			actual := trust.SQLData{}
			require.Error(t, r.Persister().DeleteGrant(s.t2Ctx, grant.ID))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, grant.ID))
			require.NoError(t, r.Persister().DeleteGrant(s.t1Ctx, grant.ID))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, grant.ID))
		})
	}
}

func (s *PersisterTestSuite) TestDeleteLoginSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			ls := flow.LoginSession{ID: uuid.Must(uuid.NewV4()).String(), Remember: true}
			persistLoginSession(s.t1Ctx, t, r.Persister(), &ls)

			require.Error(t, r.Persister().DeleteLoginSession(s.t2Ctx, ls.ID))
			_, err := r.Persister().GetRememberedLoginSession(s.t1Ctx, nil, ls.ID)
			require.NoError(t, err)

			require.NoError(t, r.Persister().DeleteLoginSession(s.t1Ctx, ls.ID))
			_, err = r.Persister().GetRememberedLoginSession(s.t1Ctx, nil, ls.ID)
			require.Error(t, err)
		})
	}
}

func (s *PersisterTestSuite) TestDeleteOpenIDConnectSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}

			authorizeCode := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateOpenIDConnectSession(s.t1Ctx, authorizeCode, request))

			actual := persistencesql.OAuth2RequestSQL{Table: "oidc"}

			require.NoError(t, r.Persister().DeleteOpenIDConnectSession(s.t2Ctx, authorizeCode))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.NoError(t, r.Persister().DeleteOpenIDConnectSession(s.t1Ctx, authorizeCode))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
		})
	}
}

func (s *PersisterTestSuite) TestDeletePKCERequestSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}

			authorizeCode := uuid.Must(uuid.NewV4()).String()
			r.Persister().CreatePKCERequestSession(s.t1Ctx, authorizeCode, request)

			actual := persistencesql.OAuth2RequestSQL{Table: "pkce"}

			require.NoError(t, r.Persister().DeletePKCERequestSession(s.t2Ctx, authorizeCode))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.NoError(t, r.Persister().DeletePKCERequestSession(s.t1Ctx, authorizeCode))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
		})
	}
}

func (s *PersisterTestSuite) TestDeleteRefreshTokenSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}

			signature := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t1Ctx, signature, request))

			actual := persistencesql.OAuth2RequestSQL{Table: "refresh"}

			require.NoError(t, r.Persister().DeleteRefreshTokenSession(s.t2Ctx, signature))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, signature))
			require.NoError(t, r.Persister().DeleteRefreshTokenSession(s.t1Ctx, signature))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, signature))
		})
	}
}

func (s *PersisterTestSuite) TestDetermineNetwork() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			store, ok := r.OAuth2Storage().(*persistencesql.Persister)
			if !ok {
				t.Fatal("type assertion failed")
			}

			r.Persister().Connection(context.Background()).Where("id <> ? AND id <> ?", s.t1NID, s.t2NID).Delete(&networkx.Network{})

			actual, err := store.DetermineNetwork(context.Background())
			require.NoError(t, err)
			require.True(t, actual.ID == s.t1NID || actual.ID == s.t2NID)
		})
	}
}

func (s *PersisterTestSuite) TestFindGrantedAndRememberedConsentRequests() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			client := &client.Client{LegacyClientID: "client-id"}
			f := newFlow(s.t1NID, client.LegacyClientID, "sub", sqlxx.NullString(sessionID))
			persistLoginSession(s.t1Ctx, t, r.Persister(), &flow.LoginSession{ID: sessionID})
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))

			req := &flow.OAuth2ConsentRequest{
				ID:             "consent-request-id",
				LoginChallenge: sqlxx.NullString(f.ID),
				Skip:           false,
				Verifier:       "verifier",
				CSRF:           "csrf",
			}

			hcr := &flow.AcceptOAuth2ConsentRequest{
				ID:        x.Must(flowctx.Encode(s.t1Ctx, r.KeyCipher(), f)),
				HandledAt: sqlxx.NullTime(time.Now()),
				Remember:  true,
			}
			require.NoError(t, r.Persister().CreateConsentRequest(s.t1Ctx, nil, req))
			_, err := r.Persister().HandleConsentRequest(s.t1Ctx, nil, hcr)
			require.NoError(t, err)

			actual, err := r.Persister().FindGrantedAndRememberedConsentRequests(s.t2Ctx, client.LegacyClientID, f.Subject)
			require.Error(t, err)
			require.Equal(t, 0, len(actual))

			actual, err = r.Persister().FindGrantedAndRememberedConsentRequests(s.t1Ctx, client.LegacyClientID, f.Subject)
			require.NoError(t, err)
			require.Equal(t, 1, len(actual))
		})
	}
}

func (s *PersisterTestSuite) TestFindSubjectsGrantedConsentRequests() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			client := &client.Client{LegacyClientID: "client-id"}
			f := newFlow(s.t1NID, client.LegacyClientID, "sub", sqlxx.NullString(sessionID))
			persistLoginSession(s.t1Ctx, t, r.Persister(), &flow.LoginSession{ID: sessionID})
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))

			req := &flow.OAuth2ConsentRequest{
				ID:             "consent-request-id",
				LoginChallenge: sqlxx.NullString(f.ID),
				Skip:           false,
				Verifier:       "verifier",
				CSRF:           "csrf",
			}

			hcr := &flow.AcceptOAuth2ConsentRequest{
				ID:        x.Must(flowctx.Encode(s.t1Ctx, r.KeyCipher(), f)),
				HandledAt: sqlxx.NullTime(time.Now()),
				Remember:  true,
			}
			require.NoError(t, r.Persister().CreateConsentRequest(s.t1Ctx, nil, req))
			_, err := r.Persister().HandleConsentRequest(s.t1Ctx, nil, hcr)
			require.NoError(t, err)

			actual, err := r.Persister().FindSubjectsGrantedConsentRequests(s.t2Ctx, f.Subject, 100, 0)
			require.Error(t, err)
			require.Equal(t, 0, len(actual))

			actual, err = r.Persister().FindSubjectsGrantedConsentRequests(s.t1Ctx, f.Subject, 100, 0)
			require.NoError(t, err)
			require.Equal(t, 1, len(actual))
		})
	}
}

func (s *PersisterTestSuite) TestFlushInactiveAccessTokens() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.RequestedAt = time.Now().UTC().Add(-24 * time.Hour)
			fr.Client = &fosite.DefaultClient{ID: client.LegacyClientID}
			require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1Ctx, sig, fr))

			actual := persistencesql.OAuth2RequestSQL{Table: "access"}

			require.NoError(t, r.Persister().FlushInactiveAccessTokens(s.t2Ctx, time.Now().Add(time.Hour), 100, 100))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, persistencesql.SignatureHash(sig)))
			require.NoError(t, r.Persister().FlushInactiveAccessTokens(s.t1Ctx, time.Now().Add(time.Hour), 100, 100))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, persistencesql.SignatureHash(sig)))
		})
	}
}

func (s *PersisterTestSuite) TestGenerateAndPersistKeySet() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			store, ok := r.OAuth2Storage().(*persistencesql.Persister)
			if !ok {
				t.Fatal("type assertion failed")
			}

			actual := &jwk.SQLData{}

			ks, err := store.GenerateAndPersistKeySet(s.t1Ctx, "ks", "kid", "RS256", "use")
			require.NoError(t, err)
			require.Error(t, r.Persister().Connection(context.Background()).Where("sid = ? AND kid = ? AND nid = ?", "ks", ks.Keys[0].KeyID, s.t2NID).First(actual))
			require.NoError(t, r.Persister().Connection(context.Background()).Where("sid = ? AND kid = ? AND nid = ?", "ks", ks.Keys[0].KeyID, s.t1NID).First(actual))
		})
	}
}

func (s *PersisterTestSuite) TestFlushInactiveGrants() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()).String(),
				ExpiresAt: time.Now().Add(-24 * time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().AddKeySet(s.t1Ctx, "ks-id", ks))
			require.NoError(t, r.Persister().CreateGrant(s.t1Ctx, grant, ks.Keys[0]))

			actual := trust.SQLData{}
			require.NoError(t, r.Persister().FlushInactiveGrants(s.t2Ctx, time.Now(), 100, 100))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, grant.ID))
			require.NoError(t, r.Persister().FlushInactiveGrants(s.t1Ctx, time.Now(), 100, 100))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, grant.ID))
		})
	}
}

func (s *PersisterTestSuite) TestFlushInactiveLoginConsentRequests() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			client := &client.Client{LegacyClientID: "client-id"}
			f := newFlow(s.t1NID, client.LegacyClientID, "sub", sqlxx.NullString(sessionID))
			f.RequestedAt = time.Now().Add(-24 * time.Hour)
			persistLoginSession(s.t1Ctx, t, r.Persister(), &flow.LoginSession{ID: sessionID})
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(f))

			actual := flow.Flow{}

			require.NoError(t, r.Persister().FlushInactiveLoginConsentRequests(s.t2Ctx, time.Now(), 100, 100))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, f.ID))
			require.NoError(t, r.Persister().FlushInactiveLoginConsentRequests(s.t1Ctx, time.Now(), 100, 100))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, f.ID))
		})
	}
}

func (s *PersisterTestSuite) TestFlushInactiveRefreshTokens() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			request := fosite.NewRequest()
			request.RequestedAt = time.Now().Add(-240 * 365 * time.Hour)
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			signature := uuid.Must(uuid.NewV4()).String()

			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t1Ctx, signature, request))

			actual := persistencesql.OAuth2RequestSQL{Table: "refresh"}

			require.NoError(t, r.Persister().FlushInactiveRefreshTokens(s.t2Ctx, time.Now(), 100, 100))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, signature))
			require.NoError(t, r.Persister().FlushInactiveRefreshTokens(s.t1Ctx, time.Now(), 100, 100))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, signature))
		})
	}
}

func (s *PersisterTestSuite) TestGetAccessTokenSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: client.LegacyClientID}
			require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1Ctx, sig, fr))

			actual, err := r.Persister().GetAccessTokenSession(s.t2Ctx, sig, &fosite.DefaultSession{})
			require.Error(t, err)
			require.Nil(t, actual)
			actual, err = r.Persister().GetAccessTokenSession(s.t1Ctx, sig, &fosite.DefaultSession{})
			require.NoError(t, err)
			require.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetAuthorizeCodeSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: client.LegacyClientID}
			require.NoError(t, r.Persister().CreateAuthorizeCodeSession(s.t1Ctx, sig, fr))

			actual, err := r.Persister().GetAuthorizeCodeSession(s.t2Ctx, sig, &fosite.DefaultSession{})
			require.Error(t, err)
			require.Nil(t, actual)
			actual, err = r.Persister().GetAuthorizeCodeSession(s.t1Ctx, sig, &fosite.DefaultSession{})
			require.NoError(t, err)
			require.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetClient() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			expected := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, expected))

			actual, err := r.Persister().GetClient(s.t2Ctx, expected.LegacyClientID)
			require.Error(t, err)
			require.Nil(t, actual)
			actual, err = r.Persister().GetClient(s.t1Ctx, expected.LegacyClientID)
			require.NoError(t, err)
			require.Equal(t, expected.LegacyClientID, actual.GetID())
		})
	}
}

func (s *PersisterTestSuite) TestGetClientAssertionJWT() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			store, ok := r.OAuth2Storage().(oauth2.AssertionJWTReader)
			if !ok {
				t.Fatal("type assertion failed")
			}
			expected := oauth2.NewBlacklistedJTI(uuid.Must(uuid.NewV4()).String(), time.Now().Add(24*time.Hour))
			require.NoError(t, r.Persister().SetClientAssertionJWT(s.t1Ctx, expected.JTI, expected.Expiry))

			_, err := store.GetClientAssertionJWT(s.t2Ctx, expected.JTI)
			require.Error(t, err)
			_, err = store.GetClientAssertionJWT(s.t1Ctx, expected.JTI)
			require.NoError(t, err)
		})
	}
}

func (s *PersisterTestSuite) TestGetClients() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			c := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, c))

			actual, err := r.Persister().GetClients(s.t2Ctx, client.Filter{Offset: 0, Limit: 100})
			require.NoError(t, err)
			require.Equal(t, 0, len(actual))
			actual, err = r.Persister().GetClients(s.t1Ctx, client.Filter{Offset: 0, Limit: 100})
			require.NoError(t, err)
			require.Equal(t, 1, len(actual))
		})
	}
}

func (s *PersisterTestSuite) TestGetConcreteClient() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			expected := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, expected))

			actual, err := r.Persister().GetConcreteClient(s.t2Ctx, expected.LegacyClientID)
			require.Error(t, err)
			require.Nil(t, actual)
			actual, err = r.Persister().GetConcreteClient(s.t1Ctx, expected.LegacyClientID)
			require.NoError(t, err)
			require.Equal(t, expected.LegacyClientID, actual.GetID())
		})
	}
}

func (s *PersisterTestSuite) TestGetConcreteGrant() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()).String(),
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().AddKeySet(s.t1Ctx, "ks-id", ks))
			require.NoError(t, r.Persister().CreateGrant(s.t1Ctx, grant, ks.Keys[0]))

			actual, err := r.Persister().GetConcreteGrant(s.t2Ctx, grant.ID)
			require.Error(t, err)
			require.Equal(t, trust.Grant{}, actual)

			actual, err = r.Persister().GetConcreteGrant(s.t1Ctx, grant.ID)
			require.NoError(t, err)
			require.NotEqual(t, trust.Grant{}, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetConsentRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			client := &client.Client{LegacyClientID: "client-id"}
			f := newFlow(s.t1NID, client.LegacyClientID, "sub", sqlxx.NullString(sessionID))
			persistLoginSession(s.t1Ctx, t, r.Persister(), &flow.LoginSession{ID: sessionID})
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))

			req := &flow.OAuth2ConsentRequest{
				ID:             x.Must(flowctx.Encode(s.t1Ctx, r.KeyCipher(), f)),
				LoginChallenge: sqlxx.NullString(f.ID),
				Skip:           false,
				Verifier:       "verifier",
				CSRF:           "csrf",
			}
			require.NoError(t, r.Persister().CreateConsentRequest(s.t1Ctx, nil, req))

			actual, err := r.Persister().GetConsentRequest(s.t2Ctx, req.ID)
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetConsentRequest(s.t1Ctx, req.ID)
			require.NoError(t, err)
			require.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetFlow() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			client := &client.Client{LegacyClientID: "client-id"}
			f := newFlow(s.t1NID, client.LegacyClientID, "sub", sqlxx.NullString(sessionID))
			persistLoginSession(s.t1Ctx, t, r.Persister(), &flow.LoginSession{ID: sessionID})
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(f))

			store, ok := r.Persister().(*persistencesql.Persister)
			if !ok {
				t.Fatal("type assertion failed")
			}

			_, err := store.GetFlow(s.t2Ctx, f.ID)
			require.Error(t, err)

			_, err = store.GetFlow(s.t1Ctx, f.ID)
			require.NoError(t, err)
		})
	}
}

func (s *PersisterTestSuite) TestGetFlowByConsentChallenge() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			client := &client.Client{LegacyClientID: "client-id"}
			f := newFlow(s.t1NID, client.LegacyClientID, "sub", sqlxx.NullString(sessionID))
			persistLoginSession(s.t1Ctx, t, r.Persister(), &flow.LoginSession{ID: sessionID})
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(f))

			store, ok := r.Persister().(*persistencesql.Persister)
			if !ok {
				t.Fatal("type assertion failed")
			}

			_, err := store.GetFlowByConsentChallenge(s.t2Ctx, x.Must(flowctx.Encode(s.t2Ctx, r.KeyCipher(), f)))
			require.Error(t, err)

			_, err = store.GetFlowByConsentChallenge(s.t1Ctx, x.Must(flowctx.Encode(s.t2Ctx, r.KeyCipher(), f)))
			require.NoError(t, err)
		})
	}
}

func (s *PersisterTestSuite) TestGetForcedObfuscatedLoginSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			session := &consent.ForcedObfuscatedLoginSession{ClientID: client.LegacyClientID}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			require.NoError(t, r.Persister().CreateForcedObfuscatedLoginSession(s.t1Ctx, session))

			actual, err := r.Persister().GetForcedObfuscatedLoginSession(s.t2Ctx, client.LegacyClientID, "")
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetForcedObfuscatedLoginSession(s.t1Ctx, client.LegacyClientID, "")
			require.NoError(t, err)
			require.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetGrants() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()).String(),
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().AddKeySet(s.t1Ctx, "ks-id", ks))
			require.NoError(t, r.Persister().CreateGrant(s.t1Ctx, grant, ks.Keys[0]))

			actual, err := r.Persister().GetGrants(s.t2Ctx, 100, 0, "")
			require.NoError(t, err)
			require.Equal(t, 0, len(actual))

			actual, err = r.Persister().GetGrants(s.t1Ctx, 100, 0, "")
			require.NoError(t, err)
			require.Equal(t, 1, len(actual))
		})
	}
}

func (s *PersisterTestSuite) TestGetLoginRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			lr := flow.LoginRequest{ID: "lr-id", ClientID: client.LegacyClientID, RequestedAt: time.Now()}

			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			f, err := r.ConsentManager().CreateLoginRequest(s.t1Ctx, &lr)
			require.NoError(t, err)
			lr.ID = x.Must(flowctx.Encode(s.t1Ctx, r.KeyCipher(), f))
			require.Equal(t, s.t1NID, f.NID)

			actual, err := r.Persister().GetLoginRequest(s.t2Ctx, lr.ID)
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetLoginRequest(s.t1Ctx, lr.ID)
			require.NoError(t, err)
			require.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetLogoutRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			lr := flow.LogoutRequest{
				ID:       uuid.Must(uuid.NewV4()).String(),
				ClientID: sql.NullString{Valid: true, String: client.LegacyClientID},
			}

			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			require.NoError(t, r.Persister().CreateLogoutRequest(s.t1Ctx, &lr))

			actual, err := r.Persister().GetLogoutRequest(s.t2Ctx, lr.ID)
			require.Error(t, err)
			require.Equal(t, &flow.LogoutRequest{}, actual)

			actual, err = r.Persister().GetLogoutRequest(s.t1Ctx, lr.ID)
			require.NoError(t, err)
			require.NotEqual(t, &flow.LogoutRequest{}, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetOpenIDConnectSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			request := fosite.NewRequest()
			request.SetID("request-id")
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			authorizeCode := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			require.NoError(t, r.Persister().CreateOpenIDConnectSession(s.t1Ctx, authorizeCode, request))

			actual, err := r.Persister().GetOpenIDConnectSession(s.t2Ctx, authorizeCode, &fosite.Request{})
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetOpenIDConnectSession(s.t1Ctx, authorizeCode, &fosite.Request{})
			require.NoError(t, err)
			require.Equal(t, request.GetID(), actual.GetID())
		})
	}
}

func (s *PersisterTestSuite) TestGetPKCERequestSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			request := fosite.NewRequest()
			request.SetID("request-id")
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			sig := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			require.NoError(t, r.Persister().CreatePKCERequestSession(s.t1Ctx, sig, request))

			actual, err := r.Persister().GetPKCERequestSession(s.t2Ctx, sig, &fosite.DefaultSession{})
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetPKCERequestSession(s.t1Ctx, sig, &fosite.DefaultSession{})
			require.NoError(t, err)
			require.Equal(t, request.GetID(), actual.GetID())
		})
	}
}

func (s *PersisterTestSuite) TestGetPublicKey() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()).String(),
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().AddKeySet(s.t1Ctx, "ks-id", ks))
			require.NoError(t, r.Persister().CreateGrant(s.t1Ctx, grant, ks.Keys[0]))

			actual, err := r.Persister().GetPublicKey(s.t2Ctx, grant.Issuer, grant.Subject, grant.PublicKey.KeyID)
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetPublicKey(s.t1Ctx, grant.Issuer, grant.Subject, grant.PublicKey.KeyID)
			require.NoError(t, err)
			require.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetPublicKeyScopes() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()).String(),
				Scope:     []string{"a", "b", "c"},
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().AddKeySet(s.t1Ctx, "ks-id", ks))
			require.NoError(t, r.Persister().CreateGrant(s.t1Ctx, grant, ks.Keys[0]))

			actual, err := r.Persister().GetPublicKeyScopes(s.t2Ctx, grant.Issuer, grant.Subject, grant.PublicKey.KeyID)
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetPublicKeyScopes(s.t1Ctx, grant.Issuer, grant.Subject, grant.PublicKey.KeyID)
			require.NoError(t, err)
			require.Equal(t, grant.Scope, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetPublicKeys() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()).String(),
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().AddKeySet(s.t1Ctx, "ks-id", ks))
			require.NoError(t, r.Persister().CreateGrant(s.t1Ctx, grant, ks.Keys[0]))

			actual, err := r.Persister().GetPublicKeys(s.t2Ctx, grant.Issuer, grant.Subject)
			require.NoError(t, err)
			require.Nil(t, actual.Keys)

			actual, err = r.Persister().GetPublicKeys(s.t1Ctx, grant.Issuer, grant.Subject)
			require.NoError(t, err)
			require.NotNil(t, actual.Keys)
		})
	}
}

func (s *PersisterTestSuite) TestGetRefreshTokenSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			cl := &client.Client{LegacyClientID: "client-id"}
			request := fosite.NewRequest()
			request.SetID("request-id")
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			sig := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, cl))
			require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t1Ctx, sig, request))

			actual, err := r.Persister().GetRefreshTokenSession(s.t2Ctx, sig, &fosite.DefaultSession{})
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetRefreshTokenSession(s.t1Ctx, sig, &fosite.DefaultSession{})
			require.NoError(t, err)
			require.Equal(t, request.GetID(), actual.GetID())
		})
	}
}

func (s *PersisterTestSuite) TestGetRememberedLoginSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			ls := flow.LoginSession{ID: uuid.Must(uuid.NewV4()).String(), Remember: true}
			require.NoError(t, r.Persister().CreateLoginSession(s.t1Ctx, &ls))

			actual, err := r.Persister().GetRememberedLoginSession(s.t2Ctx, nil, ls.ID)
			assert.Error(t, err)
			assert.Nil(t, actual)

			actual, err = r.Persister().GetRememberedLoginSession(s.t1Ctx, nil, ls.ID)
			assert.NoError(t, err)
			assert.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestHandleConsentRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			c1 := &client.Client{LegacyClientID: uuidx.NewV4().String()}
			f := newFlow(s.t1NID, c1.LegacyClientID, "sub", sqlxx.NullString(sessionID))
			persistLoginSession(s.t1Ctx, t, r.Persister(), &flow.LoginSession{ID: sessionID})
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, c1))
			c1.ID = uuid.Nil
			require.NoError(t, r.Persister().CreateClient(s.t2Ctx, c1))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(f))

			req := &flow.OAuth2ConsentRequest{
				ID:             "consent-request-id",
				LoginChallenge: sqlxx.NullString(f.ID),
				Skip:           false,
				Verifier:       "verifier",
				CSRF:           "csrf",
			}
			require.NoError(t, r.Persister().CreateConsentRequest(s.t1Ctx, nil, req))

			hcr := &flow.AcceptOAuth2ConsentRequest{
				ID:        x.Must(flowctx.Encode(s.t1Ctx, r.KeyCipher(), f)),
				HandledAt: sqlxx.NullTime(time.Now()),
				Remember:  true,
			}

			actualCR, err := r.Persister().HandleConsentRequest(s.t2Ctx, nil, hcr)
			require.Error(t, err)
			require.Nil(t, actualCR)
			actual, err := r.Persister().FindGrantedAndRememberedConsentRequests(s.t1Ctx, c1.LegacyClientID, f.Subject)
			require.Error(t, err)
			require.Equal(t, 0, len(actual))

			actualCR, err = r.Persister().HandleConsentRequest(s.t1Ctx, nil, hcr)
			require.NoError(t, err)
			require.NotNil(t, actualCR)
			actual, err = r.Persister().FindGrantedAndRememberedConsentRequests(s.t1Ctx, c1.LegacyClientID, f.Subject)
			require.NoError(t, err)
			require.Equal(t, 1, len(actual))
		})
	}
}

func (s *PersisterTestSuite) TestInvalidateAuthorizeCodeSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			cl := &client.Client{LegacyClientID: uuidx.NewV4().String()}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, cl))
			cl.ID = uuid.Nil
			require.NoError(t, r.Persister().CreateClient(s.t2Ctx, cl))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: cl.LegacyClientID}
			require.NoError(t, r.Persister().CreateAuthorizeCodeSession(s.t1Ctx, sig, fr))

			require.NoError(t, r.Persister().InvalidateAuthorizeCodeSession(s.t2Ctx, sig))
			actual := persistencesql.OAuth2RequestSQL{Table: "code"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, sig))
			require.Equal(t, true, actual.Active)

			require.NoError(t, r.Persister().InvalidateAuthorizeCodeSession(s.t1Ctx, sig))
			actual = persistencesql.OAuth2RequestSQL{Table: "code"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, sig))
			require.Equal(t, false, actual.Active)
		})
	}
}

func (s *PersisterTestSuite) TestIsJWTUsed() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			jti := oauth2.NewBlacklistedJTI(uuid.Must(uuid.NewV4()).String(), time.Now().Add(24*time.Hour))
			require.NoError(t, r.Persister().SetClientAssertionJWT(s.t1Ctx, jti.JTI, jti.Expiry))

			actual, err := r.Persister().IsJWTUsed(s.t2Ctx, jti.JTI)
			require.NoError(t, err)
			require.False(t, actual)

			actual, err = r.Persister().IsJWTUsed(s.t1Ctx, jti.JTI)
			require.NoError(t, err)
			require.True(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestListUserAuthenticatedClientsWithBackChannelLogout() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			c1 := &client.Client{LegacyClientID: "client-1", BackChannelLogoutURI: "not-null"}
			c2 := &client.Client{LegacyClientID: "client-2", BackChannelLogoutURI: "not-null"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, c1))
			c1.ID = uuid.Nil
			require.NoError(t, r.Persister().CreateClient(s.t2Ctx, c1))
			require.NoError(t, r.Persister().CreateClient(s.t2Ctx, c2))

			t1f1 := newFlow(s.t1NID, c1.LegacyClientID, "sub", sqlxx.NullString(uuid.Must(uuid.NewV4()).String()))
			t1f1.ConsentChallengeID = "t1f1-consent-challenge"
			t1f1.LoginVerifier = "t1f1-login-verifier"
			t1f1.ConsentVerifier = "t1f1-consent-verifier"

			t2f1 := newFlow(s.t2NID, c1.LegacyClientID, "sub", t1f1.SessionID)
			t2f1.ConsentChallengeID = "t2f1-consent-challenge"
			t2f1.LoginVerifier = "t2f1-login-verifier"
			t2f1.ConsentVerifier = "t2f1-consent-verifier"

			t2f2 := newFlow(s.t2NID, c2.LegacyClientID, "sub", t1f1.SessionID)
			t2f2.ConsentChallengeID = "t2f2-consent-challenge"
			t2f2.LoginVerifier = "t2f2-login-verifier"
			t2f2.ConsentVerifier = "t2f2-consent-verifier"

			persistLoginSession(s.t1Ctx, t, r.Persister(), &flow.LoginSession{ID: t1f1.SessionID.String()})

			require.NoError(t, r.Persister().Connection(context.Background()).Create(t1f1))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(t2f1))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(t2f2))

			require.NoError(t, r.Persister().CreateConsentRequest(s.t1Ctx, nil, &flow.OAuth2ConsentRequest{
				ID:             t1f1.ID,
				LoginChallenge: sqlxx.NullString(t1f1.ID),
				Skip:           false,
				Verifier:       t1f1.ConsentVerifier.String(),
				CSRF:           "csrf"}))
			require.NoError(t, r.Persister().CreateConsentRequest(s.t2Ctx, nil, &flow.OAuth2ConsentRequest{
				ID:             t2f1.ID,
				LoginChallenge: sqlxx.NullString(t2f1.ID),
				Skip:           false,
				Verifier:       t2f1.ConsentVerifier.String(),
				CSRF:           "csrf"}))
			require.NoError(t, r.Persister().CreateConsentRequest(s.t2Ctx, nil, &flow.OAuth2ConsentRequest{
				ID:             t2f2.ID,
				LoginChallenge: sqlxx.NullString(t2f2.ID),
				Skip:           false,
				Verifier:       t2f2.ConsentVerifier.String(),
				CSRF:           "csrf"}))

			_, err := r.Persister().HandleConsentRequest(s.t1Ctx, nil, &flow.AcceptOAuth2ConsentRequest{
				ID:        x.Must(flowctx.Encode(s.t1Ctx, r.KeyCipher(), t1f1)),
				HandledAt: sqlxx.NullTime(time.Now()),
				Remember:  true})
			require.NoError(t, err)
			_, err = r.Persister().HandleConsentRequest(s.t2Ctx, nil, &flow.AcceptOAuth2ConsentRequest{
				ID:        x.Must(flowctx.Encode(s.t1Ctx, r.KeyCipher(), t2f1)),
				HandledAt: sqlxx.NullTime(time.Now()),
				Remember:  true})
			require.NoError(t, err)
			_, err = r.Persister().HandleConsentRequest(s.t2Ctx, nil, &flow.AcceptOAuth2ConsentRequest{
				ID:        x.Must(flowctx.Encode(s.t1Ctx, r.KeyCipher(), t2f2)),
				HandledAt: sqlxx.NullTime(time.Now()),
				Remember:  true})
			require.NoError(t, err)

			cs, err := r.Persister().ListUserAuthenticatedClientsWithBackChannelLogout(s.t1Ctx, "sub", t1f1.SessionID.String())
			require.NoError(t, err)
			require.Equal(t, 1, len(cs))

			cs, err = r.Persister().ListUserAuthenticatedClientsWithBackChannelLogout(s.t2Ctx, "sub", t1f1.SessionID.String())
			require.NoError(t, err)
			require.Equal(t, 2, len(cs))
		})
	}
}

func (s *PersisterTestSuite) TestListUserAuthenticatedClientsWithFrontChannelLogout() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			c1 := &client.Client{LegacyClientID: "client-1", FrontChannelLogoutURI: "not-null"}
			c2 := &client.Client{LegacyClientID: "client-2", FrontChannelLogoutURI: "not-null"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, c1))
			c1.ID = uuid.Nil
			require.NoError(t, r.Persister().CreateClient(s.t2Ctx, c1))
			require.NoError(t, r.Persister().CreateClient(s.t2Ctx, c2))

			t1f1 := newFlow(s.t1NID, c1.LegacyClientID, "sub", sqlxx.NullString(uuid.Must(uuid.NewV4()).String()))
			t1f1.ConsentChallengeID = "t1f1-consent-challenge"
			t1f1.LoginVerifier = "t1f1-login-verifier"
			t1f1.ConsentVerifier = "t1f1-consent-verifier"

			t2f1 := newFlow(s.t2NID, c1.LegacyClientID, "sub", t1f1.SessionID)
			t2f1.ConsentChallengeID = "t2f1-consent-challenge"
			t2f1.LoginVerifier = "t2f1-login-verifier"
			t2f1.ConsentVerifier = "t2f1-consent-verifier"

			t2f2 := newFlow(s.t2NID, c2.LegacyClientID, "sub", t1f1.SessionID)
			t2f2.ConsentChallengeID = "t2f2-consent-challenge"
			t2f2.LoginVerifier = "t2f2-login-verifier"
			t2f2.ConsentVerifier = "t2f2-consent-verifier"

			persistLoginSession(s.t1Ctx, t, r.Persister(), &flow.LoginSession{ID: t1f1.SessionID.String()})

			require.NoError(t, r.Persister().Connection(context.Background()).Create(t1f1))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(t2f1))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(t2f2))

			require.NoError(t, r.Persister().CreateConsentRequest(s.t1Ctx, nil, &flow.OAuth2ConsentRequest{
				ID:             t1f1.ID,
				LoginChallenge: sqlxx.NullString(t1f1.ID),
				Skip:           false,
				Verifier:       t1f1.ConsentVerifier.String(),
				CSRF:           "csrf"}))
			require.NoError(t, r.Persister().CreateConsentRequest(s.t2Ctx, nil, &flow.OAuth2ConsentRequest{
				ID:             t2f1.ID,
				LoginChallenge: sqlxx.NullString(t2f1.ID),
				Skip:           false,
				Verifier:       t2f1.ConsentVerifier.String(),
				CSRF:           "csrf"}))
			require.NoError(t, r.Persister().CreateConsentRequest(s.t2Ctx, nil, &flow.OAuth2ConsentRequest{
				ID:             t2f2.ID,
				LoginChallenge: sqlxx.NullString(t2f2.ID),
				Skip:           false,
				Verifier:       t2f2.ConsentVerifier.String(),
				CSRF:           "csrf"}))

			_, err := r.Persister().HandleConsentRequest(s.t1Ctx, nil, &flow.AcceptOAuth2ConsentRequest{
				ID:        x.Must(flowctx.Encode(s.t1Ctx, r.KeyCipher(), t1f1)),
				HandledAt: sqlxx.NullTime(time.Now()),
				Remember:  true})
			require.NoError(t, err)
			_, err = r.Persister().HandleConsentRequest(s.t2Ctx, nil, &flow.AcceptOAuth2ConsentRequest{
				ID:        x.Must(flowctx.Encode(s.t2Ctx, r.KeyCipher(), t2f1)),
				HandledAt: sqlxx.NullTime(time.Now()),
				Remember:  true})
			require.NoError(t, err)
			_, err = r.Persister().HandleConsentRequest(s.t2Ctx, nil, &flow.AcceptOAuth2ConsentRequest{
				ID:        x.Must(flowctx.Encode(s.t2Ctx, r.KeyCipher(), t2f2)),
				HandledAt: sqlxx.NullTime(time.Now()),
				Remember:  true})
			require.NoError(t, err)

			cs, err := r.Persister().ListUserAuthenticatedClientsWithFrontChannelLogout(s.t1Ctx, "sub", t1f1.SessionID.String())
			require.NoError(t, err)
			require.Equal(t, 1, len(cs))

			cs, err = r.Persister().ListUserAuthenticatedClientsWithFrontChannelLogout(s.t2Ctx, "sub", t1f1.SessionID.String())
			require.NoError(t, err)
			require.Equal(t, 2, len(cs))
		})
	}
}

func (s *PersisterTestSuite) TestMarkJWTUsedForTime() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			r.Persister().SetClientAssertionJWT(s.t1Ctx, "a", time.Now().Add(-24*time.Hour))
			r.Persister().SetClientAssertionJWT(s.t2Ctx, "a", time.Now().Add(-24*time.Hour))
			r.Persister().SetClientAssertionJWT(s.t2Ctx, "b", time.Now().Add(-24*time.Hour))

			require.NoError(t, r.Persister().MarkJWTUsedForTime(s.t2Ctx, "a", time.Now().Add(48*time.Hour)))

			store, ok := r.OAuth2Storage().(oauth2.AssertionJWTReader)
			if !ok {
				t.Fatal("type assertion failed")
			}

			_, err := store.GetClientAssertionJWT(s.t1Ctx, "a")
			require.NoError(t, err)
			_, err = store.GetClientAssertionJWT(s.t2Ctx, "a")
			require.NoError(t, err)
			_, err = store.GetClientAssertionJWT(s.t2Ctx, "b")
			require.Error(t, err)
		})
	}
}

func (s *PersisterTestSuite) TestQueryWithNetwork() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			r.Persister().CreateClient(s.t1Ctx, &client.Client{LegacyClientID: "client-1", FrontChannelLogoutURI: "not-null"})

			store, ok := r.Persister().(*persistencesql.Persister)
			if !ok {
				t.Fatal("type assertion failed")
			}

			var actual []client.Client
			store.QueryWithNetwork(s.t2Ctx).All(&actual)
			require.Equal(t, 0, len(actual))
			store.QueryWithNetwork(s.t1Ctx).All(&actual)
			require.Equal(t, 1, len(actual))
		})
	}
}

func (s *PersisterTestSuite) TestRejectLogoutRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			lr := newLogoutRequest()
			require.NoError(t, r.ConsentManager().CreateLogoutRequest(s.t1Ctx, lr))

			require.Error(t, r.ConsentManager().RejectLogoutRequest(s.t2Ctx, lr.ID))
			actual, err := r.ConsentManager().GetLogoutRequest(s.t1Ctx, lr.ID)
			require.NoError(t, err)
			require.Equal(t, lr, actual)

			require.NoError(t, r.ConsentManager().RejectLogoutRequest(s.t1Ctx, lr.ID))
			actual, err = r.ConsentManager().GetLogoutRequest(s.t1Ctx, lr.ID)
			require.Error(t, err)
			require.Equal(t, &flow.LogoutRequest{}, actual)
		})
	}
}

func (s *PersisterTestSuite) TestRevokeAccessToken() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: client.LegacyClientID}
			require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1Ctx, sig, fr))
			require.NoError(t, r.Persister().RevokeAccessToken(s.t2Ctx, fr.ID))

			actual := persistencesql.OAuth2RequestSQL{Table: "access"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, persistencesql.SignatureHash(sig)))
			require.Equal(t, s.t1NID, actual.NID)

			require.NoError(t, r.Persister().RevokeAccessToken(s.t1Ctx, fr.ID))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, persistencesql.SignatureHash(sig)))
		})
	}
}

func (s *PersisterTestSuite) TestRevokeRefreshToken() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}

			signature := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t1Ctx, signature, request))

			actual := persistencesql.OAuth2RequestSQL{Table: "refresh"}

			require.NoError(t, r.Persister().RevokeRefreshToken(s.t2Ctx, request.ID))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, signature))
			require.Equal(t, true, actual.Active)
			require.NoError(t, r.Persister().RevokeRefreshToken(s.t1Ctx, request.ID))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, signature))
			require.Equal(t, false, actual.Active)
		})
	}
}

func (s *PersisterTestSuite) TestRevokeRefreshTokenMaybeGracePeriod() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			client := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}

			signature := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t1Ctx, signature, request))

			actual := persistencesql.OAuth2RequestSQL{Table: "refresh"}

			store, ok := r.Persister().(*persistencesql.Persister)
			if !ok {
				t.Fatal("type assertion failed")
			}

			require.NoError(t, store.RevokeRefreshTokenMaybeGracePeriod(s.t2Ctx, request.ID, signature))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, signature))
			require.Equal(t, true, actual.Active)
			require.NoError(t, store.RevokeRefreshTokenMaybeGracePeriod(s.t1Ctx, request.ID, signature))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, signature))
			require.Equal(t, false, actual.Active)
		})
	}
}

func (s *PersisterTestSuite) TestRevokeSubjectClientConsentSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			client := &client.Client{LegacyClientID: "client-id"}
			f := newFlow(s.t1NID, client.LegacyClientID, "sub", sqlxx.NullString(sessionID))
			f.RequestedAt = time.Now().Add(-24 * time.Hour)
			persistLoginSession(s.t1Ctx, t, r.Persister(), &flow.LoginSession{ID: sessionID})
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(f))

			actual := flow.Flow{}

			require.Error(t, r.Persister().RevokeSubjectClientConsentSession(s.t2Ctx, "sub", client.LegacyClientID))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, f.ID))
			require.NoError(t, r.Persister().RevokeSubjectClientConsentSession(s.t1Ctx, "sub", client.LegacyClientID))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, f.ID))
		})
	}
}

func (s *PersisterTestSuite) TestSetClientAssertionJWT() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			jti := oauth2.NewBlacklistedJTI(uuid.Must(uuid.NewV4()).String(), time.Now().Add(24*time.Hour))
			require.NoError(t, r.Persister().SetClientAssertionJWT(s.t1Ctx, jti.JTI, jti.Expiry))

			actual := &oauth2.BlacklistedJTI{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(actual, jti.ID))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestSetClientAssertionJWTRaw() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			store, ok := r.Persister().(*persistencesql.Persister)
			if !ok {
				t.Fatal("type assertion failed")
			}

			jti := oauth2.NewBlacklistedJTI(uuid.Must(uuid.NewV4()).String(), time.Now().Add(24*time.Hour))
			require.NoError(t, store.SetClientAssertionJWTRaw(s.t1Ctx, jti))

			actual := &oauth2.BlacklistedJTI{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(actual, jti.ID))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestUpdateClient() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			t1c1 := &client.Client{LegacyClientID: "client-id", Name: "original", Secret: "original-secret"}
			t2c1 := &client.Client{LegacyClientID: "client-id", Name: "original", Secret: "original-secret"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, t1c1))
			require.NoError(t, r.Persister().CreateClient(s.t2Ctx, t2c1))
			expectedHash := t1c1.Secret

			u1 := *t1c1
			u1.Name = "updated"
			u1.Secret = ""
			require.NoError(t, r.Persister().UpdateClient(s.t2Ctx, &u1))

			actual := &client.Client{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(actual, t1c1.ID))
			require.Equal(t, "original", actual.Name)
			require.Equal(t, expectedHash, actual.Secret)

			u2 := *t1c1
			u2.Name = "updated"
			u2.Secret = ""
			require.NoError(t, r.Persister().UpdateClient(s.t1Ctx, &u2))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(actual, t1c1.ID))
			require.Equal(t, "updated", actual.Name)
			require.Equal(t, expectedHash, actual.Secret)

			u3 := *t1c1
			u3.Name = "updated"
			u3.Secret = "updated-secret"
			require.NoError(t, r.Persister().UpdateClient(s.t1Ctx, &u3))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(actual, t1c1.ID))
			require.Equal(t, "updated", actual.Name)
			require.NotEqual(t, expectedHash, actual.Secret)
		})
	}
}

func (s *PersisterTestSuite) TestUpdateKey() {
	t := s.T()
	for k, r := range s.registries {
		t.Run("dialect="+k, func(*testing.T) {
			k1 := newKey("test-ks", "test")
			ks := "key-set"
			require.NoError(t, r.Persister().AddKey(s.t1Ctx, ks, &k1))
			actual, err := r.Persister().GetKey(s.t1Ctx, ks, k1.KeyID)
			require.NoError(t, err)
			assertx.EqualAsJSON(t, &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{k1}}, actual)

			k2 := newKey("test-ks", "test")
			r.Persister().UpdateKey(s.t2Ctx, ks, &k2)
			actual, err = r.Persister().GetKey(s.t1Ctx, ks, k1.KeyID)
			require.NoError(t, err)
			assertx.EqualAsJSON(t, &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{k1}}, actual)

			r.Persister().UpdateKey(s.t1Ctx, ks, &k2)
			actual, err = r.Persister().GetKey(s.t1Ctx, ks, k2.KeyID)
			require.NoError(t, err)
			require.NotEqual(t, &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{k1}}, actual)
		})
	}
}

func (s *PersisterTestSuite) TestUpdateKeySet() {
	t := s.T()
	for k, r := range s.registries {
		t.Run("dialect="+k, func(*testing.T) {
			ks := "key-set"
			ks1 := newKeySet(ks, "test")
			require.NoError(t, r.Persister().AddKeySet(s.t1Ctx, ks, ks1))
			actual, err := r.Persister().GetKeySet(s.t1Ctx, ks)
			require.NoError(t, err)
			requireKeySetEqual(t, ks1, actual)

			ks2 := newKeySet(ks, "test")
			r.Persister().UpdateKeySet(s.t2Ctx, ks, ks2)
			actual, err = r.Persister().GetKeySet(s.t1Ctx, ks)
			require.NoError(t, err)
			requireKeySetEqual(t, ks1, actual)

			r.Persister().UpdateKeySet(s.t1Ctx, ks, ks2)
			actual, err = r.Persister().GetKeySet(s.t1Ctx, ks)
			require.NoError(t, err)
			requireKeySetEqual(t, ks2, actual)
		})
	}
}

func (s *PersisterTestSuite) TestUpdateWithNetwork() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			t1c1 := &client.Client{LegacyClientID: "client-id", Name: "original", Secret: "original-secret"}
			t2c1 := &client.Client{LegacyClientID: "client-id", Name: "original", Secret: "original-secret", Owner: "erase-me"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, t1c1))
			require.NoError(t, r.Persister().CreateClient(s.t2Ctx, t2c1))

			store, ok := r.Persister().(*persistencesql.Persister)
			if !ok {
				t.Fatal("type assertion failed")
			}

			count, err := store.UpdateWithNetwork(s.t1Ctx, &client.Client{ID: t1c1.ID, LegacyClientID: "client-id", Name: "updated", Secret: "original-secret"})
			require.NoError(t, err)
			require.Equal(t, int64(1), count)
			actualt1, err := store.GetConcreteClient(s.t1Ctx, "client-id")
			require.NoError(t, err)
			require.Equal(t, "updated", actualt1.Name)
			require.Equal(t, "", actualt1.Owner)
			actualt2, err := store.GetConcreteClient(s.t2Ctx, "client-id")
			require.NoError(t, err)
			require.Equal(t, "original", actualt2.Name)
			require.Equal(t, "erase-me", actualt2.Owner)
		})
	}
}

func (s *PersisterTestSuite) TestVerifyAndInvalidateConsentRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			sub := uuid.Must(uuid.NewV4()).String()
			sessionID := uuid.Must(uuid.NewV4()).String()
			persistLoginSession(s.t1Ctx, t, r.Persister(), &flow.LoginSession{ID: sessionID})
			client := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			f := newFlow(s.t1NID, client.LegacyClientID, sub, sqlxx.NullString(sessionID))
			f.ConsentSkip = false
			f.GrantedScope = sqlxx.StringSliceJSONFormat{}
			f.ConsentRemember = false
			crf := 86400
			f.ConsentRememberFor = &crf
			f.ConsentError = &flow.RequestDeniedError{}
			f.SessionAccessToken = map[string]interface{}{}
			f.SessionIDToken = map[string]interface{}{}
			f.ConsentWasHandled = false
			f.State = flow.FlowStateConsentUnused
			f.ConsentVerifier = sqlxx.NullString(x.Must(flowctx.Encode(s.t1Ctx, r.KeyCipher(), f)))

			_, err := r.ConsentManager().VerifyAndInvalidateConsentRequest(s.t2Ctx, nil, f.ConsentVerifier.String())
			require.Error(t, err)
			require.Equal(t, flow.FlowStateConsentUnused, f.State)
			require.Equal(t, false, f.ConsentWasHandled)
			_, err = r.ConsentManager().VerifyAndInvalidateConsentRequest(s.t1Ctx, nil, f.ConsentVerifier.String())
			require.NoError(t, err)
			require.Equal(t, flow.FlowStateConsentUsed, f.State)
			require.Equal(t, true, f.ConsentWasHandled)
		})
	}
}

func (s *PersisterTestSuite) TestVerifyAndInvalidateLoginRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			sub := uuid.Must(uuid.NewV4()).String()
			sessionID := uuid.Must(uuid.NewV4()).String()
			persistLoginSession(s.t1Ctx, t, r.Persister(), &flow.LoginSession{ID: sessionID})
			client := &client.Client{LegacyClientID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1Ctx, client))
			f := newFlow(s.t1NID, client.LegacyClientID, sub, sqlxx.NullString(sessionID))
			f.State = flow.FlowStateLoginUnused
			f.LoginVerifier = x.Must(flowctx.Encode(s.t1Ctx, r.KeyCipher(), f))

			_, err := r.ConsentManager().VerifyAndInvalidateLoginRequest(s.t2Ctx, nil, f.LoginVerifier)
			require.Error(t, err)
			require.Equal(t, flow.FlowStateLoginUnused, f.State)
			require.Equal(t, false, f.LoginWasUsed)

			_, err = r.ConsentManager().VerifyAndInvalidateLoginRequest(s.t1Ctx, nil, f.LoginVerifier)
			require.NoError(t, err)
			require.Equal(t, flow.FlowStateLoginUsed, f.State)
			require.Equal(t, true, f.LoginWasUsed)
		})
	}
}

func (s *PersisterTestSuite) TestVerifyAndInvalidateLogoutRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			lr := newLogoutRequest()
			lr.Verifier = uuid.Must(uuid.NewV4()).String()
			lr.Accepted = true
			lr.Rejected = false
			require.NoError(t, r.ConsentManager().CreateLogoutRequest(s.t1Ctx, lr))

			expected, err := r.ConsentManager().GetLogoutRequest(s.t1Ctx, lr.ID)
			require.NoError(t, err)

			lrInvalidated, err := r.ConsentManager().VerifyAndInvalidateLogoutRequest(s.t2Ctx, lr.Verifier)
			require.Error(t, err)
			require.Equal(t, &flow.LogoutRequest{}, lrInvalidated)
			actual := &flow.LogoutRequest{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(actual, lr.ID))
			require.Equal(t, expected, actual)

			lrInvalidated, err = r.ConsentManager().VerifyAndInvalidateLogoutRequest(s.t1Ctx, lr.Verifier)
			require.NoError(t, err)
			require.NoError(t, r.Persister().Connection(context.Background()).Find(actual, lr.ID))
			require.Equal(t, lrInvalidated, actual)
			require.Equal(t, true, actual.WasHandled)
		})
	}
}

func (s *PersisterTestSuite) TestWithFallbackNetworkID() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			r.WithContextualizer(&contextx.Default{})
			store, ok := r.Persister().(*persistencesql.Persister)
			if !ok {
				t.Fatal("type assertion failed")
			}
			original := store.NetworkID(context.Background())
			expected := uuid.Must(uuid.NewV4())
			store, ok = store.WithFallbackNetworkID(expected).(*persistencesql.Persister)
			if !ok {
				t.Fatal("type assertion failed")
			}

			require.NotEqual(t, original, expected)
			require.Equal(t, expected, store.NetworkID(context.Background()))
		})
	}
}

func TestPersisterTestSuite(t *testing.T) {
	suite.Run(t, new(PersisterTestSuite))
}

func newClient() *client.Client {
	return &client.Client{
		ID: uuid.Must(uuid.NewV4()),
	}
}

func newFlow(nid uuid.UUID, clientID string, subject string, sessionID sqlxx.NullString) *flow.Flow {
	return &flow.Flow{
		NID:                nid,
		ID:                 uuid.Must(uuid.NewV4()).String(),
		ClientID:           clientID,
		Subject:            subject,
		ConsentError:       &flow.RequestDeniedError{},
		State:              flow.FlowStateConsentUnused,
		LoginError:         &flow.RequestDeniedError{},
		Context:            sqlxx.JSONRawMessage{},
		AMR:                sqlxx.StringSliceJSONFormat{},
		ConsentChallengeID: sqlxx.NullString("not-null"),
		ConsentVerifier:    sqlxx.NullString("not-null"),
		ConsentCSRF:        sqlxx.NullString("not-null"),
		SessionID:          sessionID,
		RequestedAt:        time.Now(),
	}
}

func newGrant(keySet string, keyID string) trust.Grant {
	return trust.Grant{
		ID:        uuid.Must(uuid.NewV4()).String(),
		ExpiresAt: time.Now().Add(time.Hour),
		PublicKey: trust.PublicKey{
			Set:   keySet,
			KeyID: keyID,
		},
	}
}

func newLogoutRequest() *flow.LogoutRequest {
	return &flow.LogoutRequest{
		ID: uuid.Must(uuid.NewV4()).String(),
	}
}

func newKey(ksID string, use string) jose.JSONWebKey {
	ks, err := jwk.GenerateJWK(context.Background(), jose.RS256, ksID, use)
	if err != nil {
		panic(err)
	}
	return ks.Keys[0]
}

func newKeySet(id string, use string) *jose.JSONWebKeySet {
	return x.Must(jwk.GenerateJWK(context.Background(), jose.RS256, id, use))
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
	require.NoError(t, p.CreateLoginSession(ctx, session))
	require.NoError(t, p.Connection(ctx).Create(session))
}
