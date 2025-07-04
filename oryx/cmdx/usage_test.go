// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmdx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestUsageTemplating(t *testing.T) {
	root := &cobra.Command{
		Use:   "root",
		Short: "{{ .Name }}",
	}
	cmdWithTemplate := &cobra.Command{
		Use:     "with-template",
		Long:    "{{ .Name }}",
		Example: "{{ .Name }}",
	}
	cmdWithoutTemplate := &cobra.Command{
		Use:     "without-template",
		Long:    "{{ .Name }}",
		Example: "{{ .Name }}",
	}
	root.AddCommand(cmdWithTemplate, cmdWithoutTemplate)

	EnableUsageTemplating(root)
	DisableUsageTemplating(cmdWithoutTemplate)
	assert.NotContains(t, root.UsageString(), "{{ .Name }}")
	assert.NotContains(t, cmdWithTemplate.UsageString(), "{{ .Name }}")
	assert.Contains(t, cmdWithoutTemplate.UsageString(), "{{ .Name }}")
}

func TestAssertUsageTemplates(t *testing.T) {
	var cmdsCalled []string
	AddUsageTemplateFunc("called", func(use string) string {
		cmdsCalled = append(cmdsCalled, use)
		return use
	})

	root := &cobra.Command{
		Use:   "root",
		Short: "{{ called .Use }}",
	}
	child := &cobra.Command{
		Use:  "child",
		Long: "{{ called .Use }}",
	}
	otherChild := &cobra.Command{
		Use:     "other-child",
		Example: "{{ called .Use }}",
	}
	childChild := &cobra.Command{
		Use:     "child-child",
		Example: "{{ called .Use }}",
	}
	root.AddCommand(child, otherChild)
	child.AddCommand(childChild)

	EnableUsageTemplating(root)

	require.NotPanics(t, func() {
		AssertUsageTemplates(&panicT{}, root)
	})
	assert.ElementsMatch(t, []string{root.Use, child.Use, otherChild.Use, childChild.Use}, cmdsCalled)
}

type panicT struct{}

func (t *panicT) FailNow() {
	panic("failing")
}

func (*panicT) Errorf(format string, args ...interface{}) {
	panic("erroring: " + fmt.Sprintf(format, args...))
}

var _ require.TestingT = (*panicT)(nil)
