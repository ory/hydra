package oauth2_test

import (
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/herodot"
	hc "github.com/ory/hydra/client"
	hcompose "github.com/ory/hydra/compose"
	. "github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/ory/ladon"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var hasher = &fosite.BCrypt{}

var store = &FositeMemoryStore{
	Manager: &hc.MemoryManager{
		Clients: map[string]hc.Client{},
		Hasher:  hasher,
	},
	AuthorizeCodes: make(map[string]fosite.Requester),
	IDSessions:     make(map[string]fosite.Requester),
	AccessTokens:   make(map[string]fosite.Requester),
	RefreshTokens:  make(map[string]fosite.Requester),
}

var fc = &compose.Config{
	AccessTokenLifespan: time.Second,
}

var handler = &Handler{
	OAuth2: compose.Compose(
		fc,
		store,
		&compose.CommonStrategy{
			CoreStrategy:               compose.NewOAuth2HMACStrategy(fc, []byte("some super secret secret")),
			OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(pkg.MustRSAKey()),
		},
		nil,
		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2AuthorizeImplicitFactory,
		compose.OAuth2ClientCredentialsGrantFactory,
		compose.OAuth2RefreshTokenGrantFactory,
		compose.OpenIDConnectExplicitFactory,
		compose.OpenIDConnectHybridFactory,
		compose.OpenIDConnectImplicitFactory,
		compose.OAuth2TokenRevocationFactory,
		compose.OAuth2TokenIntrospectionFactory,
	),
	Consent: &DefaultConsentStrategy{
		Issuer:                   "http://hydra.localhost",
		ConsentManager:           consentManager,
		DefaultChallengeLifespan: time.Hour,
		DefaultIDTokenLifespan:   time.Hour * 24,
	},
	CookieStore:   sessions.NewCookieStore([]byte("foo-secret")),
	ForcedHTTP:    true,
	L:             logrus.New(),
	ScopeStrategy: fosite.HierarchicScopeStrategy,
}

var router = httprouter.New()
var ts *httptest.Server
var oauthConfig *oauth2.Config
var oauthClientConfig *clientcredentials.Config

var localWarden, httpClient = hcompose.NewMockFirewall("foo", "app-client", fosite.Arguments{ConsentScope}, &ladon.DefaultPolicy{
	ID:        "1",
	Subjects:  []string{"app-client"},
	Resources: []string{"rn:hydra:oauth2:consent:requests:<.*>"},
	Actions:   []string{"get", "accept", "reject"},
	Effect:    ladon.AllowAccess,
})

var consentHandler *ConsentSessionHandler
var consentManager = NewConsentRequestMemoryManager()
var consentClient *HTTPConsentManager

func init() {
	consentHandler = &ConsentSessionHandler{
		H: herodot.NewJSONWriter(nil),
		W: localWarden,
		M: consentManager,
	}

	ts = httptest.NewServer(router)
	handler.Issuer = ts.URL

	handler.SetRoutes(router)
	consentHandler.SetRoutes(router)

	h, _ := hasher.Hash([]byte("secret"))
	u, _ := url.Parse(ts.URL + ConsentRequestPath)
	consentClient = &HTTPConsentManager{Client: httpClient, Endpoint: u}

	c, _ := url.Parse(ts.URL + "/consent")
	handler.ConsentURL = *c

	store.Manager.(*hc.MemoryManager).Clients["app-client"] = hc.Client{
		ID:            "app-client",
		Secret:        string(h),
		RedirectURIs:  []string{ts.URL + "/callback"},
		ResponseTypes: []string{"id_token", "code", "token"},
		GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
		Scope:         "hydra.* offline",
	}

	oauthConfig = &oauth2.Config{
		ClientID:     "app-client",
		ClientSecret: "secret",
		Endpoint: oauth2.Endpoint{
			AuthURL:  ts.URL + "/oauth2/auth",
			TokenURL: ts.URL + "/oauth2/token",
		},
		RedirectURL: ts.URL + "/callback",
		Scopes:      []string{"hydra.*", "offline"},
	}

	oauthClientConfig = &clientcredentials.Config{
		ClientID:     "app-client",
		ClientSecret: "secret",
		TokenURL:     ts.URL + "/oauth2/token",
		Scopes:       []string{"hydra.consent", "offline"},
	}
}
