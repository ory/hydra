// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmdx

import (
	"bytes"
	"fmt"
	"slices"
	"strconv"
	"testing"

	"github.com/spf13/cobra"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	dynamicTable struct {
		t  [][]string
		cs int
	}
	dynamicIDAbleTable struct {
		*dynamicTable
		idColumn int
	}
	dynamicRow       []string
	dynamicIDAbleRow struct {
		dynamicRow
		idColumn int
	}
)

var (
	_ Table    = (*dynamicTable)(nil)
	_ Table    = (*dynamicIDAbleTable)(nil)
	_ TableRow = (dynamicRow)(nil)
	_ TableRow = (*dynamicIDAbleRow)(nil)
)

func dynamicHeader(l int) []string {
	h := make([]string, l)
	for i := range h {
		h[i] = "C" + strconv.Itoa(i)
	}
	return h
}

func (d *dynamicTable) Header() []string {
	return dynamicHeader(d.cs)
}

func (d *dynamicTable) Table() [][]string {
	return d.t
}

func (d *dynamicTable) Interface() interface{} {
	return d.t
}

func (d *dynamicIDAbleTable) IDs() []string {
	ids := make([]string, d.Len())
	for i, row := range d.Table() {
		ids[i] = row[d.idColumn]
	}
	return ids
}

func (d *dynamicTable) Len() int {
	return len(d.t)
}

func (d dynamicRow) Header() []string {
	return dynamicHeader(len(d))
}

func (d dynamicRow) Columns() []string {
	return d
}

func (d dynamicRow) Interface() interface{} {
	return d
}

func (d *dynamicIDAbleRow) ID() string {
	return d.dynamicRow[d.idColumn]
}

func TestPrinting(t *testing.T) {
	t.Run("case=format flags", func(t *testing.T) {
		t.Run("format=no value", func(t *testing.T) {
			flags := pflag.NewFlagSet("test flags", pflag.ContinueOnError)
			RegisterFormatFlags(flags)

			require.NoError(t, flags.Parse([]string{}))
			f, err := flags.GetString(FlagFormat)
			require.NoError(t, err)

			assert.Equal(t, FormatDefault, format(f))
		})
	})

	t.Run("method=table row", func(t *testing.T) {
		t.Run("case=all formats", func(t *testing.T) {
			tr := dynamicRow{"AAA", "BBB", "CCC"}
			allFields := append(tr.Header(), tr...)

			for _, tc := range []struct {
				fArgs     []string
				contained []string
			}{
				{
					fArgs:     []string{"--" + FlagFormat, string(FormatTable)},
					contained: allFields,
				},
				{
					fArgs:     []string{"--" + FlagQuiet},
					contained: []string{tr[0]},
				},
				{
					fArgs:     []string{"--" + FlagFormat, string(FormatJSON)},
					contained: tr,
				},
				{
					fArgs:     []string{"--" + FlagFormat, string(FormatJSONPretty)},
					contained: tr,
				},
				{
					fArgs:     []string{"--" + FlagFormat, string(FormatJSONPointer) + "=/0"},
					contained: []string{"AAA"},
				},
				{
					fArgs:     []string{"--" + FlagFormat, string(FormatJSONPointer) + "=/2"},
					contained: []string{"CCC"},
				},
				{
					fArgs:     []string{"--" + FlagFormat, string(FormatJSONPointer) + "=/1"},
					contained: []string{"BBB"},
				},
				{
					fArgs:     []string{"--" + FlagFormat, string(FormatJSONPath) + "=0"},
					contained: []string{"AAA"},
				},
				{
					fArgs:     []string{"--" + FlagFormat, string(FormatJSONPath) + "=2"},
					contained: []string{"CCC"},
				},
				{
					fArgs:     []string{"--" + FlagFormat, string(FormatJSONPath) + "=[0,1]"},
					contained: []string{"AAA", "BBB"},
				},
				{
					fArgs:     []string{"--" + FlagFormat, string(FormatYAML)},
					contained: tr,
				},
			} {
				t.Run(fmt.Sprintf("format=%v", tc.fArgs), func(t *testing.T) {
					cmd := &cobra.Command{Use: "x"}
					RegisterFormatFlags(cmd.Flags())

					out := &bytes.Buffer{}
					cmd.SetOut(out)
					require.NoError(t, cmd.Flags().Parse(tc.fArgs))

					PrintRow(cmd, tr)

					for _, s := range tc.contained {
						assert.Contains(t, out.String(), s, "%s", out.String())
					}
					notContained := slices.DeleteFunc(slices.Clone(allFields), func(s string) bool {
						return slices.Contains(tc.contained, s)
					})
					for _, s := range notContained {
						assert.NotContains(t, out.String(), s, "%s", out.String())
					}

					assert.Equal(t, "\n", out.String()[len(out.String())-1:])
				})
			}
		})

		t.Run("case=uses ID()", func(t *testing.T) {
			tr := &dynamicIDAbleRow{
				dynamicRow: []string{"foo", "bar"},
				idColumn:   1,
			}

			cmd := &cobra.Command{Use: "x"}
			RegisterFormatFlags(cmd.Flags())

			out := &bytes.Buffer{}
			cmd.SetOut(out)
			require.NoError(t, cmd.Flags().Parse([]string{"--" + FlagQuiet}))

			PrintRow(cmd, tr)

			assert.Equal(t, tr.dynamicRow[1]+"\n", out.String())
		})
	})

	t.Run("method=table", func(t *testing.T) {
		t.Run("case=full table", func(t *testing.T) {
			tb := &dynamicTable{
				t: [][]string{
					{"a0", "b0", "c0"},
					{"a1", "b1", "c1"},
				},
				cs: 3,
			}
			allFields := append(tb.Header(), append(tb.t[0], tb.t[1]...)...)

			for _, tc := range []struct {
				fArgs     []string
				contained []string
			}{
				{
					fArgs:     []string{"--" + FlagFormat, string(FormatTable)},
					contained: allFields,
				},
				{
					fArgs:     []string{"--" + FlagQuiet},
					contained: []string{tb.t[0][0], tb.t[1][0]},
				},
				{
					fArgs:     []string{"--" + FlagFormat, string(FormatJSON)},
					contained: append(tb.t[0], tb.t[1]...),
				},
				{
					fArgs:     []string{"--" + FlagFormat, string(FormatJSONPretty)},
					contained: append(tb.t[0], tb.t[1]...),
				},
				{
					fArgs:     []string{"--" + FlagFormat, string(FormatJSONPath) + "=1.1"},
					contained: []string{tb.t[1][1]},
				},
				{
					fArgs:     []string{"--" + FlagFormat, string(FormatJSONPointer) + "=/1/1"},
					contained: []string{tb.t[1][1]},
				},
				{
					fArgs:     []string{"--" + FlagFormat, string(FormatYAML)},
					contained: append(tb.t[0], tb.t[1]...),
				},
			} {
				t.Run(fmt.Sprintf("format=%v", tc.fArgs), func(t *testing.T) {
					cmd := &cobra.Command{Use: "x"}
					RegisterFormatFlags(cmd.Flags())

					out := &bytes.Buffer{}
					cmd.SetOut(out)
					require.NoError(t, cmd.Flags().Parse(tc.fArgs))

					PrintTable(cmd, tb)

					for _, s := range tc.contained {
						assert.Contains(t, out.String(), s, "%s", out.String())
					}
					notContained := slices.DeleteFunc(slices.Clone(allFields), func(s string) bool {
						return slices.Contains(tc.contained, s)
					})
					for _, s := range notContained {
						assert.NotContains(t, out.String(), s, "%s", out.String())
					}

					assert.Equal(t, "\n", out.String()[len(out.String())-1:])
				})
			}
		})

		t.Run("case=empty table", func(t *testing.T) {
			tb := &dynamicTable{
				t:  nil,
				cs: 1,
			}

			for _, tc := range []struct {
				fArgs    []string
				expected string
			}{
				{
					fArgs:    []string{"--" + FlagFormat, string(FormatTable)},
					expected: "C0\t",
				},
				{
					fArgs:    []string{"--" + FlagQuiet},
					expected: "",
				},
				{
					fArgs:    []string{"--" + FlagFormat, string(FormatJSON)},
					expected: "null",
				},
				{
					fArgs:    []string{"--" + FlagFormat, string(FormatJSONPretty)},
					expected: "null",
				},
				{
					fArgs:    []string{"--" + FlagFormat, string(FormatJSONPath) + "=foo"},
					expected: "null",
				},
				{
					fArgs:    []string{"--" + FlagFormat, string(FormatYAML)},
					expected: "null",
				},
			} {
				t.Run(fmt.Sprintf("format=%v", tc.fArgs), func(t *testing.T) {
					cmd := &cobra.Command{Use: "x"}
					RegisterFormatFlags(cmd.Flags())

					out := &bytes.Buffer{}
					cmd.SetOut(out)
					require.NoError(t, cmd.Flags().Parse(tc.fArgs))

					PrintTable(cmd, tb)

					assert.Equal(t, tc.expected+"\n", out.String())
				})
			}
		})

		t.Run("case=uses IDs()", func(t *testing.T) {
			tb := &dynamicIDAbleTable{
				dynamicTable: &dynamicTable{
					t: [][]string{
						{"a0", "b0", "c0"},
						{"a1", "b1", "c1"},
					},
					cs: 3,
				},
				idColumn: 1,
			}
			cmd := &cobra.Command{Use: "x"}
			RegisterFormatFlags(cmd.Flags())

			out := &bytes.Buffer{}
			cmd.SetOut(out)
			require.NoError(t, cmd.Flags().Parse([]string{"--" + FlagQuiet}))

			PrintTable(cmd, tb)

			assert.Equal(t, tb.t[0][1]+"\n"+tb.t[1][1]+"\n", out.String())
		})
	})

	t.Run("method=jsonable", func(t *testing.T) {
		t.Run("case=nil", func(t *testing.T) {
			for _, f := range []format{FormatDefault, FormatJSON, FormatJSONPretty, FormatJSONPath, FormatJSONPointer, FormatYAML} {
				t.Run("format="+string(f), func(t *testing.T) {
					out := &bytes.Buffer{}
					cmd := &cobra.Command{}
					cmd.SetOut(out)
					RegisterJSONFormatFlags(cmd.Flags())

					PrintJSONAble(cmd, nil)

					assert.Equal(t, "null", out.String())
				})
			}
		})
	})

}
