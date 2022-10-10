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
	"github.com/spf13/cobra"

	"github.com/ory/hydra/cmd/cliclient"
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
				_, err := m.OAuth2Api.DeleteOAuth2Client(cmd.Context(), c).Execute() //nolint:bodyclose
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
