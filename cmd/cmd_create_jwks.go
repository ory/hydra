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
	"context"

	"github.com/spf13/cobra"

	hydra "github.com/ory/hydra-client-go"
	"github.com/ory/hydra/cmd/cliclient"
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
			jwks, _, err := m.JwkApi.CreateJsonWebKeySet(context.Background(), args[0]).CreateJsonWebKeySet(hydra.CreateJsonWebKeySet{
				Alg: flagx.MustGetString(cmd, alg),
				Kid: kid,
				Use: flagx.MustGetString(cmd, use),
			}).Execute()
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			cmdx.PrintTable(cmd, &outputJSONWebKeyCollection{Keys: jwks.Keys, Set: args[0]})
			return nil
		},
	}
	cmd.Root().Name()
	cmd.Flags().String(alg, "RS256", "The algorithm to be used to generated they key. Supports: RS256, RS512, ES256, ES512, EdDSA")
	cmd.Flags().String(use, "sig", "The intended use of this key. Supports: sig, enc")
	return cmd
}
