package sdk

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"

	"github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/connection"
	"github.com/ory-am/hydra/jwk"
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
	// Client offers OAuth2 Client management capabilities.
	Client   *client.HTTPManager

	// SSO offers Social Login management capabilities.
	SSO      *connection.HTTPManager

	// JWK offers JSON Web Key management capabilities.
	JWK      *jwk.HTTPManager

	// Policies offers Access Policy management capabilities.
	Policies *policy.HTTPManager

	// Warden offers Access Token and Access Request validation strategies.
	Warden   *warden.HTTPWarden

	http          *http.Client
	clusterURL    *url.URL
	clientID      string
	clientSecret  string
	skipTLSVerify bool
	scopes        []string
	credentials clientcredentials.Config
}

// Connect instantiates a new client to communicate with Hydra.
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
	c.Client = &client.HTTPManager{
		Endpoint: pkg.JoinURL(c.clusterURL, "/clients"),
		Client:   c.http,
	}

	c.SSO = &connection.HTTPManager{
		Endpoint: pkg.JoinURL(c.clusterURL, "/connections"),
		Client:   c.http,
	}

	c.JWK = &jwk.HTTPManager{
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

func (h *Client) authenticate() error {
	ctx := context.WithValue(oauth2.NoContext, oauth2.HTTPClient, h.http)
	_, err := h.credentials.Token(ctx)
	if err != nil {
		return err
	}

	h.http = h.credentials.Client(ctx)
	return nil
}
