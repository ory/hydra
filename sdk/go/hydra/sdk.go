/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package hydra

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/ory/hydra/sdk/go/hydra/swagger"
)

// CodeGenSDK contains all relevant API clients for interacting with ORY Hydra.
type CodeGenSDK struct {
	*swagger.AdminApi
	*swagger.PublicApi

	Configuration      *Configuration
	oAuth2ClientConfig *clientcredentials.Config
	oAuth2Config       *oauth2.Config
}

// Configuration configures the CodeGenSDK.
type Configuration struct {
	// ClientRegistrationPath should point to the administrative URL of ORY Hydra, for example: http://localhost:4445
	AdminURL string

	// PublicURL should point to the public url of ORY Hydra, for example: http://localhost:4444
	PublicURL string

	// ClientID is the id of the management client. The management client should have appropriate access rights
	// and the ability to request the client_credentials grant.
	ClientID string

	// ClientSecret is the secret of the management client.
	ClientSecret string

	// Scopes is a list of scopes the CodeGenSDK should request. If no scopes are given, this defaults to `hydra.*`
	Scopes []string
}

// NewSDK instantiates a new CodeGenSDK instance or returns an error.
func NewSDK(c *Configuration) (*CodeGenSDK, error) {
	if c.AdminURL == "" {
		return nil, errors.New("Please specify the ORY Hydra Admin URL")
	}

	c.AdminURL = strings.TrimLeft(c.AdminURL, "/")
	o := swagger.NewAdminApiWithBasePath(c.AdminURL)
	sdk := &CodeGenSDK{
		AdminApi:      o,
		Configuration: c,
	}

	if c.ClientSecret != "" && c.ClientID != "" && c.PublicURL != "" {
		if len(c.Scopes) == 0 {
			c.Scopes = []string{}
		}

		oAuth2ClientConfig := &clientcredentials.Config{
			ClientSecret: c.ClientSecret,
			ClientID:     c.ClientID,
			Scopes:       c.Scopes,
			TokenURL:     c.PublicURL + "/oauth2/token",
		}
		oAuth2Client := oAuth2ClientConfig.Client(context.Background())
		o.Configuration.Transport = oAuth2Client.Transport
		o.Configuration.Username = c.ClientID
		o.Configuration.Password = c.ClientSecret
		o.Configuration.Transport = oAuth2Client.Transport

		sdk.oAuth2ClientConfig = oAuth2ClientConfig
	} else if len(c.ClientSecret)+len(c.ClientID)+len(c.PublicURL) > 0 {
		return nil, errors.New("You provided one or more of client secret, ID or the public URL in the ORY Hydra SDK but not all of them. Please provide either none or all.")
	}

	return sdk, nil
}
