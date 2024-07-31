// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/hydra/v2/cmd/cliclient"
	"github.com/ory/x/cmdx"
)

func NewDeleteClientCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "oauth2-client <id-1> [<id-2> ...]",
		Aliases: []string{"client", "clients", "oauth2-clients"},
		Args:    cobra.MinimumNArgs(1),
		Short:   "Delete one or more OAuth 2.0 Clients by their ID(s)",
		Long:    "This command deletes one or more OAuth 2.0 Clients by their respective IDs.",
		Example: `{{ .CommandPath }} <client-1> <client-2> <client-3>

To delete OAuth 2.0 Clients with the owner of "foo@bar.com", run:

	{{ .CommandPath }} $({{ .Root.Name }} list oauth2-clients --format json | jq -r 'map(select(.contacts[] == "foo@bar.com")) | .[].client_id')`,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, _, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			var (
				deleted = make([]cmdx.OutputIder, 0, len(args))
				failed  = make(map[string]error)
			)

			for _, c := range args {
				_, err := m.OAuth2API.DeleteOAuth2Client(cmd.Context(), c).Execute() //nolint:bodyclose
				if err != nil {
					failed[c] = cmdx.PrintOpenAPIError(cmd, err)
					continue
				}
				deleted = append(deleted, cmdx.OutputIder(c))
			}

			if len(deleted) == 1 {
				cmdx.PrintRow(cmd, &deleted[0])
			} else if len(deleted) > 1 {
				cmdx.PrintTable(cmd, &cmdx.OutputIderCollection{Items: deleted})
			}

			cmdx.PrintErrors(cmd, failed)
			if len(failed) != 0 {
				return cmdx.FailSilently(cmd)
			}

			return nil
		},
	}
}
