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
	"fmt"
	"net/http"
	"strings" //"encoding/json"

	"github.com/ory/hydra/config" //"github.com/ory/hydra/oauth2"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx" //"context"
	"github.com/spf13/cobra"
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
	cmdx.ExactArgs(cmd, args, 1)

	c := hydra.NewAdminApiWithBasePath(h.Config.GetClusterURLWithoutTailingSlashOrFail(cmd))
	c.Configuration = configureClient(cmd, c.Configuration)

	clientID := flagx.MustGetString(cmd, "client-id")
	clientSecret := flagx.MustGetString(cmd, "client-secret")
	if clientID != "" || clientSecret != "" {
		c.Configuration.Username = clientID
		c.Configuration.Password = clientSecret
	} else {
		fmt.Println("No OAuth 2.0 Client ID an secret set, skipping authorization header. This might fail if the introspection endpoint is protected.")
	}

	result, response, err := c.IntrospectOAuth2Token(args[0], strings.Join(flagx.MustGetStringSlice(cmd, "scope"), " "))
	checkResponse(err, http.StatusOK, response)
	fmt.Println(formatResponse(result))
}
