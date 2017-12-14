// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hydra

import (
	"context"

	"github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// CodeGenSDK contains all relevant API clients for interacting with ORY Hydra.
type CodeGenSDK struct {
	*swagger.OAuth2Api
	*swagger.JsonWebKeyApi
	*swagger.WardenApi
	*swagger.PolicyApi

	Configuration      *Configuration
	oAuth2ClientConfig *clientcredentials.Config
	oAuth2Config       *oauth2.Config
}

// Configuration configures the CodeGenSDK.
type Configuration struct {
	// EndpointURL should point to the url of ORY Hydra, for example: http://localhost:4444
	EndpointURL string

	// ClientID is the id of the management client. The management client should have appropriate access rights
	// and the ability to request the client_credentials grant.
	ClientID string

	// ClientSecret is the secret of the management client.
	ClientSecret string

	// Scopes is a list of scopes the CodeGenSDK should request. If no scopes are given, this defaults to `hydra.*`
	Scopes []string
}

func removeTrailingSlash(path string) string {
	for len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	return path
}

func (s *CodeGenSDK) GetOAuth2ClientConfig() *clientcredentials.Config {
	return s.oAuth2ClientConfig
}
func (s *CodeGenSDK) GetOAuth2Config() *oauth2.Config {
	return s.oAuth2Config
}

// NewSDK instantiates a new CodeGenSDK instance or returns an error.
func NewSDK(c *Configuration) (*CodeGenSDK, error) {
	if c.EndpointURL == "" {
		return nil, errors.New("Please specify the ORY Hydra Endpoint URL")
	}
	if c.ClientSecret == "" {
		return nil, errors.New("Please specify an OAuth 2.0 Client Secret")
	}
	if c.ClientID == "" {
		return nil, errors.New("Please specify an OAuth 2.0 Client ID")
	}
	if len(c.Scopes) == 0 {
		c.Scopes = []string{"hydra.*"}
	}

	c.EndpointURL = removeTrailingSlash(c.EndpointURL)
	oAuth2Config := &oauth2.Config{
		ClientSecret: c.ClientSecret,
		ClientID:     c.ClientID,
		Scopes:       c.Scopes,
		Endpoint: oauth2.Endpoint{
			TokenURL: c.EndpointURL + "/oauth2/token",
			AuthURL:  c.EndpointURL + "/oauth2/auth",
		},
	}

	oAuth2ClientConfig := &clientcredentials.Config{
		ClientSecret: c.ClientSecret,
		ClientID:     c.ClientID,
		Scopes:       c.Scopes,
		TokenURL:     c.EndpointURL + "/oauth2/token",
	}
	oAuth2Client := oAuth2ClientConfig.Client(context.Background())

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

	sdk := &CodeGenSDK{
		OAuth2Api:          o,
		JsonWebKeyApi:      j,
		WardenApi:          w,
		PolicyApi:          p,
		Configuration:      c,
		oAuth2ClientConfig: oAuth2ClientConfig,
		oAuth2Config:       oAuth2Config,
	}

	return sdk, nil
}
