package oauth2

import (
	"context"
	"fmt"
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/x"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
	"time"
)

type cleanupRoutine = func() error

type JanitorOauthTestHelper struct {
	flushAccessRequests  []*fosite.Request
	flushRefreshRequests []*fosite.AccessRequest
	Lifespan             time.Duration
}

func NewOauthJanitorTestHelper(uniqueName string) *JanitorOauthTestHelper {
	var lifespan = time.Hour
	var tokenSignature = "4c7c7e8b3a77ad0c3ec846a21653c48b45dbfa31"

	return &JanitorOauthTestHelper{
		flushAccessRequests: []*fosite.Request{
			{
				ID:             fmt.Sprintf("%s_flush-access-1", uniqueName),
				RequestedAt:    time.Now().Round(time.Second),
				Client:         &client.Client{OutfacingID: fmt.Sprintf("%s_flush-access-1", uniqueName)},
				RequestedScope: fosite.Arguments{"fa", "ba"},
				GrantedScope:   fosite.Arguments{"fa", "ba"},
				Form:           url.Values{"foo": []string{"bar", "baz"}},
				Session:        &Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
			},
			{
				ID:             fmt.Sprintf("%s_flush-access-2", uniqueName),
				RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
				Client:         &client.Client{OutfacingID: fmt.Sprintf("%s_flush-access-2", uniqueName)},
				RequestedScope: fosite.Arguments{"fa", "ba"},
				GrantedScope:   fosite.Arguments{"fa", "ba"},
				Form:           url.Values{"foo": []string{"bar", "baz"}},
				Session:        &Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
			},
			{
				ID:             fmt.Sprintf("%s_flush-access-3", uniqueName),
				RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
				Client:         &client.Client{OutfacingID: fmt.Sprintf("%s_flush-access-3", uniqueName)},
				RequestedScope: fosite.Arguments{"fa", "ba"},
				GrantedScope:   fosite.Arguments{"fa", "ba"},
				Form:           url.Values{"foo": []string{"bar", "baz"}},
				Session:        &Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
			},
		},
		flushRefreshRequests: []*fosite.AccessRequest{
			{
				GrantTypes: []string{
					"refresh_token",
				},
				Request: fosite.Request{
					RequestedAt:    time.Now().Round(time.Second),
					ID:             fmt.Sprintf("%s_flush-refresh-1", uniqueName),
					Client:         &client.Client{OutfacingID: fmt.Sprintf("%s_flush-refresh-1", uniqueName)},
					RequestedScope: []string{"offline"},
					GrantedScope:   []string{"offline"},
					Session:        &Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
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
					Client:         &client.Client{OutfacingID: fmt.Sprintf("%s_flush-refresh-2", uniqueName)},
					RequestedScope: []string{"offline"},
					GrantedScope:   []string{"offline"},
					Session:        &Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
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
					Client:         &client.Client{OutfacingID: fmt.Sprintf("%s_flush-refresh-3", uniqueName)},
					RequestedScope: []string{"offline"},
					GrantedScope:   []string{"offline"},
					Session:        &Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
					Form: url.Values{
						"refresh_token": []string{fmt.Sprintf("%s.%s", fmt.Sprintf("%s_flush-refresh-3", uniqueName), tokenSignature)},
					},
				},
			},
		},
		Lifespan: lifespan,
	}
}

func (j *JanitorOauthTestHelper) AccessTokenNotAfter(notAfter time.Time, tokenLifespan time.Duration, cr cleanupRoutine, cl client.Manager, store x.FositeStorer) func(t *testing.T) {
	var err error

	ctx := context.Background()

	return func(t *testing.T) {

		var greater time.Time
		if notAfter.Unix() < time.Now().Add(-tokenLifespan).Unix() {
			greater = time.Now().Add(-tokenLifespan)
		} else {
			greater = notAfter
		}

		// Create access token clients and session
		for _, r := range j.flushAccessRequests {
			require.NoError(t, cl.CreateClient(ctx, r.Client.(*client.Client)))
			require.NoError(t, store.CreateAccessTokenSession(ctx, r.ID, r))
		}

		ds := new(Session)

		// run the cleanup routine
		require.NoError(t, cr())

		for _, r := range j.flushAccessRequests {
			t.Logf("access flush check: %s", r.ID)
			_, err = store.GetAccessTokenSession(ctx, r.ID, ds)
			if greater.Unix() > r.RequestedAt.Unix() {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}

	}
}

func (j *JanitorOauthTestHelper) RefreshTokenNotAfter(notAfter time.Time, tokenLifespan time.Duration, cr cleanupRoutine, cl client.Manager, store x.FositeStorer) func(t *testing.T) {
	var err error

	ctx := context.Background()

	return func(t *testing.T) {

		var greater time.Time
		if notAfter.Unix() < time.Now().Add(-tokenLifespan).Unix() {
			greater = time.Now().Add(-tokenLifespan)
		} else {
			greater = notAfter
		}

		// Create refresh token clients and session
		for _, fr := range j.flushRefreshRequests {
			require.NoError(t, cl.CreateClient(ctx, fr.Client.(*client.Client)))
			require.NoError(t, store.CreateRefreshTokenSession(ctx, fr.ID, fr))
		}

		ds := new(Session)

		// run the cleanup routine
		require.NoError(t, cr())

		for _, r := range j.flushRefreshRequests {
			t.Logf("refresh flush check: %s", r.ID)
			_, err = store.GetRefreshTokenSession(ctx, r.ID, ds)
			if greater.Unix() > r.RequestedAt.Unix() {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}

	}
}
