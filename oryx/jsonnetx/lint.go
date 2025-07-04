// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonnetx

import (
	"fmt"
	"os"

	"github.com/bmatcuk/doublestar/v2"
	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/linter"
	"github.com/spf13/cobra"

	"github.com/ory/x/cmdx"
)

// LintCommand represents the lint command
// Deprecated: use NewLintCommand instead.
var LintCommand = NewLintCommand()

func NewLintCommand() *cobra.Command {
	var verbose bool
	cmd := &cobra.Command{
		Use: "lint path/to/files/*.jsonnet [more/files.jsonnet, [supports/**/{foo,bar}.jsonnet]]",
		Long: `Lints JSONNet files using the official JSONNet linter and exits with a status code of 1 when issues are detected.

` + GlobHelp,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, pattern := range args {
				files, err := doublestar.Glob(pattern)
				if err != nil {
					_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Glob pattern %q is not valid: %s\n", pattern, err)
					return cmdx.FailSilently(cmd)
				}

				for _, file := range files {
					if fi, err := os.Stat(file); err != nil {
						_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Matching file %q could not be opened: %s\n", file, err)
						return cmdx.FailSilently(cmd)
					} else if fi.IsDir() {
						continue
					}

					if verbose {
						fmt.Printf("Processing file: %s\n", file)
					}

					//#nosec G304 -- false positive
					content, err := os.ReadFile(file)
					if err != nil {
						_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Unable to read file %q: %s\n", file, err)
						return cmdx.FailSilently(cmd)
					}

					if linter.LintSnippet(jsonnet.MakeVM(), cmd.ErrOrStderr(), []linter.Snippet{{FileName: file, Code: string(content)}}) {
						return cmdx.FailSilently(cmd)
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output.")
	return cmd
}
