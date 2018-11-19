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
	"time"

	"github.com/spf13/cobra"

	"github.com/ory/hydra/config"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

type TokenHandler struct {
	Config *config.Config
}

func (h *TokenHandler) newTokenManager(cmd *cobra.Command) *hydra.AdminApi {
	c := hydra.NewAdminApiWithBasePath(h.Config.GetClusterURLWithoutTailingSlashOrFail(cmd))
	c.Configuration = configureClientWithoutAuth(cmd, c.Configuration)
	return c
}

func newTokenHandler(c *config.Config) *TokenHandler {
	return &TokenHandler{Config: c}
}

func (h *TokenHandler) RevokeToken(cmd *cobra.Command, args []string) {
	cmdx.ExactArgs(cmd, args, 1)

	handler := hydra.NewPublicApiWithBasePath(h.Config.GetClusterURLWithoutTailingSlashOrFail(cmd))
	handler.Configuration = configureClientWithoutAuth(cmd, handler.Configuration)

	if clientID, clientSecret := flagx.MustGetString(cmd, "client-id"), flagx.MustGetString(cmd, "client-secret"); clientID == "" || clientSecret == "" {
		cmdx.Fatalf(`%s

Please provide a Client ID and Client Secret using flags --client-id and --client-secret, or environment variables OAUTH2_CLIENT_ID and OAUTH2_CLIENT_SECRET
`, cmd.UsageString())
	} else {
		handler.Configuration.Username = clientID
		handler.Configuration.Password = clientSecret
	}

	token := args[0]
	response, err := handler.RevokeOAuth2Token(args[0])
	checkResponse(err, http.StatusOK, response)
	fmt.Printf("Revoked OAuth 2.0 Access Token: %s\n", token)
}

func (h *TokenHandler) FlushTokens(cmd *cobra.Command, args []string) {
	handler := hydra.NewAdminApiWithBasePath(h.Config.GetClusterURLWithoutTailingSlashOrFail(cmd))
	handler.Configuration = configureClient(cmd, handler.Configuration)
	response, err := handler.FlushInactiveOAuth2Tokens(hydra.FlushInactiveOAuth2TokensRequest{
		NotAfter: time.Now().Add(-flagx.MustGetDuration(cmd, "min-age")),
	})
	checkResponse(err, http.StatusNoContent, response)
	fmt.Println("Successfully flushed inactive access tokens")
}
