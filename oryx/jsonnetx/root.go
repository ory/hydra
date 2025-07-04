// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonnetx

import (
	"github.com/spf13/cobra"
)

const GlobHelp = `Glob patterns supports the following special terms in the patterns:
	
	Special Terms | Meaning
	------------- | -------
	'*'           | matches any sequence of non-path-separators
	'**'          | matches any sequence of characters, including path separators
	'?'           | matches any single non-path-separator character
	'[class]'     | matches any single non-path-separator character against a class of characters ([see below](#character-classes))
	'{alt1,...}'  | matches a sequence of characters if one of the comma-separated alternatives matches
	
	Any character with a special meaning can be escaped with a backslash ('\').
	
	#### Character Classes
	
	Character classes support the following:
	
	Class      | Meaning
	---------- | -------
	'[abc]'    | matches any single character within the set
	'[a-z]'    | matches any single character in the range
	'[^class]' | matches any single character which does *not* match the class
`

// RootCommand represents the jsonnet command
// Deprecated: use NewRootCommand instead.
var RootCommand = &cobra.Command{
	Use:   "jsonnet",
	Short: "Helpers for linting and formatting JSONNet code",
}

// RegisterCommandRecursive adds all jsonnet helpers to the RootCommand
// Deprecated: use NewRootCommand instead.
func RegisterCommandRecursive(parent *cobra.Command) {
	parent.AddCommand(RootCommand)

	RootCommand.AddCommand(FormatCommand)
	RootCommand.AddCommand(LintCommand)
}

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jsonnet",
		Short: "Helpers for linting and formatting JSONNet code",
	}
	cmd.AddCommand(NewFormatCommand(), NewLintCommand())
	return cmd
}
