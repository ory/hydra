package handler

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	gojwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	accounts "github.com/ory-am/hydra/account/handler"
	"github.com/ory-am/hydra/jwt"
	"github.com/ory-am/hydra/middleware/host"
	middleware "github.com/ory-am/hydra/middleware/host"
	clients "github.com/ory-am/hydra/endpoint/client/handler"
	connections "github.com/ory-am/hydra/endpoint/connection/handler"
	oauth "github.com/ory-am/hydra/endpoint/handler"
	"github.com/ory-am/hydra/endpoint/connector"
	policies "github.com/ory-am/hydra/policy/handler"
	"github.com/ory-am/ladon/guard"

	"fmt"
	"net/http"
	"strconv"

	"crypto/tls"
	"github.com/RangelReale/osin"
	"github.com/ory-am/common/pkg"
)

type Core struct {
	Ctx               Context
	accountHandler    *accounts.Handler
	clientHandler     *clients.Handler
	connectionHandler *connections.Handler
	oauthHandler      *oauth.Handler
	policyHandler     *policies.Handler

	guard             guard.Guarder
	providers         connector.Registry

	issuer            string
	audience          string
}

func osinConfig() (conf *osin.ServerConfig, err error) {
	conf = osin.NewServerConfig()
	lifetime, err := strconv.Atoi(accessTokenLifetime)
	if err != nil {
		return nil, err
	}
	conf.AccessExpiration = int32(lifetime)

	conf.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{
		osin.CODE,
		osin.TOKEN,
	}
	conf.AllowedAccessTypes = osin.AllowedAccessType{
		osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN,
		osin.PASSWORD,
		osin.CLIENT_CREDENTIALS,
	}
	conf.AllowGetAccessRequest = false
	conf.AllowClientSecretInParams = false
	conf.ErrorStatusCode = http.StatusInternalServerError
	conf.RedirectUriSeparator = "|"
	return conf, nil
}

func (c *Core) Start(ctx *cli.Context) error {
	// Start the database backend
	if err := c.Ctx.Start(); err != nil {
		return fmt.Errorf("Could not start context: %s", err)
	}

	private, err := jwt.LoadCertificate(jwtPrivateKey)
	if err != nil {
		return fmt.Errorf("Could not load private key: %s", err)
	}

	public, err := jwt.LoadCertificate(jwtPublicKey)
	if err != nil {
		return fmt.Errorf("Could not load public key: %s", err)
	}

	fmt.Printf("Key %s", public)

	if _, err = gojwt.ParseRSAPublicKeyFromPEM(public); err != nil {
		return fmt.Errorf("Not a valid public key: %s", err)
	} else if _, err = gojwt.ParseRSAPrivateKeyFromPEM(private); err != nil {
		return fmt.Errorf("Not a valid private key: %s", err)
	}

	osinConf, err := osinConfig()
	if err != nil {
		return fmt.Errorf("Could not configure server: %s", err)
	}

	j := jwt.New(private, public)
	m := middleware.New(c.Ctx.GetPolicies(), j)
	c.guard = new(guard.Guard)
	c.accountHandler = accounts.NewHandler(c.Ctx.GetAccounts(), m)
	c.clientHandler = clients.NewHandler(c.Ctx.GetOsins(), m)
	c.connectionHandler = connections.NewHandler(c.Ctx.GetConnections(), m)
	c.providers = connector.NewRegistry(providers)
	c.policyHandler = policies.NewHandler(c.Ctx.GetPolicies(), m, c.guard, j, c.Ctx.GetOsins())
	c.oauthHandler = &oauth.Handler{
		Accounts:       c.Ctx.GetAccounts(),
		Policies:       c.Ctx.GetPolicies(),
		Guard:          c.guard,
		Connections:    c.Ctx.GetConnections(),
		Providers:      c.providers,
		Issuer:         c.issuer,
		Audience:       c.audience,
		JWT:            j,
		OAuthConfig:    osinConf,
		OAuthStore:     c.Ctx.GetOsins(),
		States:         c.Ctx.GetStates(),
		SignUpLocation: locations["signUp"],
		SignInLocation: locations["signIn"],
		Middleware:     host.New(c.Ctx.GetPolicies(), j),
	}

	extractor := m.ExtractAuthentication
	router := mux.NewRouter()
	c.accountHandler.SetRoutes(router, extractor)
	c.connectionHandler.SetRoutes(router, extractor)
	c.clientHandler.SetRoutes(router, extractor)
	c.oauthHandler.SetRoutes(router, extractor)
	c.policyHandler.SetRoutes(router, extractor)

	// TODO un-hack this, add database check, add error response
	router.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		pkg.WriteJSON(w, &struct {
			Status string `json:"status"`
		}{
			Status: "alive",
		})
	})

	log.Infoln("Hydra initialized, starting listeners...")

	if forceHTTP == "force" {
		http.Handle("/", router)
		log.Warn("You're using HTTP without TLS encryption. This is dangerously unsafe and you should not do this.")
		if err := http.ListenAndServe(listenOn, nil); err != nil {
			return fmt.Errorf("Could not serve HTTP server because %s", err)
		}
		return nil
	}

	var cert tls.Certificate
	if cert, err = tls.LoadX509KeyPair(tlsCert, tlsKey); err != nil {
		if cert, err = tls.X509KeyPair([]byte(tlsCert), []byte(tlsKey)); err != nil {
			return fmt.Errorf("Could not load or parse TLS key pair because %s", err)
		}
	}
	srv := &http.Server{
		Addr: listenOn,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{
				cert,
			},
		},
	}

	http.Handle("/", router)
	if err := srv.ListenAndServeTLS("", ""); err != nil {
		return fmt.Errorf("Could not serve HTTP/2 server because %s", err)
	}

	return nil
}
