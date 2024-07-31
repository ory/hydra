// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/ory/hydra/v2/cmd/cliclient"
	"github.com/ory/x/cmdx"
)

func NewDeleteJWKSCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "jwk <id-1> [<id-2> ...]",
		Aliases: []string{"jwks"},
		Args:    cobra.MinimumNArgs(1),
		Short:   "Delete one or more JSON Web Key Sets by their set ID",
		Long:    "This command deletes one or more JSON Web Key Sets by their respective set IDs.",
		Example: `{{ .CommandPath }} <set-1> <set-2> <set-3>`,
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
				_, err = m.JwkAPI.DeleteJsonWebKeySet(context.Background(), c).Execute() //nolint:bodyclose
				if err != nil {
					return cmdx.PrintOpenAPIError(cmd, err)
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
