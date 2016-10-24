package compose

import (
	"net/http"
	"time"

	"github.com/ory-am/fosite"
	foauth2 "github.com/ory-am/fosite/handler/oauth2"
	"github.com/ory-am/hydra/firewall"
	. "github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/hydra/warden"
	"github.com/ory-am/ladon"
	"golang.org/x/oauth2"
)

func NewFirewall(issuer string, subject string, scopes fosite.Arguments, p ...ladon.Policy) (firewall.Firewall, *http.Client) {
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
		}, conf.Client(oauth2.NoContext, &oauth2.Token{
			AccessToken: tokens[0][1],
			Expiry:      time.Now().Add(time.Hour),
			TokenType:   "bearer",
		})
}
