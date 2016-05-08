// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new OAuth2 client",
	Long: `This command creates a new OAuth2 client. Always specify at least one redirect url.

Please use the REST api to get access to all client fields.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(cmd.UsageString())
	},
}

func init() {
	clientsCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	createCmd.Flags().StringSliceP("callbacks", "c", []string{}, "REQUIRED list of allowed callback URLs")
	createCmd.Flags().StringSliceP("grant-types", "g", []string{"authorizeation_code"}, "A list of allowed grant types")
	createCmd.Flags().StringSliceP("response-types", "r", []string{"code"}, "A list of allowed response types")
	createCmd.Flags().StringP("name", "n", "", "The client's name")
}
