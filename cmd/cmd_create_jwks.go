// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"

	"github.com/ory/hydra/v2/jwk"

	"github.com/spf13/cobra"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/cmd/cliclient"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

func NewCreateJWKSCmd() *cobra.Command {
	const alg = "alg"
	const use = "use"

	cmd := &cobra.Command{
		Use:     "jwk <set-id> [<key-id>]",
		Aliases: []string{"jwks"},
		Args:    cobra.RangeArgs(1, 2),
		Example: `{{ .CommandPath }} <my-jwk-set> --alg RS256 --use sig`,
		Short:   "Create a JSON Web Key Set with a JSON Web Key",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.CommandPath()
			m, _, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			var kid string
			if len(args) == 2 {
				kid = args[1]
			}

			//nolint:bodyclose
			jwks, _, err := m.JwkAPI.CreateJsonWebKeySet(context.Background(), args[0]).CreateJsonWebKeySet(hydra.CreateJsonWebKeySet{
				Alg: flagx.MustGetString(cmd, alg),
				Kid: kid,
				Use: flagx.MustGetString(cmd, use),
			}).Execute()
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			if flagx.MustGetBool(cmd, "public") {
				jwks.Keys, err = jwk.OnlyPublicSDKKeys(jwks.Keys)
				if err != nil {
					return err
				}
			}

			cmdx.PrintTable(cmd, &outputJSONWebKeyCollection{Keys: jwks.Keys, Set: args[0]})
			return nil
		},
	}

	cmd.Flags().String(alg, "RS256", "The algorithm to be used to generated they key. Supports: RS256, RS512, ES256, ES512, EdDSA")
	cmd.Flags().String(use, "sig", "The intended use of this key. Supports: sig, enc")
	cmd.Flags().Bool("public", false, "Only return public keys")
	return cmd
}
