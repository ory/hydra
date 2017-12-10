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

package compose

import (
	"net/http"
	"time"

	"github.com/ory/fosite"
	foauth2 "github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/storage"
	"github.com/ory/hydra/firewall"
	. "github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/ory/hydra/warden"
	"github.com/ory/hydra/warden/group"
	"github.com/ory/ladon"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func NewMockFirewall(issuer string, subject string, scopes fosite.Arguments, p ...ladon.Policy) (firewall.Firewall, *http.Client) {
	return NewMockFirewallWithStore(issuer, subject, scopes, pkg.FositeStore(), p...)
}

func NewMockFirewallWithStore(issuer string, subject string, scopes fosite.Arguments, storage *storage.MemoryStore, p ...ladon.Policy) (firewall.Firewall, *http.Client) {
	tokens := pkg.Tokens(1)

	fositeStore := storage
	ps := map[string]ladon.Policy{}

	for _, x := range p {
		ps[x.GetID()] = x
	}
	ladonWarden := pkg.LadonWarden(ps)

	ar := fosite.NewAccessRequest(NewSession(subject))
	ar.GrantedScopes = scopes
	fositeStore.CreateAccessTokenSession(nil, tokens[0][0], ar)

	conf := &oauth2.Config{Scopes: scopes, Endpoint: oauth2.Endpoint{}}
	l := logrus.New()
	l.Level = logrus.DebugLevel

	return &warden.LocalWarden{
			Warden: ladonWarden,
			OAuth2: &fosite.Fosite{
				Store: fositeStore,
				TokenIntrospectionHandlers: fosite.TokenIntrospectionHandlers{
					&foauth2.CoreValidator{
						CoreStrategy:  pkg.HMACStrategy,
						CoreStorage:   fositeStore,
						ScopeStrategy: fosite.HierarchicScopeStrategy,
					},
				},
				ScopeStrategy: fosite.HierarchicScopeStrategy,
			},
			Issuer:              issuer,
			AccessTokenLifespan: time.Hour,
			Groups:              group.NewMemoryManager(),
			L:                   l,
		}, conf.Client(oauth2.NoContext, &oauth2.Token{
			AccessToken: tokens[0][1],
			Expiry:      time.Now().UTC().Add(time.Hour),
			TokenType:   "bearer",
		})
}
