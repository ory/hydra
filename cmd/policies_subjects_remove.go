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

// policiesSubjectsRemoveCmd represents the remove command
var policiesSubjectsRemoveCmd = &cobra.Command{
	Use:   "remove <policy> <subject> [<subject>...]",
	Short: "Remove subjects from the regex matching list",
	Long: `You can use regular expressions in your matches. Encapsulate them in < >.

Example:
  hydra policies subjects remove my-policy john@org.com <[peter|max]>@org.com`,
	Run: cmdHandler.Policies.RemoveSubjectFromPolicy,
}

func init() {
	policiesSubjectsCmd.AddCommand(policiesSubjectsRemoveCmd)
}
