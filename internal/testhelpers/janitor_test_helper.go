// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/oauth2/trust"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/configx"
	"github.com/ory/x/contextx"
	"github.com/ory/x/logrusx"
)

type JanitorConsentTestHelper struct {
	uniqueName           string
	flushLoginRequests   []*flow.LoginRequest
	flushConsentRequests []*flow.OAuth2ConsentRequest
	flushAccessRequests  []*fosite.Request
	flushRefreshRequests []*fosite.AccessRequest
	flushGrants          []*createGrantRequest
	conf                 *config.DefaultProvider
}

type createGrantRequest struct {
	grant trust.Grant
	pk    jose.JSONWebKey
}

const lifespan = time.Hour

func NewConsentJanitorTestHelper(t *testing.T, uniqueName string, opts ...configx.OptionModifier) *JanitorConsentTestHelper {
	conf := NewConfigurationWithDefaults(t, append([]configx.OptionModifier{configx.WithValues(map[string]any{
		config.KeyScopeStrategy:        "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY",
		config.KeyIssuerURL:            "http://hydra.localhost",
		config.KeyAccessTokenLifespan:  lifespan,
		config.KeyRefreshTokenLifespan: lifespan,
		config.KeyConsentRequestMaxAge: lifespan,
		config.KeyLogLevel:             "trace",
	})}, opts...)...)

	return &JanitorConsentTestHelper{
		uniqueName:           uniqueName,
		conf:                 conf,
		flushAccessRequests:  getAccessRequests(uniqueName, lifespan),
		flushRefreshRequests: getRefreshRequests(uniqueName, lifespan),
		flushGrants:          getGrantRequests(uniqueName, lifespan),
	}
}

func (j *JanitorConsentTestHelper) GetDSN() string {
	return j.conf.DSN()
}

func (j *JanitorConsentTestHelper) GetConfig() *config.DefaultProvider {
	return j.conf
}

var NotAfterTestCycles = map[string]time.Duration{
	"notAfter24h":   lifespan * 24,
	"notAfter1h30m": lifespan + time.Hour/2,
	"notAfterNow":   0,
}

func (j *JanitorConsentTestHelper) GetNotAfterTestCycles() map[string]time.Duration {
	return map[string]time.Duration{}
}

func (j *JanitorConsentTestHelper) GetRegistry(ctx context.Context, dbname string) (driver.Registry, error) {
	j.conf.MustSet(ctx, config.KeyDSN, fmt.Sprintf("sqlite://file:%s?mode=memory&_fk=true&cache=shared", dbname))
	return driver.NewRegistryFromDSN(ctx, j.conf, logrusx.New("test_hydra", "master"), false, true, &contextx.Default{})
}

func (j *JanitorConsentTestHelper) AccessTokenNotAfterSetup(ctx context.Context, cl client.Manager, store x.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		// Create access token clients and session
		for _, r := range j.flushAccessRequests {
			require.NoError(t, cl.CreateClient(ctx, r.Client.(*client.Client)))
			require.NoError(t, store.CreateAccessTokenSession(ctx, r.ID, r))
		}

	}
}

func (j *JanitorConsentTestHelper) AccessTokenNotAfterValidate(ctx context.Context, notAfter time.Time, store x.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		var err error
		ds := new(oauth2.Session)

		accessTokenLifespan := time.Now().Round(time.Second).Add(-j.conf.GetAccessTokenLifespan(ctx))

		for _, r := range j.flushAccessRequests {
			t.Logf("access flush check: %s", r.ID)
			_, err = store.GetAccessTokenSession(ctx, r.ID, ds)
			if j.notAfterCheck(notAfter, accessTokenLifespan, r.RequestedAt) {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}
	}
}

func (j *JanitorConsentTestHelper) RefreshTokenNotAfterSetup(ctx context.Context, cl client.Manager, store x.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		// Create refresh token clients and session
		for _, fr := range j.flushRefreshRequests {
			require.NoError(t, cl.CreateClient(ctx, fr.Client.(*client.Client)))
			require.NoError(t, store.CreateRefreshTokenSession(ctx, fr.ID, "", fr))
		}
	}
}

func (j *JanitorConsentTestHelper) RefreshTokenNotAfterValidate(ctx context.Context, notAfter time.Time, store x.FositeStorer) func(t *testing.T) {
	return func(t *testing.T) {
		var err error
		ds := new(oauth2.Session)

		refreshTokenLifespan := time.Now().Round(time.Second).Add(-j.conf.GetRefreshTokenLifespan(ctx))

		for _, r := range j.flushRefreshRequests {
			t.Logf("refresh flush check: %s", r.ID)
			_, err = store.GetRefreshTokenSession(ctx, r.ID, ds)
			if j.notAfterCheck(notAfter, refreshTokenLifespan, r.RequestedAt) {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}
	}
}

func (j *JanitorConsentTestHelper) GrantNotAfterSetup(ctx context.Context, gr trust.GrantManager) func(t *testing.T) {
	return func(t *testing.T) {
		for _, fg := range j.flushGrants {
			require.NoError(t, gr.CreateGrant(ctx, fg.grant, fg.pk))
		}
	}
}

func (j *JanitorConsentTestHelper) GrantNotAfterValidate(ctx context.Context, notAfter time.Time, gr trust.GrantManager) func(t *testing.T) {
	return func(t *testing.T) {
		var err error

		// flush won't delete grants that have not yet expired, so use now to check that
		deleteUntil := time.Now().Round(time.Second)
		if deleteUntil.After(notAfter) {
			deleteUntil = notAfter
		}

		for _, r := range j.flushGrants {
			t.Logf("grant flush check: %s", r.grant.Issuer)
			_, err = gr.GetConcreteGrant(ctx, r.grant.ID)

			if deleteUntil.After(r.grant.ExpiresAt) {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}
	}
}

func (j *JanitorConsentTestHelper) GetConsentRequestLifespan(ctx context.Context) time.Duration {
	return j.conf.ConsentRequestMaxAge(ctx)
}

func (j *JanitorConsentTestHelper) GetAccessTokenLifespan(ctx context.Context) time.Duration {
	return j.conf.GetAccessTokenLifespan(ctx)
}

func (j *JanitorConsentTestHelper) GetRefreshTokenLifespan(ctx context.Context) time.Duration {
	return j.conf.GetRefreshTokenLifespan(ctx)
}

func (j *JanitorConsentTestHelper) notAfterCheck(notAfter time.Time, lifespan time.Time, requestedAt time.Time) bool {
	// The database deletes where requested_at time is smaller than the lowest between notAfter and consent-request-lifespan
	// thus we get the lowest value here first to compare later to requested_at
	var lesser time.Time
	// if the lowest between notAfter and consent-request-lifespan is greater than requested_at
	// then the it should expect the value to be deleted.
	if notAfter.Unix() < lifespan.Unix() {
		lesser = notAfter
	} else {
		lesser = lifespan
	}

	// true: value has been deleted
	// false: value still exists
	return lesser.Unix() > requestedAt.Unix()
}

func getAccessRequests(uniqueName string, lifespan time.Duration) []*fosite.Request {
	return []*fosite.Request{
		{
			ID:             fmt.Sprintf("%s_flush-access-1", uniqueName),
			RequestedAt:    time.Now().Round(time.Second),
			Client:         &client.Client{ID: fmt.Sprintf("%s_flush-access-1", uniqueName)},
			RequestedScope: fosite.Arguments{"fa", "ba"},
			GrantedScope:   fosite.Arguments{"fa", "ba"},
			Form:           url.Values{"foo": []string{"bar", "baz"}},
			Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
		},
		{
			ID:             fmt.Sprintf("%s_flush-access-2", uniqueName),
			RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
			Client:         &client.Client{ID: fmt.Sprintf("%s_flush-access-2", uniqueName)},
			RequestedScope: fosite.Arguments{"fa", "ba"},
			GrantedScope:   fosite.Arguments{"fa", "ba"},
			Form:           url.Values{"foo": []string{"bar", "baz"}},
			Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
		},
		{
			ID:             fmt.Sprintf("%s_flush-access-3", uniqueName),
			RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
			Client:         &client.Client{ID: fmt.Sprintf("%s_flush-access-3", uniqueName)},
			RequestedScope: fosite.Arguments{"fa", "ba"},
			GrantedScope:   fosite.Arguments{"fa", "ba"},
			Form:           url.Values{"foo": []string{"bar", "baz"}},
			Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
		},
	}
}

func getRefreshRequests(uniqueName string, lifespan time.Duration) []*fosite.AccessRequest {
	var tokenSignature = "4c7c7e8b3a77ad0c3ec846a21653c48b45dbfa31" //nolint:gosec
	return []*fosite.AccessRequest{
		{
			GrantTypes: []string{
				"refresh_token",
			},
			Request: fosite.Request{
				RequestedAt:    time.Now().Round(time.Second),
				ID:             fmt.Sprintf("%s_flush-refresh-1", uniqueName),
				Client:         &client.Client{ID: fmt.Sprintf("%s_flush-refresh-1", uniqueName)},
				RequestedScope: []string{"offline"},
				GrantedScope:   []string{"offline"},
				Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
				Form: url.Values{
					"refresh_token": []string{fmt.Sprintf("%s.%s", fmt.Sprintf("%s_flush-refresh-1", uniqueName), tokenSignature)},
				},
			},
		},
		{
			GrantTypes: []string{
				"refresh_token",
			},
			Request: fosite.Request{
				RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
				ID:             fmt.Sprintf("%s_flush-refresh-2", uniqueName),
				Client:         &client.Client{ID: fmt.Sprintf("%s_flush-refresh-2", uniqueName)},
				RequestedScope: []string{"offline"},
				GrantedScope:   []string{"offline"},
				Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
				Form: url.Values{
					"refresh_token": []string{fmt.Sprintf("%s.%s", fmt.Sprintf("%s_flush-refresh-2", uniqueName), tokenSignature)},
				},
			},
		},
		{
			GrantTypes: []string{
				"refresh_token",
			},
			Request: fosite.Request{
				RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
				ID:             fmt.Sprintf("%s_flush-refresh-3", uniqueName),
				Client:         &client.Client{ID: fmt.Sprintf("%s_flush-refresh-3", uniqueName)},
				RequestedScope: []string{"offline"},
				GrantedScope:   []string{"offline"},
				Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
				Form: url.Values{
					"refresh_token": []string{fmt.Sprintf("%s.%s", fmt.Sprintf("%s_flush-refresh-3", uniqueName), tokenSignature)},
				},
			},
		},
	}
}

func getGrantRequests(uniqueName string, lifespan time.Duration) []*createGrantRequest {
	return []*createGrantRequest{
		{
			grant: trust.Grant{
				ID:      uuid.Must(uuid.NewV4()),
				Issuer:  fmt.Sprintf("%s_flush-grant-iss-1", uniqueName),
				Subject: fmt.Sprintf("%s_flush-grant-sub-1", uniqueName),
				Scope:   []string{"foo", "bar"},
				PublicKey: trust.PublicKey{
					Set:   fmt.Sprintf("%s_flush-grant-iss-1", uniqueName),
					KeyID: fmt.Sprintf("%s_flush-grant-kid-1", uniqueName),
				},
				CreatedAt: time.Now().Round(time.Second),
				ExpiresAt: time.Now().Round(time.Second).Add(lifespan),
			},
			pk: jose.JSONWebKey{
				Key:   []byte("asdf"),
				KeyID: fmt.Sprintf("%s_flush-grant-kid-1", uniqueName),
			},
		},
		{
			grant: trust.Grant{
				ID:      uuid.Must(uuid.NewV4()),
				Issuer:  fmt.Sprintf("%s_flush-grant-iss-2", uniqueName),
				Subject: fmt.Sprintf("%s_flush-grant-sub-2", uniqueName),
				Scope:   []string{"foo", "bar"},
				PublicKey: trust.PublicKey{
					Set:   fmt.Sprintf("%s_flush-grant-iss-2", uniqueName),
					KeyID: fmt.Sprintf("%s_flush-grant-kid-2", uniqueName),
				},
				CreatedAt: time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
				ExpiresAt: time.Now().Round(time.Second).Add(-(lifespan + time.Minute)).Add(lifespan),
			},
			pk: jose.JSONWebKey{
				Key:   []byte("asdf"),
				KeyID: fmt.Sprintf("%s_flush-grant-kid-2", uniqueName),
			},
		},
		{
			grant: trust.Grant{
				ID:      uuid.Must(uuid.NewV4()),
				Issuer:  fmt.Sprintf("%s_flush-grant-iss-3", uniqueName),
				Subject: fmt.Sprintf("%s_flush-grant-sub-3", uniqueName),
				Scope:   []string{"foo", "bar"},
				PublicKey: trust.PublicKey{
					Set:   fmt.Sprintf("%s_flush-grant-iss-3", uniqueName),
					KeyID: fmt.Sprintf("%s_flush-grant-kid-3", uniqueName),
				},
				CreatedAt: time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
				ExpiresAt: time.Now().Round(time.Second).Add(-(lifespan + time.Hour)).Add(lifespan),
			},
			pk: jose.JSONWebKey{
				Key:   []byte("asdf"),
				KeyID: fmt.Sprintf("%s_flush-grant-kid-3", uniqueName),
			},
		},
	}
}
