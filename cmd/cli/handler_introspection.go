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

package cli

import (
	"crypto/tls"
	//"context"
	//"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ory/hydra/config"
	//"github.com/ory/hydra/oauth2"

	"github.com/spf13/cobra"

	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
)

type IntrospectionHandler struct {
	Config *config.Config
}

func newIntrospectionHandler(c *config.Config) *IntrospectionHandler {
	return &IntrospectionHandler{
		Config: c,
	}
}

func (h *IntrospectionHandler) Introspect(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Print(cmd.UsageString())
		return
	}

	c := hydra.NewOAuth2ApiWithBasePath(h.Config.GetClusterURLWithoutTailingSlashOrFail(cmd))

	clientID, _ := cmd.Flags().GetString("client-id")
	clientSecret, _ := cmd.Flags().GetString("client-secret")
	if clientID != "" || clientSecret != "" {
		c.Configuration.Username = clientID
		c.Configuration.Password = clientSecret
	} else {
		fmt.Println("No OAuth 2.0 Client ID an secret set, skipping authorization header. This might fail if the introspection endpoint is protected.")
	}

	skipTLSTermination, _ := cmd.Flags().GetBool("skip-tls-verify")
	c.Configuration.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLSTermination},
	}

	if term, _ := cmd.Flags().GetBool("fake-tls-termination"); term {
		c.Configuration.DefaultHeader["X-Forwarded-Proto"] = "https"
	}

	scopes, _ := cmd.Flags().GetStringSlice("scope")
	result, response, err := c.IntrospectOAuth2Token(args[0], strings.Join(scopes, " "))
	checkResponse(response, err, http.StatusOK)
	fmt.Printf("%s\n", formatResponse(result))
}
