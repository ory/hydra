package oauth2_test

import (
	"net/http/httptest"
	"time"

	"fmt"
	"net/url"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-errors/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/core"
	"github.com/ory-am/fosite/handler/core/client"
	"github.com/ory-am/fosite/handler/core/explicit"
	"github.com/ory-am/fosite/handler/core/strategy"
	"github.com/ory-am/fosite/hash"
	"github.com/ory-am/fosite/token/hmac"
	"github.com/ory-am/hydra/key"
	. "github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var store = pkg.FositeStore()

var keyStrategy = &key.DefaultKeyStrategy{
	AsymmetricKeyStrategy: &key.RSAPEMStrategy{},
	SymmetricKeyStrategy:  &key.SHAStrategy{},
}

var keyManager = &key.MemoryManager{
	AsymmetricKeys: map[string]*key.AsymmetricKey{},
	SymmetricKeys:  map[string]*key.SymmetricKey{},
	Strategy:       keyStrategy,
}

var hmacStrategy = &strategy.HMACSHAStrategy{
	Enigma: &hmac.HMACStrategy{
		GlobalSecret: []byte("some-super-cool-secret-that-nobody-knows"),
	},
}

var authCodeHandler = &explicit.AuthorizeExplicitGrantTypeHandler{
	AccessTokenStrategy:       hmacStrategy,
	RefreshTokenStrategy:      hmacStrategy,
	AuthorizeCodeStrategy:     hmacStrategy,
	AuthorizeCodeGrantStorage: store,
	AuthCodeLifespan:          time.Hour,
	AccessTokenLifespan:       time.Hour,
}

var hasher = &hash.BCrypt{}

var handler = &Handler{
	OAuth2: &fosite.Fosite{
		Store:          store,
		MandatoryScope: "hydra",
		AuthorizeEndpointHandlers: fosite.AuthorizeEndpointHandlers{
			authCodeHandler,
		},
		TokenEndpointHandlers: fosite.TokenEndpointHandlers{
			authCodeHandler,
			&client.ClientCredentialsGrantHandler{
				HandleHelper: &core.HandleHelper{
					AccessTokenStrategy: hmacStrategy,
					AccessTokenStorage:  store,
					AccessTokenLifespan: time.Hour,
				},
			},
		},
		AuthorizedRequestValidators: fosite.AuthorizedRequestValidators{},
		Hasher: hasher,
	},
	Consent: &DefaultConsentStrategy{
		Issuer:     "https://hydra.localhost",
		KeyManager: keyManager,
	},
}

var router = httprouter.New()

var ts *httptest.Server

var oauthConfig *oauth2.Config

var oauthClientConfig *clientcredentials.Config

func init() {
	keyManager.CreateAsymmetricKey(ConsentChallengeKey)
	keyManager.CreateAsymmetricKey(ConsentEndpointKey)
	ts = httptest.NewServer(router)

	handler.SetRoutes(router)
	store.Clients["app"] = &fosite.DefaultClient{
		ID:           "app",
		Secret:       []byte("secret"),
		RedirectURIs: []string{ts.URL + "/callback"},
	}

	c, _ := url.Parse(ts.URL + "/consent")
	handler.ConsentURL = *c

	h, _ := hasher.Hash([]byte("secret"))
	store.Clients["app-client"] = &fosite.DefaultClient{
		ID:           "app-client",
		Secret:       h,
		RedirectURIs: []string{ts.URL + "/callback"},
	}

	oauthConfig = &oauth2.Config{
		ClientID:     "app-client",
		ClientSecret: "secret",
		Endpoint: oauth2.Endpoint{
			AuthURL:  ts.URL + "/oauth2/auth",
			TokenURL: ts.URL + "/oauth2/token",
		},
		RedirectURL: ts.URL + "/callback",
		Scopes:      []string{"hydra"},
	}

	oauthClientConfig = &clientcredentials.Config{
		ClientID:     "app-client",
		ClientSecret: "secret",
		TokenURL:     ts.URL + "/oauth2/token",
		Scopes:       []string{"hydra"},
	}
}

func signConsentToken(claims map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = claims

	key, err := keyManager.GetAsymmetricKey(ConsentEndpointKey)
	if err != nil {
		return "", errors.New(err)
	}

	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(key.Private)
	if err != nil {
		return "", errors.New(err)
	}

	var signature, encoded string
	if encoded, err = token.SigningString(); err != nil {
		return "", errors.New(err)
	} else if signature, err = token.Method.Sign(encoded, rsaKey); err != nil {
		return "", errors.New(err)
	}

	return fmt.Sprintf("%s.%s", encoded, signature), nil
}
