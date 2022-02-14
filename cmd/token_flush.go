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
	"os"
	"time"

	"github.com/ory/hydra/cmd/cli"

	"github.com/spf13/cobra"
)

// flushCmd represents the flush command
func NewTokenFlushCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flush",
		Short: "Removes inactive access tokens from the database",
		Run:   cli.NewHandler().Token.FlushTokens,
	}

	cmd.Flags().Duration("min-age", time.Duration(0), "Skip removing tokens which do not satisfy the minimum age (1s, 1m, 1h)")
	cmd.Flags().String("access-token", os.Getenv("OAUTH2_ACCESS_TOKEN"), "Set an access token to be used in the Authorization header, defaults to environment variable OAUTH2_ACCESS_TOKEN")
	cmd.Flags().String("endpoint", os.Getenv("HYDRA_ADMIN_URL"), "Set the URL where Ory Hydra is hosted, defaults to environment variable HYDRA_ADMIN_URL")
	return cmd
}
