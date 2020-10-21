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
	"fmt"
	"github.com/ory/hydra/cmd/cli"
	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "get <id>",
		Args:    cobra.ExactArgs(1),
		Short:   "Get an OAuth 2.0 Client",
		Long:    "This command retrieves an OAuth 2.0 Clients by its ID.",
		Example: "$ hydra clients get client-1",
		RunE: func(cmd *cobra.Command, args []string) error {
			m := cli.ConfigureClient(cmd)

			response, err := m.Admin.GetOAuth2Client(admin.NewGetOAuth2ClientParams().WithID(args[0]))
			if err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "The request failed with the following error message:\n%s", cli.FormatSwaggerError(err))
				return cmdx.FailSilently(cmd)
			}

			cmdx.PrintRow(cmd, (*outputOAuth2Client)(response.Payload))
			return nil
		},
	}
}
