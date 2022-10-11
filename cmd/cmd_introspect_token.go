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
	"strings"

	"github.com/ory/hydra/cmd/cliclient"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"

	"github.com/spf13/cobra"
)

func NewIntrospectTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token the-token",
		Args:    cobra.ExactArgs(1),
		Example: `{{ .CommandPath }} --client-id a0184d6c-b313-4e70-a0b9-905b581e9218 --client-secret Hh1BjioNNm ciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNT`,
		Short:   "Introspect an OAuth 2.0 Access or Refresh Token",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			result, _, err := client.OAuth2Api.IntrospectOAuth2Token(cmd.Context()).
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
