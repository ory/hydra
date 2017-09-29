package hydra

import (
	"github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/pkg/errors"
	"context"
	"golang.org/x/oauth2/clientcredentials"
)

// SDK contains all relevant API clients for interacting with ORY Hydra.
type SDK struct {
	*swagger.OAuth2Api
	*swagger.JsonWebKeyApi
	*swagger.WardenApi
	*swagger.PolicyApi

	Configuration *Configuration
}

// Configuration configures the SDK.
type Configuration struct {
	// EndpointURL should point to the url of ORY Hydra, for example: http://localhost:4444
	EndpointURL     string

	// ClientID is the id of the management client. The management client should have appropriate access rights
	// and the ability to request the client_credentials grant.
	ClientID     string

	// ClientSecret is the secret of the management client.
	ClientSecret string

	// Scopes is a list of scopes the SDK should request. If no scopes are given, this defaults to `hydra.*`
	Scopes       []string
}

func removeTrailingSlash(path string) string {
	for len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0: len(path)-1]
	}
	return path
}

// NewSDK instantiates a new SDK instance or returns an error.
func NewSDK(c *Configuration) (*SDK, error) {
	if c.EndpointURL == "" {
		return nil, errors.New("Please specify an EndpointURL url")
	}
	if c.ClientSecret == "" {
		return nil, errors.New("Please specify a client id")
	}
	if c.ClientID == "" {
		return nil, errors.New("Please specify a client secret")
	}
	if len(c.Scopes) == 0 {
		c.Scopes = []string{"hydra.*"}
	}

	c.EndpointURL = removeTrailingSlash(c.EndpointURL)

	oAuth2Config := clientcredentials.Config{
		ClientSecret: c.ClientSecret,
		ClientID:     c.ClientID,
		Scopes:       c.Scopes,
		TokenURL:     c.EndpointURL + "/oauth2/token",
	}
	oAuth2Client := oAuth2Config.Client(context.Background())

	o := swagger.NewOAuth2ApiWithBasePath(c.EndpointURL)
	o.Configuration.Transport = oAuth2Client.Transport
	o.Configuration.Username = c.ClientID
	o.Configuration.Password = c.ClientSecret

	j := swagger.NewJsonWebKeyApiWithBasePath(c.EndpointURL)
	j.Configuration.Transport = oAuth2Client.Transport

	w := swagger.NewWardenApiWithBasePath(c.EndpointURL)
	w.Configuration.Transport = oAuth2Client.Transport

	p := swagger.NewPolicyApiWithBasePath(c.EndpointURL)
	p.Configuration.Transport = oAuth2Client.Transport

	sdk := &SDK{
		OAuth2Api:     o,
		JsonWebKeyApi: j,
		WardenApi:     w,
		PolicyApi:     p,
		Configuration: c,
	}
	
	return sdk, nil
}
