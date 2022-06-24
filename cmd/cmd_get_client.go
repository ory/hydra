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

package cmd

import (
	"fmt"
	hydra "github.com/ory/hydra-client-go"
	"github.com/ory/hydra/cmd/cliclient"
	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"
)

func NewGetClientsCmd(root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "client [id-1] [id-2] [id-n]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Get one or more OAuth 2.0 Clients by their ID(s)",
		Long:  fmt.Sprintf(`This command gets all the details about an OAuth 2.0 Client. You can use this command in combination with jq.`),
		Example: fmt.Sprintf(`To get the OAuth 2.0 Client's secret, run:

	%s get client <your-client-id> | jq -r '.client_secret'`, root.Use),
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			clients := make([]hydra.OAuth2Client, 0, len(args))
			for _, id := range args {
				client, _, err := m.AdminApi.GetOAuth2Client(cmd.Context(), id).Execute()
				if err != nil {
					return cmdx.PrintOpenAPIError(cmd, err)
				}
				clients = append(clients, *client)
			}

			if len(clients) == 1 {
				cmdx.PrintRow(cmd, (*outputOAuth2Client)(&clients[0]))
			} else if len(clients) > 1 {
				cmdx.PrintTable(cmd, &outputOAuth2ClientCollection{clients})
			}

			return nil
		},
	}
}
