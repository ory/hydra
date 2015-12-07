package handler

import (
	"github.com/codegangsta/cli"
	"github.com/gorilla/mux"
	accounts "github.com/ory-am/hydra/account/handler"
	"github.com/ory-am/hydra/jwt"
	"github.com/ory-am/hydra/middleware/host"
	middleware "github.com/ory-am/hydra/middleware/host"
	clients "github.com/ory-am/hydra/oauth/client/handler"
	connections "github.com/ory-am/hydra/oauth/connection/handler"
	oauth "github.com/ory-am/hydra/oauth/handler"
	"github.com/ory-am/hydra/oauth/provider"
	policies "github.com/ory-am/hydra/policy/handler"
	"github.com/ory-am/ladon/guard"

	"fmt"
	"golang.org/x/net/http2"
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

func (c *Core) Start(ctx *cli.Context) error {
	c.Ctx.Start()

	private, err := jwt.LoadCertificate(jwtPrivateKeyPath)
	if err != nil {
		return fmt.Errorf("Could not load private key: %s", err)
	}

	public, err := jwt.LoadCertificate(jwtPublicKeyPath)
	if err != nil {
		return fmt.Errorf("Could not load public key: %s", err)
	}

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
		Middleware:     host.New(c.Ctx.Policies, j),
	}

	extractor := m.ExtractAuthentication
	router := mux.NewRouter()
	c.accountHandler.SetRoutes(router, extractor)
	c.connectionHandler.SetRoutes(router, extractor)
	c.clientHandler.SetRoutes(router, extractor)
	c.oauthHandler.SetRoutes(router, extractor)

	if forceHTTP == "force" {
		http.Handle("/", router)
		http.ListenAndServe(listenOn, nil)
		return nil
	}

	http.Handle("/", router)
	srv := &http.Server{
		Addr: listenOn,
	}
	http2.ConfigureServer(srv, &http2.Server{})
	err = srv.ListenAndServeTLS(tlsCertPath, tlsKeyPath)
	if err != nil {
		return fmt.Errorf("Could not serve HTTP/2 server because %s", err)
	}
	return nil
}
