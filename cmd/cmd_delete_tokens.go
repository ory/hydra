// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/hydra/v2/cmd/cliclient"
	"github.com/ory/x/cmdx"
)

func NewDeleteAccessTokensCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "access-tokens <client-id>",
		Args:    cobra.ExactArgs(1),
		Example: `{{ .CommandPath }} <client-id>`,
		Short:   "Delete all OAuth2 Access Tokens of an OAuth2 Client",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			clientID := args[0]
			_, err = client.OAuth2API.DeleteOAuth2Token(cmd.Context()).ClientId(clientID).Execute() //nolint:bodyclose
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			cmdx.PrintRow(cmd, cmdx.OutputIder(clientID))
			return nil
		},
	}
	return cmd
}
