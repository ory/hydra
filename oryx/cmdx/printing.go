// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmdx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/go-openapi/jsonpointer"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tidwall/gjson"
)

type (
	TableHeader interface {
		Header() []string
	}
	TableRow interface {
		TableHeader
		Columns() []string
		Interface() interface{}
	}
	Table interface {
		TableHeader
		Table() [][]string
		Interface() interface{}
		Len() int
	}
	Nil struct{}

	format string
)

const (
	FormatQuiet       format = "quiet"
	FormatTable       format = "table"
	FormatJSON        format = "json"
	FormatJSONPath    format = "jsonpath"
	FormatJSONPointer format = "jsonpointer"
	FormatJSONPretty  format = "json-pretty"
	FormatYAML        format = "yaml"
	FormatDefault     format = "default"

	FlagFormat = "format"

	None = "<none>"
)

func (Nil) String() string {
	return "null"
}

func (Nil) Interface() interface{} {
	return nil
}

func PrintErrors(cmd *cobra.Command, errs map[string]error) {
	for src, err := range errs {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s: %s\n", src, err.Error())
	}
}

func PrintRow(cmd *cobra.Command, row TableRow) {
	f := getFormat(cmd)

	switch f {
	case FormatQuiet:
		if idAble, ok := row.(interface{ ID() string }); ok {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), idAble.ID())
			break
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), row.Columns()[0])
	case FormatJSON:
		printJSON(cmd.OutOrStdout(), row.Interface(), false, "")
	case FormatYAML:
		printYAML(cmd.OutOrStdout(), row.Interface())
	case FormatJSONPretty:
		printJSON(cmd.OutOrStdout(), row.Interface(), true, "")
	case FormatJSONPath:
		printJSON(cmd.OutOrStdout(), row.Interface(), true, getPath(cmd))
	case FormatJSONPointer:
		printJSON(cmd.OutOrStdout(), filterJSONPointer(cmd, row.Interface()), true, "")
	case FormatTable, FormatDefault:
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 8, 1, '\t', 0)

		fields := row.Columns()
		for i, h := range row.Header() {
			_, _ = fmt.Fprintf(w, "%s\t%s\t\n", h, fields[i])
		}

		_ = w.Flush()
	}
}

func filterJSONPointer(cmd *cobra.Command, data any) any {
	f, err := cmd.Flags().GetString(FlagFormat)
	// unexpected error
	Must(err, "flag access error: %s", err)
	_, jsonptr, found := strings.Cut(f, "=")
	if !found {
		_, _ = fmt.Fprintf(os.Stderr,
			"Format %s is missing a JSON pointer, e.g., --%s=%s=<jsonpointer>. The path syntax is described at https://datatracker.ietf.org/doc/html/draft-ietf-appsawg-json-pointer-07.",
			f, FlagFormat, f)
		os.Exit(1)
	}
	ptr, err := jsonpointer.New(jsonptr)
	Must(err, "invalid JSON pointer: %s", err)

	result, _, err := ptr.Get(data)
	Must(err, "failed to apply JSON pointer: %s", err)

	return result
}

func PrintTable(cmd *cobra.Command, table Table) {
	f := getFormat(cmd)

	switch f {
	case FormatQuiet:
		if table.Len() == 0 {
			fmt.Fprintln(cmd.OutOrStdout())
		}

		if idAble, ok := table.(interface{ IDs() []string }); ok {
			for _, row := range idAble.IDs() {
				fmt.Fprintln(cmd.OutOrStdout(), row)
			}
			break
		}

		for _, row := range table.Table() {
			fmt.Fprintln(cmd.OutOrStdout(), row[0])
		}
	case FormatJSON:
		printJSON(cmd.OutOrStdout(), table.Interface(), false, "")
	case FormatJSONPretty:
		printJSON(cmd.OutOrStdout(), table.Interface(), true, "")
	case FormatJSONPath:
		printJSON(cmd.OutOrStdout(), table.Interface(), true, getPath(cmd))
	case FormatJSONPointer:
		printJSON(cmd.OutOrStdout(), filterJSONPointer(cmd, table.Interface()), true, "")
	case FormatYAML:
		printYAML(cmd.OutOrStdout(), table.Interface())
	default:
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 8, 1, '\t', 0)

		for _, h := range table.Header() {
			fmt.Fprintf(w, "%s\t", h)
		}
		fmt.Fprintln(w)

		for _, row := range table.Table() {
			fmt.Fprintln(w, strings.Join(row, "\t")+"\t")
		}

		_ = w.Flush()
	}
}

type interfacer interface{ Interface() interface{} }

func PrintJSONAble(cmd *cobra.Command, d interface{ String() string }) {
	var path string
	if d == nil {
		d = Nil{}
	}
	switch getFormat(cmd) {
	default:
		_, _ = fmt.Fprint(cmd.OutOrStdout(), d.String())
	case FormatJSON:
		var v interface{} = d
		if i, ok := d.(interfacer); ok {
			v = i
		}
		printJSON(cmd.OutOrStdout(), v, false, "")
	case FormatJSONPath:
		path = getPath(cmd)
		fallthrough
	case FormatJSONPretty:
		var v interface{} = d
		if i, ok := d.(interfacer); ok {
			v = i
		}
		printJSON(cmd.OutOrStdout(), v, true, path)
	case FormatJSONPointer:
		var v interface{} = d
		if i, ok := d.(interfacer); ok {
			v = i
		}
		printJSON(cmd.OutOrStdout(), filterJSONPointer(cmd, v), true, "")
	case FormatYAML:
		var v interface{} = d
		if i, ok := d.(interfacer); ok {
			v = i
		}
		printYAML(cmd.OutOrStdout(), v)
	}
}

func getQuiet(cmd *cobra.Command) bool {
	// ignore the error here as we use this function also when the flag might not be registered
	q, _ := cmd.Flags().GetBool(FlagQuiet)
	return q
}

func getFormat(cmd *cobra.Command) format {
	if getQuiet(cmd) {
		return FormatQuiet
	}

	// ignore the error here as we use this function also when the flag might not be registered
	f, _ := cmd.Flags().GetString(FlagFormat)

	switch {
	case f == string(FormatTable):
		return FormatTable
	case f == string(FormatJSON):
		return FormatJSON
	case strings.HasPrefix(f, string(FormatJSONPath)):
		return FormatJSONPath
	case strings.HasPrefix(f, string(FormatJSONPointer)):
		return FormatJSONPointer
	case f == string(FormatJSONPretty):
		return FormatJSONPretty
	case f == string(FormatYAML):
		return FormatYAML
	default:
		return FormatDefault
	}
}

func getPath(cmd *cobra.Command) string {
	f, err := cmd.Flags().GetString(FlagFormat)
	// unexpected error
	Must(err, "flag access error: %s", err)
	_, path, found := strings.Cut(f, "=")
	if !found {
		_, _ = fmt.Fprintf(os.Stderr,
			"Format %s is missing a path, e.g., --%s=%s=<path>. The path syntax is described at https://github.com/tidwall/gjson/blob/master/SYNTAX.md",
			f, FlagFormat, f)
		os.Exit(1)
	}

	return path
}

func printJSON(w io.Writer, v interface{}, pretty bool, path string) {
	if path != "" {
		temp, err := json.Marshal(v)
		Must(err, "Error encoding JSON: %s", err)
		v = gjson.GetBytes(temp, path).Value()
	}

	e := json.NewEncoder(w)
	if pretty {
		e.SetIndent("", "  ")
	}
	err := e.Encode(v)
	// unexpected error
	Must(err, "Error encoding JSON: %s", err)
}

func printYAML(w io.Writer, v interface{}) {
	j, err := json.Marshal(v)
	Must(err, "Error encoding JSON: %s", err)
	e, err := yaml.JSONToYAML(j)
	Must(err, "Error encoding YAML: %s", err)
	_, _ = w.Write(e)
}

func RegisterJSONFormatFlags(flags *pflag.FlagSet) {
	flags.String(FlagFormat, string(FormatDefault), fmt.Sprintf("Set the output format. One of %s, %s, %s, %s, %s and %s.", FormatDefault, FormatJSON, FormatYAML, FormatJSONPretty, FormatJSONPath, FormatJSONPointer))
}

func RegisterFormatFlags(flags *pflag.FlagSet) {
	RegisterNoiseFlags(flags)
	flags.String(FlagFormat, string(FormatDefault), fmt.Sprintf("Set the output format. One of %s, %s, %s, %s, %s and %s.", FormatTable, FormatJSON, FormatYAML, FormatJSONPretty, FormatJSONPath, FormatJSONPointer))
}

func PrintOpenAPIError(cmd *cobra.Command, err error) error {
	if err == nil {
		return nil
	}

	var be interface {
		Body() []byte
	}
	if !errors.As(err, &be) {
		return err
	}

	body := be.Body()
	didPrettyPrint := false
	if message := gjson.GetBytes(body, "error.message"); message.Exists() {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s\n", message.String())
		didPrettyPrint = true
	}
	if reason := gjson.GetBytes(body, "error.reason"); reason.Exists() {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s\n", reason.String())
		didPrettyPrint = true
	}

	if didPrettyPrint {
		return FailSilently(cmd)
	}

	if body, err := json.MarshalIndent(json.RawMessage(body), "", "  "); err == nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s\nFailed to execute API request, see error above.\n", body)
		return FailSilently(cmd)
	}

	return err
}
