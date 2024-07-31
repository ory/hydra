// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"strings"

	"github.com/ory/hydra/v2/cmd/cliclient"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"

	"github.com/spf13/cobra"
)

func NewIntrospectTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token the-token",
		Args:    cobra.ExactArgs(1),
		Example: `{{ .CommandPath }} AYjcyMzY3ZDhiNmJkNTY --project 32197be3-8e57-4009-becd-9d38dbde129c`,
		Short:   "Introspect an OAuth 2.0 Access or Refresh Token",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			result, _, err := client.OAuth2API.IntrospectOAuth2Token(cmd.Context()).
				Token(args[0]).
				Scope(strings.Join(flagx.MustGetStringSlice(cmd, "scope"), " ")).Execute() //nolint:bodyclose
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			cmdx.PrintRow(cmd, outputOAuth2TokenIntrospection(*result))
			return nil
		},
	}
	cmd.Flags().StringSlice("scope", []string{}, "Additionally check if the scope was granted.")
	return cmd
}
