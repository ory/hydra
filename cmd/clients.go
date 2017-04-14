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
	"github.com/spf13/cobra"
)

// clientsCmd represents the clients command
var clientsCmd = &cobra.Command{
	Use:   "clients <command>",
	Short: "Manage OAuth2 clients",
	Long:  `Use this command to create, modify or delete OAuth2 clients.`,
}

func init() {
	RootCmd.AddCommand(clientsCmd)
	clientsCmd.PersistentFlags().Bool("dry", false, "do not execute the command but show the corresponding curl command instead")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
