// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/cmd/cliclient"
	"github.com/ory/x/cmdx"
)

func NewGetClientsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "oauth2-client <id-1> [<id-2> ...]",
		Aliases: []string{"client", "clients", "oauth2-clients"},
		Args:    cobra.MinimumNArgs(1),
		Short:   "Get one or more OAuth 2.0 Clients by their ID(s)",
		Long:    `This command gets all the details about an OAuth 2.0 Client. You can use this command in combination with jq.`,
		Example: `To get the OAuth 2.0 Client's name, run:

	{{ .CommandPath }} <your-client-id> --format json | jq -r '.client_name'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, _, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			clients := make([]hydra.OAuth2Client, 0, len(args))
			for _, id := range args {
				client, _, err := m.OAuth2API.GetOAuth2Client(cmd.Context(), id).Execute() //nolint:bodyclose
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
