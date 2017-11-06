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

package cli

import (
	//"context"
	//"encoding/json"
	"fmt"

	"github.com/ory/hydra/config"
	//"github.com/ory/hydra/oauth2"

	"net/http"
	"strings"

	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
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

func (h *IntrospectionHandler) IsAuthorized(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Print(cmd.UsageString())
		return
	}

	c := hydra.NewOAuth2ApiWithBasePath(h.Config.ClusterURL)
	c.Configuration.Transport = h.Config.OAuth2Client(cmd).Transport

	if term, _ := cmd.Flags().GetBool("fake-tls-termination"); term {
		c.Configuration.DefaultHeader["X-Forwarded-Proto"] = "https"
	}

	scopes, _ := cmd.Flags().GetStringSlice("scopes")
	result, response, err := c.IntrospectOAuth2Token(args[0], strings.Join(scopes, " "))
	checkResponse(response, err, http.StatusOK)
	fmt.Printf("%s\n", formatResponse(result))
}
