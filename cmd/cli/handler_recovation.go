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
	"fmt"
	"net/http"

	"crypto/tls"

	"github.com/ory/hydra/config"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/spf13/cobra"
)

type RevocationHandler struct {
	Config *config.Config
}

func newRevocationHandler(c *config.Config) *RevocationHandler {
	return &RevocationHandler{Config: c}
}

func (h *RevocationHandler) RevokeToken(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Print(cmd.UsageString())
		return
	}

	handler := hydra.NewOAuth2ApiWithBasePath(h.Config.ClusterURL)
	handler.Configuration.Username = h.Config.ClientID
	handler.Configuration.Password = h.Config.ClientSecret

	if skip, _ := cmd.Flags().GetBool("skip-tls-verify"); skip {
		handler.Configuration.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	if term, _ := cmd.Flags().GetBool("fake-tls-termination"); term {
		handler.Configuration.DefaultHeader["X-Forwarded-Proto"] = "https"
	}

	token := args[0]
	response, err := handler.RevokeOAuth2Token(args[0])
	checkResponse(response, err, http.StatusOK)
	fmt.Printf("Revoked token %s", token)
}
