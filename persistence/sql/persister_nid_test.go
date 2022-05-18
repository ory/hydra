package sql_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/instana/testify/require"
	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/flow"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/oauth2/trust"
	persistencesql "github.com/ory/hydra/persistence/sql"
	"github.com/ory/hydra/x"
	"github.com/ory/hydra/x/contextx"
	"github.com/ory/x/dbal"
	"github.com/ory/x/networkx"
	"github.com/ory/x/sqlxx"
	"github.com/stretchr/testify/suite"
	"gopkg.in/square/go-jose.v2"
)

type PersisterTestSuite struct {
	suite.Suite
	registries map[string]driver.Registry
	clean      func(*testing.T)
	t1         context.Context
	t2         context.Context
	t1NID      uuid.UUID
	t2NID      uuid.UUID
}

var _ PersisterTestSuite = PersisterTestSuite{}

func (s *PersisterTestSuite) SetupSuite() {
	s.registries = map[string]driver.Registry{
		"memory": internal.NewRegistrySQLFromURL(s.T(), dbal.SQLiteSharedInMemory, true, &contextx.DefaultContextualizer{}),
	}

	if !testing.Short() {
		s.registries["postgres"], s.registries["mysql"], s.registries["cockroach"], s.clean = internal.ConnectDatabases(s.T(), true, &contextx.DefaultContextualizer{})
	}

	s.t1NID, s.t2NID = uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4())
	s.t1 = contextx.SetNIDContext(context.Background(), s.t1NID)
	s.t2 = contextx.SetNIDContext(context.Background(), s.t2NID)

	for _, r := range s.registries {
		require.NoError(s.T(), r.Persister().Connection(context.Background()).Create(&networkx.Network{ID: s.t1NID}))
		require.NoError(s.T(), r.Persister().Connection(context.Background()).Create(&networkx.Network{ID: s.t2NID}))
		r.WithContextualizer(&contextx.TestContextualizer{})
	}
}

func (s *PersisterTestSuite) TearDownTest() {
	for _, r := range s.registries {
		x.DeleteHydraRows(s.T(), r.Persister().Connection(context.Background()))
	}
}

func (s *PersisterTestSuite) TestAcceptLogoutRequest() {
	t := s.T()
	lr := newLogoutRequest()

	for k, r := range s.registries {
		t.Run("dialect="+k, func(*testing.T) {
			require.NoError(t, r.ConsentManager().CreateLogoutRequest(s.t1, lr))
			lrAccepted, err := r.ConsentManager().AcceptLogoutRequest(s.t2, lr.ID)
			require.Error(t, err)
			require.Equal(t, &consent.LogoutRequest{}, lrAccepted)

			actual, err := r.ConsentManager().GetLogoutRequest(s.t1, lr.ID)
			require.NoError(t, err)
			require.Equal(t, lr, actual)

			lrActual, err := r.ConsentManager().GetLogoutRequest(s.t1, lr.ID)
			require.NoError(t, err)
			require.Equal(t, false, lrActual.Accepted)
		})
	}
}

func (s *PersisterTestSuite) TestAddKeyGetKeyDeleteKey() {
	t := s.T()
	key := newKey("test-ks", "test")
	for k, r := range s.registries {
		t.Run("dialect="+k, func(*testing.T) {
			ks := "key-set"
			require.NoError(t, r.Persister().AddKey(s.t1, ks, &key))
			actual, err := r.Persister().GetKey(s.t2, ks, key.KeyID)
			require.Error(t, err)
			require.Equal(t, (*jose.JSONWebKeySet)(nil), actual)
			actual, err = r.Persister().GetKey(s.t1, ks, key.KeyID)
			require.NoError(t, err)
			require.Equal(t, &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{key}}, actual)

			r.Persister().DeleteKey(s.t2, ks, key.KeyID)
			_, err = r.Persister().GetKey(s.t1, ks, key.KeyID)
			require.NoError(t, err)
			r.Persister().DeleteKey(s.t1, ks, key.KeyID)
			_, err = r.Persister().GetKey(s.t1, ks, key.KeyID)
			require.Error(t, err)
		})
	}
}

func (s *PersisterTestSuite) TestAddKeySetGetKeySetDeleteKeySet() {
	t := s.T()
	ks := newKeySet("test-ks", "test")
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			ksID := "key-set"
			r.Persister().AddKeySet(s.t1, ksID, ks)
			actual, err := r.Persister().GetKeySet(s.t2, ksID)
			require.Error(t, err)
			require.Equal(t, (*jose.JSONWebKeySet)(nil), actual)
			actual, err = r.Persister().GetKeySet(s.t1, ksID)
			require.NoError(t, err)
			if actual.Keys[0].KeyID == ks.Keys[1].KeyID {
				actual.Keys[0], actual.Keys[1] = actual.Keys[1], actual.Keys[0]
			}
			require.Equal(t, ks, actual)

			r.Persister().DeleteKeySet(s.t2, ksID)
			_, err = r.Persister().GetKeySet(s.t1, ksID)
			require.NoError(t, err)
			r.Persister().DeleteKeySet(s.t1, ksID)
			_, err = r.Persister().GetKeySet(s.t1, ksID)
			require.Error(t, err)
		})
	}
}

func (s *PersisterTestSuite) TestConfirmLoginSession() {
	t := s.T()
	ls := newLoginSession()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			require.NoError(t, r.Persister().CreateLoginSession(s.t1, ls))
			expected := &consent.LoginSession{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(expected, ls.ID))

			require.NoError(t, r.Persister().ConfirmLoginSession(s.t2, expected.ID, time.Now(), expected.Subject, !expected.Remember))
			actual := &consent.LoginSession{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(actual, ls.ID))
			require.Equal(t, expected, actual)

			require.NoError(t, r.Persister().ConfirmLoginSession(s.t1, expected.ID, time.Now(), expected.Subject, !expected.Remember))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(actual, ls.ID))
			require.NotEqual(t, expected, actual)
		})
	}
}

func (s *PersisterTestSuite) TestCreateSession() {
	t := s.T()
	ls := newLoginSession()
	for k, r := range s.registries {
		t.Run(k, func(t *testing.T) {
			require.NoError(t, r.Persister().CreateLoginSession(s.t1, ls))
			actual := &consent.LoginSession{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(actual, ls.ID))
			require.Equal(t, s.t1NID, actual.NID)
			ls.NID = actual.NID
			require.Equal(t, ls, actual)
		})
	}
}

func (s *PersisterTestSuite) TestCountClients() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			count, err := r.Persister().CountClients(s.t1)
			require.NoError(t, err)
			require.Equal(t, 0, count)

			count, err = r.Persister().CountClients(s.t2)
			require.NoError(t, err)
			require.Equal(t, 0, count)

			require.NoError(t, r.Persister().CreateClient(s.t1, newClient()))

			count, err = r.Persister().CountClients(s.t1)
			require.NoError(t, err)
			require.Equal(t, 1, count)

			count, err = r.Persister().CountClients(s.t2)
			require.NoError(t, err)
			require.Equal(t, 0, count)
		})
	}
}

func (s *PersisterTestSuite) TestCountGrants() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			count, err := r.Persister().CountGrants(s.t1)
			require.NoError(t, err)
			require.Equal(t, 0, count)

			count, err = r.Persister().CountGrants(s.t2)
			require.NoError(t, err)
			require.Equal(t, 0, count)

			keySet := uuid.Must(uuid.NewV4()).String()
			publicKey := newKey(keySet, "use")
			grant := newGrant(keySet, publicKey.KeyID)
			require.NoError(t, r.Persister().AddKey(s.t1, keySet, &publicKey))
			require.NoError(t, r.Persister().CreateGrant(s.t1, grant, publicKey))

			count, err = r.Persister().CountGrants(s.t1)
			require.NoError(t, err)
			require.Equal(t, 1, count)

			count, err = r.Persister().CountGrants(s.t2)
			require.NoError(t, err)
			require.Equal(t, 0, count)
		})
	}
}

func (s *PersisterTestSuite) TestCountSubjectsGrantedConsentRequests() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			sub := uuid.Must(uuid.NewV4()).String()
			count, err := r.Persister().CountSubjectsGrantedConsentRequests(s.t1, sub)
			require.NoError(t, err)
			require.Equal(t, 0, count)

			count, err = r.Persister().CountSubjectsGrantedConsentRequests(s.t2, sub)
			require.NoError(t, err)
			require.Equal(t, 0, count)

			sessionID := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateLoginSession(s.t1, &consent.LoginSession{ID: sessionID}))
			client := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			f := newFlow(s.t1NID, client.OutfacingID, sub, sqlxx.NullString(sessionID))
			f.ConsentSkip = false
			f.ConsentError = &consent.RequestDeniedError{}
			f.State = flow.FlowStateConsentUnused
			require.NoError(t, r.Persister().Connection(context.Background()).Create(f))

			count, err = r.Persister().CountSubjectsGrantedConsentRequests(s.t1, sub)
			require.NoError(t, err)
			require.Equal(t, 1, count)

			count, err = r.Persister().CountSubjectsGrantedConsentRequests(s.t2, sub)
			require.NoError(t, err)
			require.Equal(t, 0, count)
		})
	}
}

func (s *PersisterTestSuite) TestCreateAccessTokenSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.Persister().CreateClient(s.t2, client))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: client.OutfacingID}
			require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1, sig, fr))
			actual := persistencesql.OAuth2RequestSQL{Table: "access"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, sig))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateAuthorizeCodeSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.Persister().CreateClient(s.t2, client))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: client.OutfacingID}
			require.NoError(t, r.Persister().CreateAuthorizeCodeSession(s.t1, sig, fr))
			actual := persistencesql.OAuth2RequestSQL{Table: "code"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, sig))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateClient() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			expected := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, expected))
			actual := client.Client{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, expected.ID))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateConsentRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			client := &client.Client{OutfacingID: "client-id"}
			f := newFlow(s.t1NID, client.OutfacingID, "sub", sqlxx.NullString(sessionID))
			require.NoError(t, r.Persister().CreateLoginSession(s.t1, &consent.LoginSession{ID: sessionID}))
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(f))

			req := &consent.ConsentRequest{
				ID:             "consent-request-id",
				LoginChallenge: sqlxx.NullString(f.ID),
				Skip:           false,
				Verifier:       "verifier",
				CSRF:           "csrf",
			}
			require.NoError(t, r.Persister().CreateConsentRequest(s.t1, req))

			actual := flow.Flow{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, f.ID))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateForcedObfuscatedLoginSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			session := &consent.ForcedObfuscatedLoginSession{ClientID: client.OutfacingID}
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.Persister().CreateForcedObfuscatedLoginSession(s.t1, session))
			actual, err := r.Persister().GetForcedObfuscatedLoginSession(s.t1, client.OutfacingID, "")
			require.NoError(t, err)
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateGrant() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()).String(),
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().AddKeySet(s.t1, "ks-id", ks))
			require.NoError(t, r.Persister().CreateGrant(s.t1, grant, ks.Keys[0]))
			actual := trust.SQLData{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, grant.ID))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateLoginRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			lr := consent.LoginRequest{ID: "lr-id", ClientID: client.OutfacingID, RequestedAt: time.Now()}

			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.ConsentManager().CreateLoginRequest(s.t1, &lr))
			f := flow.Flow{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&f, lr.ID))
			require.Equal(t, s.t1NID, f.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateLoginSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			ls := consent.LoginSession{ID: uuid.Must(uuid.NewV4()).String(), Remember: true}
			require.NoError(t, r.Persister().CreateLoginSession(s.t1, &ls))
			actual, err := r.Persister().GetRememberedLoginSession(s.t1, ls.ID)
			require.NoError(t, err)
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateLogoutRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			lr := consent.LogoutRequest{
				// TODO there is not FK for SessionID so we don't need it here; TODO make sure the missing FK is intentional
				ID:       uuid.Must(uuid.NewV4()).String(),
				ClientID: sql.NullString{Valid: true, String: client.OutfacingID},
			}

			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.Persister().CreateLogoutRequest(s.t1, &lr))
			actual, err := r.Persister().GetLogoutRequest(s.t1, lr.ID)
			require.NoError(t, err)
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateOpenIDConnectSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, client))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}

			requester := fosite.AccessRequest{
				GrantTypes:       nil,
				HandledGrantType: nil,
				Request:          *request,
			}
			authorizeCode := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateOpenIDConnectSession(s.t1, authorizeCode, &requester))

			actual := persistencesql.OAuth2RequestSQL{Table: "oidc"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreatePKCERequestSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, client))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}

			requester := fosite.AccessRequest{
				GrantTypes:       nil,
				HandledGrantType: nil,
				Request:          *request,
			}

			authorizeCode := uuid.Must(uuid.NewV4()).String()
			r.Persister().CreatePKCERequestSession(s.t1, authorizeCode, &requester)

			actual := persistencesql.OAuth2RequestSQL{Table: "pkce"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) TestCreateRefreshTokenSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, client))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}

			requester := fosite.AccessRequest{
				GrantTypes:       nil,
				HandledGrantType: nil,
				Request:          *request,
			}

			authorizeCode := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t1, authorizeCode, &requester))

			actual := persistencesql.OAuth2RequestSQL{Table: "refresh"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.Equal(t, s.t1NID, actual.NID)
		})
	}
}

func (s *PersisterTestSuite) DeleteAccessTokenSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: client.OutfacingID}
			require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1, sig, fr))
			require.NoError(t, r.Persister().DeleteAccessTokenSession(s.t2, sig))

			actual := persistencesql.OAuth2RequestSQL{Table: "access"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, sig))
			require.Equal(t, s.t1NID, actual.NID)

			require.NoError(t, r.Persister().DeleteAccessTokenSession(s.t1, sig))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, sig))
		})
	}
}

func (s *PersisterTestSuite) TestDeleteAccessTokens() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: client.OutfacingID}
			require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1, sig, fr))
			require.NoError(t, r.Persister().DeleteAccessTokens(s.t2, client.OutfacingID))

			actual := persistencesql.OAuth2RequestSQL{Table: "access"}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, sig))
			require.Equal(t, s.t1NID, actual.NID)

			require.NoError(t, r.Persister().DeleteAccessTokens(s.t1, client.OutfacingID))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, sig))
		})
	}
}

func (s *PersisterTestSuite) TestDeleteClient() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			c := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, c))
			actual := client.Client{}
			require.Error(t, r.Persister().DeleteClient(s.t2, c.OutfacingID))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, c.ID))
			require.NoError(t, r.Persister().DeleteClient(s.t1, c.OutfacingID))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, c.ID))
		})
	}
}

func (s *PersisterTestSuite) TestDeleteGrant() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()).String(),
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().AddKeySet(s.t1, "ks-id", ks))
			require.NoError(t, r.Persister().CreateGrant(s.t1, grant, ks.Keys[0]))

			actual := trust.SQLData{}
			require.Error(t, r.Persister().DeleteGrant(s.t2, grant.ID))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, grant.ID))
			require.NoError(t, r.Persister().DeleteGrant(s.t1, grant.ID))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, grant.ID))
		})
	}
}

func (s *PersisterTestSuite) TestDeleteLoginSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			ls := consent.LoginSession{ID: uuid.Must(uuid.NewV4()).String(), Remember: true}
			require.NoError(t, r.Persister().CreateLoginSession(s.t1, &ls))

			require.Error(t, r.Persister().DeleteLoginSession(s.t2, ls.ID))
			_, err := r.Persister().GetRememberedLoginSession(s.t1, ls.ID)
			require.NoError(t, err)

			require.NoError(t, r.Persister().DeleteLoginSession(s.t1, ls.ID))
			_, err = r.Persister().GetRememberedLoginSession(s.t1, ls.ID)
			require.Error(t, err)
		})
	}
}

func (s *PersisterTestSuite) TestDeleteOpenIDConnectSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, client))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}

			requester := fosite.AccessRequest{
				GrantTypes:       nil,
				HandledGrantType: nil,
				Request:          *request,
			}
			authorizeCode := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateOpenIDConnectSession(s.t1, authorizeCode, &requester))

			actual := persistencesql.OAuth2RequestSQL{Table: "oidc"}

			require.NoError(t, r.Persister().DeleteOpenIDConnectSession(s.t2, authorizeCode))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.NoError(t, r.Persister().DeleteOpenIDConnectSession(s.t1, authorizeCode))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
		})
	}
}

func (s *PersisterTestSuite) TestDeletePKCERequestSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, client))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}

			requester := fosite.AccessRequest{
				GrantTypes:       nil,
				HandledGrantType: nil,
				Request:          *request,
			}

			authorizeCode := uuid.Must(uuid.NewV4()).String()
			r.Persister().CreatePKCERequestSession(s.t1, authorizeCode, &requester)

			actual := persistencesql.OAuth2RequestSQL{Table: "pkce"}

			require.NoError(t, r.Persister().DeletePKCERequestSession(s.t2, authorizeCode))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
			require.NoError(t, r.Persister().DeletePKCERequestSession(s.t1, authorizeCode))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, authorizeCode))
		})
	}
}

func (s *PersisterTestSuite) TestDeleteRefreshTokenSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, client))

			request := fosite.NewRequest()
			request.Client = &fosite.DefaultClient{ID: "client-id"}

			requester := fosite.AccessRequest{
				GrantTypes:       nil,
				HandledGrantType: nil,
				Request:          *request,
			}

			signature := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t1, signature, &requester))

			actual := persistencesql.OAuth2RequestSQL{Table: "refresh"}

			require.NoError(t, r.Persister().DeleteRefreshTokenSession(s.t2, signature))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, signature))
			require.NoError(t, r.Persister().DeleteRefreshTokenSession(s.t1, signature))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, signature))
		})
	}
}

func (s *PersisterTestSuite) TestFindGrantedAndRememberedConsentRequests() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			client := &client.Client{OutfacingID: "client-id"}
			f := newFlow(s.t1NID, client.OutfacingID, "sub", sqlxx.NullString(sessionID))
			require.NoError(t, r.Persister().CreateLoginSession(s.t1, &consent.LoginSession{ID: sessionID}))
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(f))

			req := &consent.ConsentRequest{
				ID:             "consent-request-id",
				LoginChallenge: sqlxx.NullString(f.ID),
				Skip:           false,
				Verifier:       "verifier",
				CSRF:           "csrf",
			}

			hcr := &consent.HandledConsentRequest{
				ID:        req.ID,
				HandledAt: sqlxx.NullTime(time.Now()),
				Remember:  true,
			}
			require.NoError(t, r.Persister().CreateConsentRequest(s.t1, req))
			_, err := r.Persister().HandleConsentRequest(s.t1, hcr)
			require.NoError(t, err)

			actual, err := r.Persister().FindGrantedAndRememberedConsentRequests(s.t2, client.OutfacingID, f.Subject)
			require.Error(t, err)
			require.Equal(t, 0, len(actual))

			actual, err = r.Persister().FindGrantedAndRememberedConsentRequests(s.t1, client.OutfacingID, f.Subject)
			require.NoError(t, err)
			require.Equal(t, 1, len(actual))
		})
	}
}

func (s *PersisterTestSuite) TestFindSubjectsGrantedConsentRequests() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			client := &client.Client{OutfacingID: "client-id"}
			f := newFlow(s.t1NID, client.OutfacingID, "sub", sqlxx.NullString(sessionID))
			require.NoError(t, r.Persister().CreateLoginSession(s.t1, &consent.LoginSession{ID: sessionID}))
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(f))

			req := &consent.ConsentRequest{
				ID:             "consent-request-id",
				LoginChallenge: sqlxx.NullString(f.ID),
				Skip:           false,
				Verifier:       "verifier",
				CSRF:           "csrf",
			}

			hcr := &consent.HandledConsentRequest{
				ID:        req.ID,
				HandledAt: sqlxx.NullTime(time.Now()),
				Remember:  true,
			}
			require.NoError(t, r.Persister().CreateConsentRequest(s.t1, req))
			_, err := r.Persister().HandleConsentRequest(s.t1, hcr)
			require.NoError(t, err)

			actual, err := r.Persister().FindSubjectsGrantedConsentRequests(s.t2, f.Subject, 100, 0)
			require.Error(t, err)
			require.Equal(t, 0, len(actual))

			actual, err = r.Persister().FindSubjectsGrantedConsentRequests(s.t1, f.Subject, 100, 0)
			require.NoError(t, err)
			require.Equal(t, 1, len(actual))
		})
	}
}

func (s *PersisterTestSuite) TestFlushInactiveAccessTokens() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.RequestedAt = time.Now().UTC().Add(-24 * time.Hour)
			fr.Client = &fosite.DefaultClient{ID: client.OutfacingID}
			require.NoError(t, r.Persister().CreateAccessTokenSession(s.t1, sig, fr))

			actual := persistencesql.OAuth2RequestSQL{Table: "access"}

			require.NoError(t, r.Persister().FlushInactiveAccessTokens(s.t2, time.Now().Add(time.Hour), 100, 100))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, sig))
			require.NoError(t, r.Persister().FlushInactiveAccessTokens(s.t1, time.Now().Add(time.Hour), 100, 100))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, sig))
		})
	}
}

func (s *PersisterTestSuite) TestFlushInactiveGrants() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()).String(),
				ExpiresAt: time.Now().Add(-24 * time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().AddKeySet(s.t1, "ks-id", ks))
			require.NoError(t, r.Persister().CreateGrant(s.t1, grant, ks.Keys[0]))

			actual := trust.SQLData{}
			require.NoError(t, r.Persister().FlushInactiveGrants(s.t2, time.Now(), 100, 100))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, grant.ID))
			require.NoError(t, r.Persister().FlushInactiveGrants(s.t1, time.Now(), 100, 100))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, grant.ID))
		})
	}
}

func (s *PersisterTestSuite) TestFlushInactiveLoginConsentRequests() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			client := &client.Client{OutfacingID: "client-id"}
			f := newFlow(s.t1NID, client.OutfacingID, "sub", sqlxx.NullString(sessionID))
			f.RequestedAt = time.Now().Add(-24 * time.Hour)
			require.NoError(t, r.Persister().CreateLoginSession(s.t1, &consent.LoginSession{ID: sessionID}))
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(f))

			actual := flow.Flow{}

			require.NoError(t, r.Persister().FlushInactiveLoginConsentRequests(s.t2, time.Now(), 100, 100))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, f.ID))
			require.NoError(t, r.Persister().FlushInactiveLoginConsentRequests(s.t1, time.Now(), 100, 100))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, f.ID))
		})
	}
}

func (s *PersisterTestSuite) TestFlushInactiveRefreshTokens() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			request := fosite.NewRequest()
			request.RequestedAt = time.Now().Add(-240 * 365 * time.Hour)
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			signature := uuid.Must(uuid.NewV4()).String()

			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t1, signature, request))

			actual := persistencesql.OAuth2RequestSQL{Table: "refresh"}

			require.NoError(t, r.Persister().FlushInactiveRefreshTokens(s.t2, time.Now(), 100, 100))
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&actual, signature))
			require.NoError(t, r.Persister().FlushInactiveRefreshTokens(s.t1, time.Now(), 100, 100))
			require.Error(t, r.Persister().Connection(context.Background()).Find(&actual, signature))
		})
	}
}

func (s *PersisterTestSuite) TestGetAccessTokenSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: client.OutfacingID}
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
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			sig := uuid.Must(uuid.NewV4()).String()
			fr := fosite.NewRequest()
			fr.Client = &fosite.DefaultClient{ID: client.OutfacingID}
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
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			expected := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, expected))

			actual, err := r.Persister().GetClient(s.t2, expected.OutfacingID)
			require.Error(t, err)
			require.Equal(t, "", actual.GetID())
			actual, err = r.Persister().GetClient(s.t1, expected.OutfacingID)
			require.NoError(t, err)
			require.Equal(t, expected.OutfacingID, actual.GetID())
		})
	}
}

func (s *PersisterTestSuite) TestGetClientAssertionJWT() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			store, ok := r.OAuth2Storage().(oauth2.AssertionJWTReader)
			if !ok {
				t.Fatal("type assertion failed")
			}
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
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			c := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, c))

			actual, err := r.Persister().GetClients(s.t2, client.Filter{Offset: 0, Limit: 100})
			require.NoError(t, err)
			require.Equal(t, 0, len(actual))
			actual, err = r.Persister().GetClients(s.t1, client.Filter{Offset: 0, Limit: 100})
			require.NoError(t, err)
			require.Equal(t, 1, len(actual))
		})
	}
}

func (s *PersisterTestSuite) TestGetConcreteClient() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			expected := &client.Client{OutfacingID: "client-id"}
			require.NoError(t, r.Persister().CreateClient(s.t1, expected))

			actual, err := r.Persister().GetConcreteClient(s.t2, expected.OutfacingID)
			require.Error(t, err)
			require.Equal(t, "", actual.GetID())
			actual, err = r.Persister().GetConcreteClient(s.t1, expected.OutfacingID)
			require.NoError(t, err)
			require.Equal(t, expected.OutfacingID, actual.GetID())
		})
	}
}

func (s *PersisterTestSuite) TestGetConcreteGrant() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()).String(),
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().AddKeySet(s.t1, "ks-id", ks))
			require.NoError(t, r.Persister().CreateGrant(s.t1, grant, ks.Keys[0]))

			actual, err := r.Persister().GetConcreteGrant(s.t2, grant.ID)
			require.Error(t, err)
			require.Equal(t, trust.Grant{}, actual)

			actual, err = r.Persister().GetConcreteGrant(s.t1, grant.ID)
			require.NoError(t, err)
			require.NotEqual(t, trust.Grant{}, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetConsentRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			client := &client.Client{OutfacingID: "client-id"}
			f := newFlow(s.t1NID, client.OutfacingID, "sub", sqlxx.NullString(sessionID))
			require.NoError(t, r.Persister().CreateLoginSession(s.t1, &consent.LoginSession{ID: sessionID}))
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(f))

			req := &consent.ConsentRequest{
				ID:             "consent-request-id",
				LoginChallenge: sqlxx.NullString(f.ID),
				Skip:           false,
				Verifier:       "verifier",
				CSRF:           "csrf",
			}
			require.NoError(t, r.Persister().CreateConsentRequest(s.t1, req))

			actual, err := r.Persister().GetConsentRequest(s.t2, req.ID)
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetConsentRequest(s.t1, req.ID)
			require.NoError(t, err)
			require.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetFlow() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			client := &client.Client{OutfacingID: "client-id"}
			f := newFlow(s.t1NID, client.OutfacingID, "sub", sqlxx.NullString(sessionID))
			require.NoError(t, r.Persister().CreateLoginSession(s.t1, &consent.LoginSession{ID: sessionID}))
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(f))

			store, ok := r.Persister().(*persistencesql.Persister)
			if !ok {
				t.Fatal("type assertion failed")
			}

			_, err := store.GetFlow(s.t2, f.ID)
			require.Error(t, err)

			_, err = store.GetFlow(s.t1, f.ID)
			require.NoError(t, err)
		})
	}
}

func (s *PersisterTestSuite) TestGetFlowByConsentChallenge() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			client := &client.Client{OutfacingID: "client-id"}
			f := newFlow(s.t1NID, client.OutfacingID, "sub", sqlxx.NullString(sessionID))
			require.NoError(t, r.Persister().CreateLoginSession(s.t1, &consent.LoginSession{ID: sessionID}))
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(f))

			store, ok := r.Persister().(*persistencesql.Persister)
			if !ok {
				t.Fatal("type assertion failed")
			}

			_, err := store.GetFlowByConsentChallenge(s.t2, f.ConsentChallengeID.String())
			require.Error(t, err)

			_, err = store.GetFlowByConsentChallenge(s.t1, f.ConsentChallengeID.String())
			require.NoError(t, err)
		})
	}
}

func (s *PersisterTestSuite) TestGetForcedObfuscatedLoginSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			session := &consent.ForcedObfuscatedLoginSession{ClientID: client.OutfacingID}
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.Persister().CreateForcedObfuscatedLoginSession(s.t1, session))

			actual, err := r.Persister().GetForcedObfuscatedLoginSession(s.t2, client.OutfacingID, "")
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetForcedObfuscatedLoginSession(s.t1, client.OutfacingID, "")
			require.NoError(t, err)
			require.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetGrants() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()).String(),
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().AddKeySet(s.t1, "ks-id", ks))
			require.NoError(t, r.Persister().CreateGrant(s.t1, grant, ks.Keys[0]))

			actual, err := r.Persister().GetGrants(s.t2, 100, 0, "")
			require.NoError(t, err)
			require.Equal(t, 0, len(actual))

			actual, err = r.Persister().GetGrants(s.t1, 100, 0, "")
			require.NoError(t, err)
			require.Equal(t, 1, len(actual))
		})
	}
}

func (s *PersisterTestSuite) TestGetLoginRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			lr := consent.LoginRequest{ID: "lr-id", ClientID: client.OutfacingID, RequestedAt: time.Now()}

			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.ConsentManager().CreateLoginRequest(s.t1, &lr))
			f := flow.Flow{}
			require.NoError(t, r.Persister().Connection(context.Background()).Find(&f, lr.ID))
			require.Equal(t, s.t1NID, f.NID)

			actual, err := r.Persister().GetLoginRequest(s.t2, lr.ID)
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetLoginRequest(s.t1, lr.ID)
			require.NoError(t, err)
			require.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetLogoutRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			lr := consent.LogoutRequest{
				ID:       uuid.Must(uuid.NewV4()).String(),
				ClientID: sql.NullString{Valid: true, String: client.OutfacingID},
			}

			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.Persister().CreateLogoutRequest(s.t1, &lr))

			actual, err := r.Persister().GetLogoutRequest(s.t2, lr.ID)
			require.Error(t, err)
			require.Equal(t, &consent.LogoutRequest{}, actual)

			actual, err = r.Persister().GetLogoutRequest(s.t1, lr.ID)
			require.NoError(t, err)
			require.NotEqual(t, &consent.LogoutRequest{}, actual)
		})
	}
}

func (s *PersisterTestSuite) TestGetOpenIDConnectSession() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			request := fosite.NewRequest()
			request.SetID("request-id")
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			authorizeCode := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
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
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			request := fosite.NewRequest()
			request.SetID("request-id")
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			sig := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
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
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()).String(),
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().AddKeySet(s.t1, "ks-id", ks))
			require.NoError(t, r.Persister().CreateGrant(s.t1, grant, ks.Keys[0]))

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
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()).String(),
				Scope:     []string{"a", "b", "c"},
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().AddKeySet(s.t1, "ks-id", ks))
			require.NoError(t, r.Persister().CreateGrant(s.t1, grant, ks.Keys[0]))

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
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			ks := newKeySet("ks-id", "use")
			grant := trust.Grant{
				ID:        uuid.Must(uuid.NewV4()).String(),
				ExpiresAt: time.Now().Add(time.Hour),
				PublicKey: trust.PublicKey{Set: "ks-id", KeyID: ks.Keys[0].KeyID},
			}
			require.NoError(t, r.Persister().AddKeySet(s.t1, "ks-id", ks))
			require.NoError(t, r.Persister().CreateGrant(s.t1, grant, ks.Keys[0]))

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
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			client := &client.Client{OutfacingID: "client-id"}
			request := fosite.NewRequest()
			request.SetID("request-id")
			request.Client = &fosite.DefaultClient{ID: "client-id"}
			sig := uuid.Must(uuid.NewV4()).String()
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.Persister().CreateRefreshTokenSession(s.t1, sig, request))

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
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			ls := consent.LoginSession{ID: uuid.Must(uuid.NewV4()).String(), Remember: true}
			require.NoError(t, r.Persister().CreateLoginSession(s.t1, &ls))

			actual, err := r.Persister().GetRememberedLoginSession(s.t2, ls.ID)
			require.Error(t, err)
			require.Nil(t, actual)

			actual, err = r.Persister().GetRememberedLoginSession(s.t1, ls.ID)
			require.NoError(t, err)
			require.NotNil(t, actual)
		})
	}
}

func (s *PersisterTestSuite) TestHandleConsentRequest() {
	t := s.T()
	for k, r := range s.registries {
		t.Run(k, func(*testing.T) {
			sessionID := uuid.Must(uuid.NewV4()).String()
			client := &client.Client{OutfacingID: "client-id"}
			f := newFlow(s.t1NID, client.OutfacingID, "sub", sqlxx.NullString(sessionID))
			require.NoError(t, r.Persister().CreateLoginSession(s.t1, &consent.LoginSession{ID: sessionID}))
			require.NoError(t, r.Persister().CreateClient(s.t1, client))
			require.NoError(t, r.Persister().CreateClient(s.t2, client))
			require.NoError(t, r.Persister().Connection(context.Background()).Create(f))

			req := &consent.ConsentRequest{
				ID:             "consent-request-id",
				LoginChallenge: sqlxx.NullString(f.ID),
				Skip:           false,
				Verifier:       "verifier",
				CSRF:           "csrf",
			}

			hcr := &consent.HandledConsentRequest{
				ID:        req.ID,
				HandledAt: sqlxx.NullTime(time.Now()),
				Remember:  true,
			}
			require.NoError(t, r.Persister().CreateConsentRequest(s.t1, req))

			actualCR, err := r.Persister().HandleConsentRequest(s.t2, hcr)
			require.Error(t, err)
			require.Nil(t, actualCR)
			actual, err := r.Persister().FindGrantedAndRememberedConsentRequests(s.t1, client.OutfacingID, f.Subject)
			require.Error(t, err)
			require.Equal(t, 0, len(actual))

			actualCR, err = r.Persister().HandleConsentRequest(s.t1, hcr)
			require.NoError(t, err)
			require.NotNil(t, actualCR)
			actual, err = r.Persister().FindGrantedAndRememberedConsentRequests(s.t1, client.OutfacingID, f.Subject)
			require.NoError(t, err)
			require.Equal(t, 1, len(actual))
		})
	}
}

func TestExampleTestSuite(t *testing.T) {
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
		ConsentError:       &consent.RequestDeniedError{},
		State:              flow.FlowStateConsentUnused,
		LoginError:         &consent.RequestDeniedError{},
		Context:            sqlxx.JSONRawMessage{},
		AMR:                sqlxx.StringSlicePipeDelimiter{},
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

func newLogoutRequest() *consent.LogoutRequest {
	return &consent.LogoutRequest{
		ID: uuid.Must(uuid.NewV4()).String(),
	}
}

func newKey(ksID string, use string) jose.JSONWebKey {
	kg := &jwk.RS256Generator{}
	ks, err := kg.Generate(ksID, use)
	if err != nil {
		panic(err)
	}
	return ks.Keys[0]
}

func newKeySet(id string, use string) *jose.JSONWebKeySet {
	kg := &jwk.RS256Generator{}
	ks, err := kg.Generate(id, use)
	if err != nil {
		panic(err)
	}
	return ks
}

func newLoginSession() *consent.LoginSession {
	return &consent.LoginSession{
		ID:              uuid.Must(uuid.NewV4()).String(),
		AuthenticatedAt: sqlxx.NullTime(time.Time{}),
		Subject:         uuid.Must(uuid.NewV4()).String(),
		Remember:        false,
	}
}
