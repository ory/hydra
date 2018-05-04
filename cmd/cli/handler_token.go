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
	"fmt"
	"net/http"
	"time"

	"github.com/ory/hydra/config"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/spf13/cobra"
)

type TokenHandler struct {
	Config *config.Config
}

func (h *TokenHandler) newTokenManager(cmd *cobra.Command) *hydra.OAuth2Api {
	c := hydra.NewOAuth2ApiWithBasePath(h.Config.GetClusterURLWithoutTailingSlash(cmd))

	skipTLSTermination, _ := cmd.Flags().GetBool("skip-tls-verify")
	c.Configuration.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLSTermination},
	}

	if term, _ := cmd.Flags().GetBool("fake-tls-termination"); term {
		c.Configuration.DefaultHeader["X-Forwarded-Proto"] = "https"
	}

	if token, _ := cmd.Flags().GetString("access-token"); token != "" {
		c.Configuration.DefaultHeader["Authorization"] = "Bearer " + token
	}

	return c
}

func newTokenHandler(c *config.Config) *TokenHandler {
	return &TokenHandler{Config: c}
}

func (h *TokenHandler) RevokeToken(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Print(cmd.UsageString())
		return
	}

	handler := hydra.NewOAuth2ApiWithBasePath(h.Config.GetClusterURLWithoutTailingSlash(cmd))

	skipTLSTermination, _ := cmd.Flags().GetBool("skip-tls-verify")
	handler.Configuration.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLSTermination},
	}

	clientID, _ := cmd.Flags().GetString("client-id")
	clientSecret, _ := cmd.Flags().GetString("client-secret")
	if clientID == "" || clientSecret == "" {
		fmt.Print(cmd.UsageString())
		fmt.Println("Please provide a Client ID and Client Secret using flags --client-id and --client-secret, or environment variables OAUTH2_CLIENT_ID and OAUTH2_CLIENT_SECRET.")
		return
	}

	handler.Configuration.Username = clientID
	handler.Configuration.Password = clientSecret

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

func (h *TokenHandler) FlushTokens(cmd *cobra.Command, args []string) {
	handler := h.newTokenManager(cmd)

	minAge, _ := cmd.Flags().GetDuration("min-age")

	response, err := handler.FlushInactiveOAuth2Tokens(hydra.FlushInactiveOAuth2TokensRequest{NotAfter: time.Now().Add(-minAge)})
	checkResponse(response, err, http.StatusNoContent)
	fmt.Println("Successfully flushed inactive access tokens.")
}
