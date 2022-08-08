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

	"github.com/spf13/cobra"

	"github.com/ory/hydra/cmd/cliclient"
	"github.com/ory/x/cmdx"
)

func NewDeleteClientCmd(root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "client id-1 [id-2] [id-n]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Delete one or more OAuth 2.0 Clients by their ID(s)",
		Long:  "This command deletes one or more OAuth 2.0 Clients by their respective IDs.",
		Example: fmt.Sprintf(`%[1]s delete client client-1 client-2 client-3

To delete OAuth 2.0 Clients with the owner of "foo@bar.com", run:

	%[1]s delete client $(%[1]s list clients --format json | jq -r 'map(select(.contacts[] == "foo@bar.com")) | .[].client_id')`, root.Use),
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			var (
				deleted = make([]cmdx.OutputIder, 0, len(args))
				failed  = make(map[string]error)
			)

			for _, c := range args {
				_, err := m.V0alpha2Api.AdminDeleteOAuth2Client(cmd.Context(), c).Execute() //nolint:bodyclose
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
