package compose

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/ory/fosite"
	foauth2 "github.com/ory/fosite/handler/oauth2"
	"github.com/ory/hydra/firewall"
	. "github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/ory/hydra/warden"
	"github.com/ory/hydra/warden/group"
	"github.com/ory/ladon"
	"golang.org/x/oauth2"
)

func NewMockFirewall(issuer string, subject string, scopes fosite.Arguments, p ...ladon.Policy) (firewall.Firewall, *http.Client) {
	tokens := pkg.Tokens(1)

	fositeStore := pkg.FositeStore()
	ps := map[string]ladon.Policy{}

	for _, x := range p {
		ps[x.GetID()] = x
	}
	ladonWarden := pkg.LadonWarden(ps)

	ar := fosite.NewAccessRequest(NewSession(subject))
	ar.GrantedScopes = scopes
	fositeStore.CreateAccessTokenSession(nil, tokens[0][0], ar)

	conf := &oauth2.Config{Scopes: scopes, Endpoint: oauth2.Endpoint{}}

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
			L:                   logrus.New(),
		}, conf.Client(oauth2.NoContext, &oauth2.Token{
			AccessToken: tokens[0][1],
			Expiry:      time.Now().Add(time.Hour),
			TokenType:   "bearer",
		})
}
