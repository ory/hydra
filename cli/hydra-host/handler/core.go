package handler

import (
	"github.com/codegangsta/cli"
	"github.com/gorilla/mux"
	accounts "github.com/ory-am/hydra/account/handler"
	"github.com/ory-am/hydra/jwt"
	"github.com/ory-am/hydra/middleware"
	clients "github.com/ory-am/hydra/oauth/client/handler"
	connections "github.com/ory-am/hydra/oauth/connection/handler"
	oauth "github.com/ory-am/hydra/oauth/handler"
	"github.com/ory-am/hydra/oauth/provider"
	policies "github.com/ory-am/hydra/policy/handler"
	"github.com/ory-am/ladon/guard"
	"net/http"
)

type Core struct {
	Ctx               *Context
	accountHandler    *accounts.Handler
	clientHandler     *clients.Handler
	connectionHandler *connections.Handler
	oauthHandler      *oauth.Handler
	policyHandler     *policies.Handler

	guard     guard.Guarder
	providers provider.Registry

	issuer   string
	audience string
}

func (c *Core) Start(ctx *cli.Context) {
	c.Ctx.Start()

	var private, public []byte
	j := jwt.New(private, public)
	m := middleware.New(c.Ctx.Policies, j)
	c.accountHandler = accounts.NewHandler(c.Ctx.Accounts, m)
	c.clientHandler = clients.NewHandler(c.Ctx.Osins, m)
	c.connectionHandler = connections.NewHandler(c.Ctx.Connections, m)
	c.providers = provider.NewRegistry(providers)
	c.oauthHandler = &oauth.Handler{
		Accounts:       c.Ctx.Accounts,
		Policies:       c.Ctx.Policies,
		Guard:          c.guard,
		Connections:    c.Ctx.Connections,
		Providers:      c.providers,
		Issuer:         c.issuer,
		Audience:       c.audience,
		JWT:            j,
		OAuthConfig:    oauth.DefaultConfig(),
		OAuthStore:     c.Ctx.Osins,
		States:         c.Ctx.States,
		SignUpLocation: locations["signUp"],
		SignInLocation: locations["signIn"],
	}

	extractor := m.ExtractAuthentication
	router := mux.NewRouter()
	c.accountHandler.SetRoutes(router, extractor)
	c.connectionHandler.SetRoutes(router, extractor)
	c.clientHandler.SetRoutes(router, extractor)
	c.oauthHandler.SetRoutes(router)

	http.Handle("/", router)
	http.ListenAndServe(listenOn, nil)
}
