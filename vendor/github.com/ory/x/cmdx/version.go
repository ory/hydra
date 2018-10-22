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
 * @Copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 *
 */

package cmdx

import (
	"fmt"

	"os"

	"github.com/spf13/cobra"
)

func Version(gitTag, gitHash, buildTime *string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the build version, build time, and git hash",
		Run: func(cmd *cobra.Command, args []string) {
			if len(*gitTag) == 0 {
				fmt.Fprintln(os.Stderr, "Unable to determine version because the build process did not properly configure it.")
			} else {
				fmt.Printf("Version:			%s\n", *gitTag)
			}

			if len(*gitHash) == 0 {
				fmt.Fprintln(os.Stderr, "Unable to determine build commit because the build process did not properly configure it.")
			} else {
				fmt.Printf("Build Commit:	%s\n", *gitHash)
			}

			if len(*buildTime) == 0 {
				fmt.Fprintln(os.Stderr, "Unable to determine build timestamp because the build process did not properly configure it.")
			} else {
				fmt.Printf("Build Timestamp:	%s\n", *buildTime)
			}
		},
	}
}
