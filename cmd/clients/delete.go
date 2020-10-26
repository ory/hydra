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

package clients

import (
	"errors"
	"fmt"
	"github.com/ory/hydra/cmd/cli"
	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"
)

// newDeleteCmd returns the delete command
func newDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <id> [<id>...]",
		Args:    cobra.MinimumNArgs(1),
		Short:   "Delete OAuth 2.0 Clients",
		Long:    "This command deletes one or more OAuth 2.0 Clients by their respective IDs.",
		Example: "$ hydra clients delete client-1 client-2 client-3",
		RunE: func(cmd *cobra.Command, args []string) error {
			m := cli.ConfigureClient(cmd)

			errs := make(map[string]error)
			for _, c := range args {
				_, err := m.Admin.DeleteOAuth2Client(admin.NewDeleteOAuth2ClientParams().WithID(c))
				if err != nil {
					errs[c] = errors.New(cli.FormatSwaggerError(err))
					continue
				}

				_, _ = fmt.Fprintln(cmd.OutOrStdout(), c)
			}

			if len(errs) != 0 {
				cmdx.PrintErrors(cmd, errs)
				return cmdx.FailSilently(cmd)
			}
			return nil
		},
	}

	cmd.Flags().AddFlagSet(packageFlags)

	return cmd
}
