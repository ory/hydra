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
	"os"

	"github.com/spf13/cobra"

	"github.com/ory/x/viperx"

	"github.com/ory/hydra/cmd/cli"
)

var (
	Version = "master"
	Date    = "undefined"
	Commit  = "undefined"
)

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hydra",
	Short: "Hydra is a cloud native high throughput OAuth2 and OpenID Connect provider",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

var cmdHandler = cli.NewHandler()

// Execute adds all child commands to the root command sets flags appropriately.
// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		viperx.InitializeConfig("hydra", "", nil)
	})

	viperx.RegisterConfigFlag(RootCmd, "hydra")

	RootCmd.PersistentFlags().Bool("skip-tls-verify", false, "Foolishly accept TLS certificates signed by unkown certificate authorities")

}
