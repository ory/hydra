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
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/hydra/internal/httpclient/models"
	"github.com/ory/hydra/x"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/pointerx"
)

type ClientHandler struct{}

func newClientHandler() *ClientHandler {
	return &ClientHandler{}
}

func (h *ClientHandler) ImportClients(cmd *cobra.Command, args []string) {
	cmdx.MinArgs(cmd, args, 1)
	m := configureClient(cmd)

	ek, encryptSecret, err := newEncryptionKey(cmd, nil)
	cmdx.Must(err, "Failed to load encryption key: %s", err)

	for _, path := range args {
		reader, err := os.Open(path)
		cmdx.Must(err, "Could not open file %s: %s", path, err)

		var c models.OAuth2Client
		err = json.NewDecoder(reader).Decode(&c)
		cmdx.Must(err, "Could not parse JSON from file %s: %s", path, err)

		response, err := m.Admin.CreateOAuth2Client(admin.NewCreateOAuth2ClientParams().WithBody(&c))
		cmdx.Must(err, "The request failed with the following error message:\n%s", formatSwaggerError(err))
		result := response.Payload

		if c.ClientSecret == "" {
			if encryptSecret {
				enc, err := ek.Encrypt([]byte(result.ClientSecret))
				if err == nil {
					fmt.Printf("Imported OAuth 2.0 Client %s from: %s\n", result.ClientID, path)
					fmt.Printf("OAuth 2.0 Encrypted Client Secret: %s\n\n", enc.Base64Encode())
					continue
				}

				fmt.Printf("Imported OAuth 2.0 Client %s:%s from: %s\n", result.ClientID, result.ClientSecret, path)
				cmdx.Must(err, "Failed to encrypt client secret: %s", err)
			}

			fmt.Printf("Imported OAuth 2.0 Client %s:%s from: %s\n", result.ClientID, result.ClientSecret, path)
		} else {
			fmt.Printf("Imported OAuth 2.0 Client %s from: %s\n", result.ClientID, path)
		}
	}
}

func (h *ClientHandler) CreateClient(cmd *cobra.Command, args []string) {
	var err error
	m := configureClient(cmd)
	secret := flagx.MustGetString(cmd, "secret")

	var echoSecret bool
	if secret == "" {
		var secretb []byte
		secretb, err = x.GenerateSecret(26)
		cmdx.Must(err, "Could not generate OAuth 2.0 Client Secret: %s", err)
		secret = string(secretb)

		echoSecret = true
	} else {
		fmt.Println("You should not provide secrets using command line flags, the secret might leak to bash history and similar systems")
	}

	ek, encryptSecret, err := newEncryptionKey(cmd, nil)
	cmdx.Must(err, "Failed to load encryption key: %s", err)

	cc := models.OAuth2Client{
		ClientID:                          flagx.MustGetString(cmd, "id"),
		ClientSecret:                      secret,
		ResponseTypes:                     flagx.MustGetStringSlice(cmd, "response-types"),
		Scope:                             strings.Join(flagx.MustGetStringSlice(cmd, "scope"), " "),
		GrantTypes:                        flagx.MustGetStringSlice(cmd, "grant-types"),
		RedirectUris:                      flagx.MustGetStringSlice(cmd, "callbacks"),
		ClientName:                        flagx.MustGetString(cmd, "name"),
		TokenEndpointAuthMethod:           flagx.MustGetString(cmd, "token-endpoint-auth-method"),
		JwksURI:                           flagx.MustGetString(cmd, "jwks-uri"),
		TosURI:                            flagx.MustGetString(cmd, "tos-uri"),
		PolicyURI:                         flagx.MustGetString(cmd, "policy-uri"),
		LogoURI:                           flagx.MustGetString(cmd, "logo-uri"),
		ClientURI:                         flagx.MustGetString(cmd, "client-uri"),
		AllowedCorsOrigins:                flagx.MustGetStringSlice(cmd, "allowed-cors-origins"),
		SubjectType:                       flagx.MustGetString(cmd, "subject-type"),
		Audience:                          flagx.MustGetStringSlice(cmd, "audience"),
		PostLogoutRedirectUris:            flagx.MustGetStringSlice(cmd, "post-logout-callbacks"),
		BackchannelLogoutSessionRequired:  flagx.MustGetBool(cmd, "backchannel-logout-session-required"),
		BackchannelLogoutURI:              flagx.MustGetString(cmd, "backchannel-logout-callback"),
		FrontchannelLogoutSessionRequired: flagx.MustGetBool(cmd, "frontchannel-logout-session-required"),
		FrontchannelLogoutURI:             flagx.MustGetString(cmd, "frontchannel-logout-callback"),
	}

	response, err := m.Admin.CreateOAuth2Client(admin.NewCreateOAuth2ClientParams().WithBody(&cc))
	cmdx.Must(err, "The request failed with the following error message:\n%s", formatSwaggerError(err))
	result := response.Payload

	fmt.Printf("OAuth 2.0 Client ID: %s\n", result.ClientID)
	if result.ClientSecret == "" {
		fmt.Println("This OAuth 2.0 Client has no secret")
	} else {
		if echoSecret {
			if encryptSecret {
				enc, err := ek.Encrypt([]byte(result.ClientSecret))
				if err == nil {
					fmt.Printf("OAuth 2.0 Encrypted Client Secret: %s\n", enc.Base64Encode())
					return
				}

				// executes this at last to print raw client secret
				// because if executes immediately, nobody knows client secret
				defer cmdx.Must(err, "Failed to encrypt client secret: %s", err)
			}

			fmt.Printf("OAuth 2.0 Client Secret: %s\n", result.ClientSecret)
		}
	}
}

func (h *ClientHandler) UpdateClient(cmd *cobra.Command, args []string) {
	cmdx.ExactArgs(cmd, args, 1)
	m := configureClient(cmd)
	newSecret := flagx.MustGetString(cmd, "secret")

	var echoSecret bool
	if newSecret != "" {
		echoSecret = true
		fmt.Println("You should not provide secrets using command line flags, the secret might leak to bash history and similar systems")
	}

	ek, encryptSecret, err := newEncryptionKey(cmd, nil)
	cmdx.Must(err, "Failed to load encryption key: %s", err)

	id := args[0]
	cc := models.OAuth2Client{
		ClientID:                          id,
		ClientSecret:                      newSecret,
		ResponseTypes:                     flagx.MustGetStringSlice(cmd, "response-types"),
		Scope:                             strings.Join(flagx.MustGetStringSlice(cmd, "scope"), " "),
		GrantTypes:                        flagx.MustGetStringSlice(cmd, "grant-types"),
		RedirectUris:                      flagx.MustGetStringSlice(cmd, "callbacks"),
		ClientName:                        flagx.MustGetString(cmd, "name"),
		TokenEndpointAuthMethod:           flagx.MustGetString(cmd, "token-endpoint-auth-method"),
		JwksURI:                           flagx.MustGetString(cmd, "jwks-uri"),
		TosURI:                            flagx.MustGetString(cmd, "tos-uri"),
		PolicyURI:                         flagx.MustGetString(cmd, "policy-uri"),
		LogoURI:                           flagx.MustGetString(cmd, "logo-uri"),
		ClientURI:                         flagx.MustGetString(cmd, "client-uri"),
		AllowedCorsOrigins:                flagx.MustGetStringSlice(cmd, "allowed-cors-origins"),
		SubjectType:                       flagx.MustGetString(cmd, "subject-type"),
		Audience:                          flagx.MustGetStringSlice(cmd, "audience"),
		PostLogoutRedirectUris:            flagx.MustGetStringSlice(cmd, "post-logout-callbacks"),
		BackchannelLogoutSessionRequired:  flagx.MustGetBool(cmd, "backchannel-logout-session-required"),
		BackchannelLogoutURI:              flagx.MustGetString(cmd, "backchannel-logout-callback"),
		FrontchannelLogoutSessionRequired: flagx.MustGetBool(cmd, "frontchannel-logout-session-required"),
		FrontchannelLogoutURI:             flagx.MustGetString(cmd, "frontchannel-logout-callback"),
	}

	response, err := m.Admin.UpdateOAuth2Client(admin.NewUpdateOAuth2ClientParams().WithID(id).WithBody(&cc))
	cmdx.Must(err, "The request failed with the following error message:\n%s", formatSwaggerError(err))
	result := response.Payload
	fmt.Printf("%s OAuth 2.0 Client updated\n", result.ClientID)

	if echoSecret {
		if encryptSecret {
			enc, err := ek.Encrypt([]byte(result.ClientSecret))
			if err == nil {
				fmt.Printf("OAuth 2.0 Encrypted Client Secret: %s\n", enc.Base64Encode())
				return
			}

			// executes this at last to print raw client secret
			// because if executes immediately, nobody knows client secret
			defer cmdx.Must(err, "Failed to encrypt client secret: %s", err)
		}
		fmt.Printf("Updated OAuth 2.0 Client Secret: %s\n", result.ClientSecret)
	}
}

func (h *ClientHandler) DeleteClient(cmd *cobra.Command, args []string) {
	cmdx.MinArgs(cmd, args, 1)
	m := configureClient(cmd)

	for _, c := range args {
		_, err := m.Admin.DeleteOAuth2Client(admin.NewDeleteOAuth2ClientParams().WithID(c))
		cmdx.Must(err, "The request failed with the following error message:\n%s", formatSwaggerError(err))
	}

	fmt.Println("OAuth 2.0 Client(s) deleted")
}

func (h *ClientHandler) GetClient(cmd *cobra.Command, args []string) {

	m := configureClient(cmd)
	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	response, err := m.Admin.GetOAuth2Client(admin.NewGetOAuth2ClientParams().WithID(args[0]))
	cmdx.Must(err, "The request failed with the following error message:\n%s", formatSwaggerError(err))
	cl := response.Payload
	fmt.Println(cmdx.FormatResponse(cl))
}

func (h *ClientHandler) ListClients(cmd *cobra.Command, args []string) {
	m := configureClient(cmd)

	limit := flagx.MustGetInt(cmd, "limit")
	page := flagx.MustGetInt(cmd, "page")
	offset := (limit * page) - limit

	response, err := m.Admin.ListOAuth2Clients(admin.NewListOAuth2ClientsParams().WithLimit(pointerx.Int64(int64(limit))).WithOffset(pointerx.Int64(int64(offset))))
	cmdx.Must(err, "The request failed with the following error message:\n%s", formatSwaggerError(err))
	cls := response.Payload

	table := newTable()
	table.SetHeader([]string{
		"Client ID",
		"Name",
		"Response Types",
		"Scope",
		"Redirect Uris",
		"Grant Types",
		"Token Endpoint Auth Method",
	})

	data := make([][]string, len(cls))
	for i, cl := range cls {
		data[i] = []string{
			cl.ClientID,
			cl.ClientName,
			strings.Join(cl.ResponseTypes, ","),
			cl.Scope,
			strings.Join(cl.RedirectUris, "\n"),
			strings.Join(cl.GrantTypes, ","),
			cl.TokenEndpointAuthMethod,
		}
	}

	table.AppendBulk(data)
	table.Render()
}
