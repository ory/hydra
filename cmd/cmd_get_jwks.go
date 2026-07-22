// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/x/flagx"

	"github.com/ory/hydra/v2/cmd/cliclient"
	"github.com/ory/x/cmdx"
)

func NewGetJWKSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "jwk set-1 [set-2] ...",
		Aliases: []string{"jwks"},
		Args:    cobra.MinimumNArgs(1),
		Short:   "Get one or more JSON Web Key Set by its ID(s)",
		Long:    `This command gets all the details about an JSON Web Key. You can use this command in combination with jq.`,
		Example: `To get the JSON Web Key Set's use, run:

	{{ .CommandPath }} <set-id> | jq -r '.[].use'

To get the JSON Web Key Set as only public keys:

	{{ .CommandPath }} --public <set-id>'
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, _, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			var sets outputJSONWebKeyCollection
			for _, set := range args {
				key, _, err := m.JwkAPI.GetJsonWebKeySet(cmd.Context(), set).Execute() //nolint:bodyclose
				if err != nil {
					return cmdx.PrintOpenAPIError(cmd, err)
				}

				sets.Keys = append(sets.Keys, newOutputJsonWebKeys(set, key.Keys)...)
			}
			if len(args) == 1 {
				sets.Set = args[0]
			}

			if flagx.MustGetBool(cmd, "public") {
				keys := make([]hydra.JsonWebKey, len(sets.Keys))
				for i, key := range sets.Keys {
					keys[i] = key.JsonWebKey
				}
				keys, err = jwk.OnlyPublicSDKKeys(keys)
				if err != nil {
					return err
				}
				// OnlyPublicSDKKeys preserves order, so the set names still line up.
				for i, key := range keys {
					sets.Keys[i].JsonWebKey = key
				}
			}

			if len(sets.Keys) == 1 {
				cmdx.PrintRow(cmd, sets.Keys[0])
			} else if len(sets.Keys) > 1 {
				cmdx.PrintTable(cmd, sets)
			}

			return nil
		},
	}
	cmd.Flags().Bool("public", false, "Only return public keys")
	return cmd
}
