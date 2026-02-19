// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonnetx

import (
	"fmt"
	"os"

	"github.com/bmatcuk/doublestar/v2"
	"github.com/google/go-jsonnet/formatter"
	"github.com/spf13/cobra"

	"github.com/ory/x/cmdx"
)

func NewFormatCommand() *cobra.Command {
	var verbose, write bool
	cmd := &cobra.Command{
		Use: "format path/to/files/*.jsonnet [more/files.jsonnet, [supports/**/{foo,bar}.jsonnet]]",
		Long: `Formats JSONNet files using the official JSONNet formatter.

Use -w or --write to write output back to files instead of stdout.

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

					output, err := formatter.Format(file, string(content), formatter.DefaultOptions())
					if err != nil {
						_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "File %q could not be formatted: %s", file, err)
					}

					if write {
						err := os.WriteFile(file, []byte(output), 0o644) // #nosec
						if err != nil {
							_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Unable to write file %q: %s\n", file, err)
							return cmdx.FailSilently(cmd)
						}
					} else {
						fmt.Println(output)
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&write, "write", "w", false, "Write formatted output back to file.")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output.")
	return cmd
}
