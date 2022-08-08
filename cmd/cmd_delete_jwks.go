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
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ory/hydra/cmd/cliclient"
	"github.com/ory/x/cmdx"
)

func NewDeleteJWKSCommand(root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:     "jwks id-1 [id-2] [id-n]",
		Args:    cobra.MinimumNArgs(1),
		Short:   "Delete one or more JSON Web Key Sets by their set ID",
		Long:    "This command deletes one or more JSON Web Key Sets by their respective set IDs.",
		Example: fmt.Sprintf(`%[1]s delete jwks set-1 set-2 set-3`, root.Use),
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
				_, err = m.V0alpha2Api.AdminDeleteJsonWebKeySet(context.Background(), c).Execute() //nolint:bodyclose
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
