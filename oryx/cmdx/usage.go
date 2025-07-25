// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmdx

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var usageTemplateFuncs = sprig.TxtFuncMap()

// AddUsageTemplateFunc adds a template function to the usage template.
func AddUsageTemplateFunc(name string, f interface{}) {
	usageTemplateFuncs[name] = f
}

const (
	helpTemplate = `{{insertTemplate . (or .Long .Short) | trimTrailingWhitespaces}}

{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
	usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{insertTemplate . .Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
)

// EnableUsageTemplating enables gotemplates for usage strings, i.e. cmd.Short, cmd.Long, and cmd.Example.
// The data for the template is the command itself. Especially useful are `.Root.Name` and `.CommandPath`.
// This will be inherited by all subcommands, so enabling it on the root command is sufficient.
func EnableUsageTemplating(cmds ...*cobra.Command) {
	cobra.AddTemplateFunc("insertTemplate", TemplateCommandField)
	for _, cmd := range cmds {
		cmd.SetHelpTemplate(helpTemplate)
		cmd.SetUsageTemplate(usageTemplate)
	}
}

func TemplateCommandField(cmd *cobra.Command, field string) (string, error) {
	t := template.New("")
	t.Funcs(usageTemplateFuncs)
	t, err := t.Parse(field)
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	if err := t.Execute(&out, cmd); err != nil {
		return "", err
	}
	return out.String(), nil
}

// DisableUsageTemplating resets the commands usage template to the default.
// This can be used to undo the effects of EnableUsageTemplating, specifically for a subcommand.
func DisableUsageTemplating(cmds ...*cobra.Command) {
	defaultCmd := new(cobra.Command)
	for _, cmd := range cmds {
		cmd.SetHelpTemplate(defaultCmd.HelpTemplate())
		cmd.SetUsageTemplate(defaultCmd.UsageTemplate())
	}
}

// AssertUsageTemplates asserts that the usage string of the commands are properly templated.
func AssertUsageTemplates(t require.TestingT, cmd *cobra.Command) {
	var usage, help string
	require.NotPanics(t, func() {
		usage = cmd.UsageString()

		out, err := cmd.OutOrStdout(), cmd.ErrOrStderr()
		bb := new(bytes.Buffer)

		cmd.SetOut(bb)
		cmd.SetErr(bb)
		require.NoError(t, cmd.Help())
		help = bb.String()

		cmd.SetOut(out)
		cmd.SetErr(err)
	})
	assert.NotContains(t, usage, "{{")
	assert.NotContains(t, help, "{{")
	for _, child := range cmd.Commands() {
		AssertUsageTemplates(t, child)
	}
}
