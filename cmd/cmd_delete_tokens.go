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

func NewDeleteAccessTokensCmd(parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "access-tokens client-id",
		Args:    cobra.ExactArgs(1),
		Example: fmt.Sprintf(`%s delete access-tokens 33137249-dd2c-49e6-a066-75ad2a72f221`, parent.Use),
		Short:   "Invalidate all OAuth2 Access Tokens of an OAuth2 Client",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			clientID := args[0]
			_, err = client.V0alpha2Api.AdminDeleteOAuth2Token(cmd.Context()).ClientId(clientID).Execute() //nolint:bodyclose
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			cmdx.PrintRow(cmd, cmdx.OutputIder(clientID))
			return nil
		},
	}
	return cmd
}
