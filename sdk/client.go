package sdk

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"

	"github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/jwk"
	hoauth2 "github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/hydra/policy"
	"github.com/ory-am/hydra/warden"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type option func(*Client) error

// default options for hydra client
var defaultOptions = []option{
	ClusterURL(os.Getenv("HYDRA_CLUSTER_URL")),
	ClientID(os.Getenv("HYDRA_CLIENT_ID")),
	ClientSecret(os.Getenv("HYDRA_CLIENT_SECRET")),
	Scopes("hydra"),
}

// Client offers easy use of all HTTP clients.
type Client struct {
	// Clients offers OAuth2 Client management capabilities.
	Clients *client.HTTPManager

	// JSONWebKeys offers JSON Web Key management capabilities.
	JSONWebKeys *jwk.HTTPManager

	// Policies offers Access Policy management capabilities.
	Policies *policy.HTTPManager

	// Warden offers Access Token and Access Request validation strategies (for first-party resource servers).
	Warden *warden.HTTPWarden

	// Introspection offers Access Token and Access Request introspection strategies (according to RFC 7662).
	Introspection *hoauth2.HTTPIntrospector

	// Revocation offers OAuth2 Token Revocation.
	Revocator *hoauth2.HTTPRecovator

	http          *http.Client
	clusterURL    *url.URL
	clientID      string
	clientSecret  string
	skipTLSVerify bool
	scopes        []string
	credentials   clientcredentials.Config
}

// Connect instantiates a new client to communicate with Hydra.
//
//  import "github.com/ory-am/hydra/sdk"
//
//  var hydra, err = sdk.Connect(
// 	sdk.ClientID("client-id"),
// 	sdk.ClientSecret("client-secret"),
//  	sdk.ClusterURL("https://localhost:4444"),
//  )
func Connect(opts ...option) (*Client, error) {
	c := &Client{}

	var err error
	// apply default options
	for _, opt := range defaultOptions {
		err = opt(c)
		if err != nil {
			return nil, err
		}
	}

	// override any default values with given options
	for _, opt := range opts {
		err = opt(c)
		if err != nil {
			return nil, err
		}
	}

	c.credentials = clientcredentials.Config{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		TokenURL:     pkg.JoinURL(c.clusterURL, "oauth2/token").String(),
		Scopes:       c.scopes,
	}

	c.http = http.DefaultClient

	if c.skipTLSVerify {
		c.http = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	err = c.authenticate()
	if err != nil {
		return nil, err
	}

	// initialize service endpoints
	c.Clients = &client.HTTPManager{
		Endpoint: pkg.JoinURL(c.clusterURL, "/clients"),
		Client:   c.http,
	}

	c.Revocator = &hoauth2.HTTPRecovator{
		Endpoint: pkg.JoinURL(c.clusterURL, hoauth2.RevocationPath),
		Config:   &c.credentials,
	}

	c.Introspection = &hoauth2.HTTPIntrospector{
		Endpoint: pkg.JoinURL(c.clusterURL, hoauth2.IntrospectPath),
		Client:   c.http,
	}

	c.JSONWebKeys = &jwk.HTTPManager{
		Endpoint: pkg.JoinURL(c.clusterURL, "/keys"),
		Client:   c.http,
	}

	c.Policies = &policy.HTTPManager{
		Endpoint: pkg.JoinURL(c.clusterURL, "/policies"),
		Client:   c.http,
	}

	c.Warden = &warden.HTTPWarden{
		Client:   c.http,
		Endpoint: c.clusterURL,
	}

	return c, nil
}

// OAuth2Config returns an oauth2 config instance which you can use to initiate various oauth2 flows.
//
//  config := client.OAuth2Config("https://mydomain.com/oauth2_callback", "photos", "contacts.read")
//  redirectRequestTo := oauth2.AuthCodeURL()
//
//  // in callback handler...
//  token, err := config.Exchange(oauth2.NoContext, authorizeCode)
func (h *Client) OAuth2Config(redirectURL string, scopes ...string) *oauth2.Config {
	return &oauth2.Config{
		ClientSecret: h.clientSecret,
		ClientID:     h.clientID,
		Endpoint: oauth2.Endpoint{
			TokenURL: pkg.JoinURL(h.clusterURL, "/oauth2/token").String(),
			AuthURL:  pkg.JoinURL(h.clusterURL, "/oauth2/auth").String(),
		},
		Scopes:      scopes,
		RedirectURL: redirectURL,
	}
}

func (h *Client) authenticate() error {
	ctx := context.WithValue(oauth2.NoContext, oauth2.HTTPClient, h.http)
	_, err := h.credentials.Token(ctx)
	if err != nil {
		return err
	}

	h.http = h.credentials.Client(ctx)
	return nil
}
