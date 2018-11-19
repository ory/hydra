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
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ory/hydra/config"
	"github.com/ory/hydra/pkg"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

type ClientHandler struct {
	Config *config.Config
}

func newClientHandler(c *config.Config) *ClientHandler {
	return &ClientHandler{
		Config: c,
	}
}

func (h *ClientHandler) newClientManager(cmd *cobra.Command) *hydra.AdminApi {
	c := hydra.NewAdminApiWithBasePath(h.Config.GetClusterURLWithoutTailingSlashOrFail(cmd))
	c.Configuration = configureClient(cmd, c.Configuration)
	return c
}

func (h *ClientHandler) ImportClients(cmd *cobra.Command, args []string) {
	cmdx.MinArgs(cmd, args, 1)
	m := h.newClientManager(cmd)

	for _, path := range args {
		reader, err := os.Open(path)
		cmdx.Must(err, "Could not open file %s: %s", path, err)

		var c hydra.OAuth2Client
		err = json.NewDecoder(reader).Decode(&c)
		cmdx.Must(err, "Could not parse JSON from file %s: %s", path, err)

		result, response, err := m.CreateOAuth2Client(c)
		checkResponse(err, http.StatusCreated, response)

		if c.ClientSecret == "" {
			fmt.Printf("Imported OAuth 2.0 Client %s:%s from: %s\n", result.ClientId, result.ClientSecret, path)
		} else {
			fmt.Printf("Imported OAuth 2.0 Client %s from: %s\n", result.ClientId, path)
		}
	}
}

func (h *ClientHandler) CreateClient(cmd *cobra.Command, args []string) {
	var err error
	m := h.newClientManager(cmd)
	secret := flagx.MustGetString(cmd, "secret")

	var echoSecret bool
	if secret == "" {
		var secretb []byte
		secretb, err = pkg.GenerateSecret(26)
		cmdx.Must(err, "Could not generate OAuth 2.0 Client Secret: %s", err)
		secret = string(secretb)

		echoSecret = true
	} else {
		fmt.Println("You should not provide secrets using command line flags, the secret might leak to bash history and similar systems")
	}

	cc := hydra.OAuth2Client{
		ClientId:                flagx.MustGetString(cmd, "id"),
		ClientSecret:            secret,
		ResponseTypes:           flagx.MustGetStringSlice(cmd, "response-types"),
		Scope:                   strings.Join(flagx.MustGetStringSlice(cmd, "scope"), " "),
		GrantTypes:              flagx.MustGetStringSlice(cmd, "grant-types"),
		RedirectUris:            flagx.MustGetStringSlice(cmd, "callbacks"),
		ClientName:              flagx.MustGetString(cmd, "name"),
		TokenEndpointAuthMethod: flagx.MustGetString(cmd, "token-endpoint-auth-method"),
		JwksUri:                 flagx.MustGetString(cmd, "jwks-uri"),
		TosUri:                  flagx.MustGetString(cmd, "tos-uri"),
		PolicyUri:               flagx.MustGetString(cmd, "policy-uri"),
		LogoUri:                 flagx.MustGetString(cmd, "logo-uri"),
		ClientUri:               flagx.MustGetString(cmd, "client-uri"),
		SubjectType:             flagx.MustGetString(cmd, "subject-type"),
		Audience:                flagx.MustGetStringSlice(cmd, "audience"),
	}

	result, response, err := m.CreateOAuth2Client(cc)
	checkResponse(err, http.StatusCreated, response)

	fmt.Printf("OAuth 2.0 Client ID: %s\n", result.ClientId)
	if result.ClientSecret == "" {
		fmt.Println("This OAuth 2.0 Client has no secret")
	} else {
		if echoSecret {
			fmt.Printf("OAuth 2.0 Client Secret: %s\n", result.ClientSecret)
		}
	}
}

func (h *ClientHandler) DeleteClient(cmd *cobra.Command, args []string) {
	cmdx.MinArgs(cmd, args, 1)
	m := h.newClientManager(cmd)

	for _, c := range args {
		response, err := m.DeleteOAuth2Client(c)
		checkResponse(err, http.StatusNoContent, response)
	}

	fmt.Println("OAuth2 client(s) deleted")
}

func (h *ClientHandler) GetClient(cmd *cobra.Command, args []string) {
	m := h.newClientManager(cmd)

	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	cl, response, err := m.GetOAuth2Client(args[0])
	checkResponse(err, http.StatusOK, response)
	fmt.Println(cmdx.FormatResponse(&cl))
}
