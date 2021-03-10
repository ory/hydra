package cli

import (
	"context"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
	hydra "github.com/ory/hydra/internal/httpclient/client"
	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/hydra/internal/httpclient/models"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
	"github.com/ory/x/urlx"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
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
var tokenSignature = "4c7c7e8b3a77ad0c3ec846a21653c48b45dbfa31"
var flushRefreshRequests = []*fosite.AccessRequest{
	{
		GrantTypes: []string{
			"refresh_token",
		},
		Request: fosite.Request{
			RequestedAt:    time.Now().Round(time.Second),
			ID:             "flush-refresh-1",
			Client:         &client.Client{OutfacingID: "foobar"},
			RequestedScope: []string{"offline"},
			GrantedScope:   []string{"offline"},
			Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
			Form: url.Values{
				"refresh_token": []string{fmt.Sprintf("%s.%s", "flush-refresh-1", tokenSignature)},
			},
		},
	},
	{
		GrantTypes: []string{
			"refresh_token",
		},
		Request: fosite.Request{
			RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
			ID:             "flush-refresh-2",
			Client:         &client.Client{OutfacingID: "foobar"},
			RequestedScope: []string{"offline"},
			GrantedScope:   []string{"offline"},
			Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
			Form: url.Values{
				"refresh_token": []string{fmt.Sprintf("%s.%s", "flush-refresh-2", tokenSignature)},
			},
		},
	},
	{
		GrantTypes: []string{
			"refresh_token",
		},
		Request: fosite.Request{
			RequestedAt:    time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
			ID:             "flush-refresh-3",
			Client:         &client.Client{OutfacingID: "foobar"},
			RequestedScope: []string{"offline"},
			GrantedScope:   []string{"offline"},
			Session:        &oauth2.Session{DefaultSession: &openid.DefaultSession{Subject: "bar"}},
			Form: url.Values{
				"refresh_token": []string{fmt.Sprintf("%s.%s", "flush-refresh-3", tokenSignature)},
			},
		},
	},
}

func TestJanitorHandler_Purge(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	conf.MustSet(config.KeyScopeStrategy, "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY")
	conf.MustSet(config.KeyIssuerURL, "http://hydra.localhost")
	conf.MustSet(config.KeyRefreshTokenLifespan, time.Hour*1)
	reg := internal.NewRegistryMemory(t, conf)

	cl := reg.ClientManager()
	store := reg.OAuth2Storage()

	// initialise the login/consent requests here

	h := oauth2.NewHandler(reg, conf)

	for _, r := range flushRequests {
		_ = cl.CreateClient(context.Background(), r.Client.(*client.Client))
		require.NoError(t, store.CreateAccessTokenSession(context.Background(), r.ID, r))
		// == Create here the login/register requests ==
	}

	for _, fr := range flushRefreshRequests {
		_ = cl.CreateClient(context.Background(), fr.Client.(*client.Client))
		require.NoError(t, store.CreateRefreshTokenSession(context.Background(), fr.ID, fr))
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

	_, err = store.GetRefreshTokenSession(ctx, "flush-refresh-1", ds)
	require.NoError(t, err)
	_, err = store.GetRefreshTokenSession(ctx, "flush-refresh-2", ds)
	require.NoError(t, err)
	_, err = store.GetRefreshTokenSession(ctx, "flush-refresh-3", ds)
	require.NoError(t, err)

	// == Add a check here to see if the data has been changed.==

	_, err = c.Admin.FlushInactiveOAuth2Tokens(admin.NewFlushInactiveOAuth2TokensParams().WithBody(&models.FlushInactiveOAuth2TokensRequest{NotAfter: strfmt.DateTime(time.Now().Add(-(lifespan + time.Hour/2)))}))
	require.NoError(t, err)

	_, err = store.GetAccessTokenSession(ctx, "flush-1", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-2", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-3", ds)
	require.Error(t, err)

	_, err = store.GetRefreshTokenSession(ctx, "flush-refresh-1", ds)
	require.NoError(t, err)
	_, err = store.GetRefreshTokenSession(ctx, "flush-refresh-2", ds)
	require.NoError(t, err)
	_, err = store.GetRefreshTokenSession(ctx, "flush-refresh-3", ds)
	require.Error(t, err)

	// == Add a check here to see if the data has been changed.==

	_, err = c.Admin.FlushInactiveOAuth2Tokens(admin.NewFlushInactiveOAuth2TokensParams().WithBody(&models.FlushInactiveOAuth2TokensRequest{NotAfter: strfmt.DateTime(time.Now())}))
	require.NoError(t, err)

	_, err = store.GetAccessTokenSession(ctx, "flush-1", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-2", ds)
	require.Error(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-3", ds)
	require.Error(t, err)

	_, err = store.GetRefreshTokenSession(ctx, "flush-refresh-1", ds)
	require.NoError(t, err)
	_, err = store.GetRefreshTokenSession(ctx, "flush-refresh-2", ds)
	require.Error(t, err)
	_, err = store.GetRefreshTokenSession(ctx, "flush-refresh-3", ds)
	require.Error(t, err)

	// == Add a check here to see if the data has been changed.==
}
