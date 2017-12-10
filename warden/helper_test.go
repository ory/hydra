// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package warden_test

import (
	"testing"
	"time"

	"os"

	"github.com/ory/fosite"
	"github.com/ory/fosite/storage"
	"github.com/ory/hydra/firewall"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/ory/hydra/warden"
	"github.com/ory/hydra/warden/group"
	"github.com/ory/ladon"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	accessRequestTestCases = []struct {
		req       *firewall.AccessRequest
		expectErr bool
		assert    func(*firewall.Context)
	}{
		{
			req: &firewall.AccessRequest{
				Subject:  "alice",
				Resource: "other-thing",
				Action:   "create",
				Context:  ladon.Context{},
			},
			expectErr: true,
		},
		{
			req: &firewall.AccessRequest{
				Subject:  "alice",
				Resource: "matrix",
				Action:   "delete",
				Context:  ladon.Context{},
			},
			expectErr: true,
		},
		{
			req: &firewall.AccessRequest{
				Subject:  "alice",
				Resource: "matrix",
				Action:   "create",
				Context:  ladon.Context{},
			},
			expectErr: false,
		},
	}
	wardens     = map[string]firewall.Firewall{}
	ladonWarden = pkg.LadonWarden(map[string]ladon.Policy{
		"1": &ladon.DefaultPolicy{
			ID:        "1",
			Subjects:  []string{"alice", "group1"},
			Resources: []string{"matrix", "forbidden_matrix", "rn:hydra:token<.*>"},
			Actions:   []string{"create", "decide"},
			Effect:    ladon.AllowAccess,
		},
		"2": &ladon.DefaultPolicy{
			ID:        "2",
			Subjects:  []string{"siri"},
			Resources: []string{"<.*>"},
			Actions:   []string{"decide"},
			Effect:    ladon.AllowAccess,
		},
		"3": &ladon.DefaultPolicy{
			ID:        "3",
			Subjects:  []string{"group1"},
			Resources: []string{"forbidden_matrix", "rn:hydra:token<.*>"},
			Actions:   []string{"create", "decide"},
			Effect:    ladon.DenyAccess,
		},
	})
	fositeStore                 = pkg.FositeStore()
	now                         = time.Now().UTC().Round(time.Second)
	tokens                      = pkg.Tokens(4)
	accessRequestTokenTestCases = []struct {
		token     string
		req       *firewall.TokenAccessRequest
		scopes    []string
		expectErr bool
		assert    func(*testing.T, *firewall.Context)
	}{
		{
			token:     "invalid",
			req:       &firewall.TokenAccessRequest{},
			scopes:    []string{},
			expectErr: true,
		},
		{
			token:     tokens[0][1],
			req:       &firewall.TokenAccessRequest{},
			scopes:    []string{"core"},
			expectErr: true,
		},
		{
			token:     tokens[0][1],
			req:       &firewall.TokenAccessRequest{},
			scopes:    []string{"foo"},
			expectErr: true,
		},
		{
			token: tokens[0][1],
			req: &firewall.TokenAccessRequest{
				Resource: "matrix",
				Action:   "create",
				Context:  ladon.Context{},
			},
			scopes:    []string{"foo"},
			expectErr: true,
		},
		{
			token: tokens[0][1],
			req: &firewall.TokenAccessRequest{
				Resource: "matrix",
				Action:   "delete",
				Context:  ladon.Context{},
			},
			scopes:    []string{"core"},
			expectErr: true,
		},
		{
			token: tokens[0][1],
			req: &firewall.TokenAccessRequest{
				Resource: "matrix",
				Action:   "create",
				Context:  ladon.Context{},
			},
			scopes:    []string{"illegal"},
			expectErr: true,
		},
		{
			token: tokens[0][1],
			req: &firewall.TokenAccessRequest{
				Resource: "matrix",
				Action:   "create",
				Context:  ladon.Context{},
			},
			scopes:    []string{"core"},
			expectErr: false,
			assert: func(t *testing.T, c *firewall.Context) {
				assert.Equal(t, "siri", c.ClientID)
				assert.Equal(t, "alice", c.Subject)
				assert.Equal(t, "tests", c.Issuer)
				assert.Equal(t, now.Add(time.Hour).Unix(), c.ExpiresAt.Unix())
				assert.Equal(t, now.Unix(), c.IssuedAt.Unix())
			},
		},
		{
			token: tokens[3][1],
			req: &firewall.TokenAccessRequest{
				Resource: "forbidden_matrix",
				Action:   "create",
				Context:  ladon.Context{},
			},
			scopes:    []string{"core"},
			expectErr: true,
		},
		{
			token: tokens[3][1],
			req: &firewall.TokenAccessRequest{
				Resource: "matrix",
				Action:   "create",
				Context:  ladon.Context{},
			},
			scopes:    []string{"core"},
			expectErr: false,
			assert: func(t *testing.T, c *firewall.Context) {
				assert.Equal(t, "siri", c.ClientID)
				assert.Equal(t, "ken", c.Subject)
				assert.Equal(t, "tests", c.Issuer)
				assert.Equal(t, now.Add(time.Hour).Unix(), c.ExpiresAt.Unix())
				assert.Equal(t, now.Unix(), c.IssuedAt.Unix())
			},
		},
	}
)

func createAccessTokenSession(subject, client string, token string, expiresAt time.Time, fs *storage.MemoryStore, scopes fosite.Arguments) {
	ar := fosite.NewAccessRequest(oauth2.NewSession(subject))
	ar.GrantedScopes = fosite.Arguments{"core"}
	if scopes != nil {
		ar.GrantedScopes = scopes
	}
	ar.RequestedAt = time.Now().UTC().Round(time.Second)
	ar.Client = &fosite.DefaultClient{ID: client}
	ar.Session.SetExpiresAt(fosite.AccessToken, expiresAt)
	ar.Session.(*oauth2.Session).Extra = map[string]interface{}{"foo": "bar"}
	fs.CreateAccessTokenSession(nil, token, ar)
}

func TestMain(m *testing.M) {
	wardens["local"] = &warden.LocalWarden{
		Warden: ladonWarden,
		L:      logrus.New(),
		OAuth2: &fosite.Fosite{
			Store: fositeStore,
			TokenIntrospectionHandlers: fosite.TokenIntrospectionHandlers{
				&warden.TokenValidator{
					CoreStrategy:  pkg.HMACStrategy,
					CoreStorage:   fositeStore,
					ScopeStrategy: fosite.HierarchicScopeStrategy,
				},
			},
			ScopeStrategy: fosite.HierarchicScopeStrategy,
		},
		Groups: &group.MemoryManager{
			Groups: map[string]group.Group{
				"group1": {
					ID:      "group1",
					Members: []string{"ken"},
				},
			},
		},
		Issuer:              "tests",
		AccessTokenLifespan: time.Hour,
	}

	createAccessTokenSession("alice", "siri", tokens[0][0], now.Add(time.Hour), fositeStore, fosite.Arguments{"core", "hydra.warden"})
	createAccessTokenSession("siri", "bob", tokens[1][0], now.Add(time.Hour), fositeStore, fosite.Arguments{"core", "hydra.warden"})
	createAccessTokenSession("siri", "doesnt-exist", tokens[2][0], now.Add(-time.Hour), fositeStore, fosite.Arguments{"core", "hydra.warden"})
	createAccessTokenSession("ken", "siri", tokens[3][0], now.Add(time.Hour), fositeStore, fosite.Arguments{"core", "hydra.warden"})

	os.Exit(m.Run())
}
