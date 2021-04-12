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
	"time"

	"github.com/spf13/cobra"
)

// tokenCmd represents the token command
func NewTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Issue and Manage OAuth2 tokens",
	}
	//tokenCmd.PersistentFlags().Bool("dry", false, "do not execute the command but show the corresponding curl command instead")
	cmd.PersistentFlags().Duration("fail-after", time.Minute, `Stop retrying after the specified duration`)
	cmd.PersistentFlags().Bool("fake-tls-termination", false, `fake tls termination by adding "X-Forwarded-Proto: https" to http headers`)
	cmd.PersistentFlags().Bool("skip-tls-verify", false, "Foolishly accept TLS certificates signed by unknown certificate authorities")
	return cmd
}
