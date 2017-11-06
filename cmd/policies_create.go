// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
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

// policiesCreateCmd represents the create command
var policiesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new policy",
	Long: `To create a policy, either specify the files flag or pass arguments to create it directly from the CLI.

Example
  hydra policies create -f policy-a.json,policy-b.json
  hydra policies create -s peter,max -r blog,users -a post,ban --allow`,
	Run: cmdHandler.Policies.CreatePolicy,
}

func init() {
	policiesCmd.AddCommand(policiesCreateCmd)

	policiesCreateCmd.Flags().StringSliceP("files", "f", []string{}, "A list of paths to JSON encoded policy files")
	policiesCreateCmd.Flags().StringP("id", "i", "", "The policy's id")
	policiesCreateCmd.Flags().StringP("description", "d", "", "The policy's description")
	policiesCreateCmd.Flags().StringSliceP("resources", "r", []string{}, "A list of resource regex strings this policy will match to (required)")
	policiesCreateCmd.Flags().StringSliceP("subjects", "s", []string{}, "A list of subject regex strings this policy will match to (required)")
	policiesCreateCmd.Flags().StringSliceP("actions", "a", []string{}, "A list of action regex strings this policy will match to (required)")
	policiesCreateCmd.Flags().Bool("allow", false, "A list of action regex strings this policy will match to")
}
