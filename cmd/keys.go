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

	"github.com/spf13/cobra"
)

// cmd represents the keys command
func NewKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage JSON Web Keys",
	}
	//cmd.PersistentFlags().Bool("dry", false, "do not execute the command but show the corresponding curl command instead")
	cmd.PersistentFlags().Bool("fake-tls-termination", false, `fake tls termination by adding "X-Forwarded-Proto: https" to http headers`)
	cmd.PersistentFlags().Duration("fail-after", time.Minute, `Stop retrying after the specified duration`)
	cmd.PersistentFlags().String("access-token", os.Getenv("OAUTH2_ACCESS_TOKEN"), "Set an access token to be used in the Authorization header, defaults to environment variable OAUTH2_ACCESS_TOKEN")
	cmd.PersistentFlags().String("endpoint", os.Getenv("HYDRA_ADMIN_URL"), "Set the URL where Ory Hydra is hosted, defaults to environment variable HYDRA_ADMIN_URL")
	cmd.PersistentFlags().Bool("skip-tls-verify", false, "Foolishly accept TLS certificates signed by unknown certificate authorities")
	return cmd
}
