// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonnetsecure

import (
	"bufio"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	GiB uint64 = 1024 * 1024 * 1024
	// Generous limit on virtual memory including the peak memory allocated by the Go runtime, the Jsonnet VM,
	// and the Jsonnet script.
	// This number was acquired by running:
	// Found by trial and error with:
	// `ulimit -Sv 1048576 && echo '{"Snippet": "{user_id: std.repeat(\'a\', 1000)}"}' | kratos jsonnet -0`
	// NOTE: Ideally we'd like to limit RSS but that is not possible on Linux with `ulimit/setrlimit(2)` - only with cgroups.
	virtualMemoryLimitBytes = 2 * GiB
)

func NewJsonnetCmd() *cobra.Command {
	var null bool
	cmd := &cobra.Command{
		Use:    "jsonnet",
		Short:  "Run Jsonnet as a CLI command",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			// This could fail because current limits are lower than what we tried to set,
			// so we still continue in this case.
			SetVirtualMemoryLimit(virtualMemoryLimitBytes)

			if null {
				return scan(cmd.OutOrStdout(), cmd.InOrStdin())
			}

			input, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return errors.Wrap(err, "failed to read from stdin")
			}

			json, err := eval(input)
			if err != nil {
				return errors.Wrap(err, "failed to evaluate jsonnet")
			}

			if _, err := io.WriteString(cmd.OutOrStdout(), json); err != nil {
				return errors.Wrap(err, "failed to write json output")
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&null, "null", "0", false,
		`Read multiple snippets and parameters from stdin separated by null bytes.
Output will be in the same order as inputs, separated by null bytes.
Evaluation errors will also be reported to stdout, separated by null bytes.
Non-recoverable errors are written to stderr and the program will terminate with a non-zero exit code.`)

	return cmd
}

func scan(w io.Writer, r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(splitNull)
	for scanner.Scan() {
		json, err := eval(scanner.Bytes())
		if err != nil {
			json = fmt.Sprintf("ERROR: %s", err)
		}
		if _, err := fmt.Fprintf(w, "%s%c", json, 0); err != nil {
			return errors.Wrap(err, "failed to write json output")
		}
	}
	return errors.Wrap(scanner.Err(), "failed to read from stdin")
}

func eval(input []byte) (json string, err error) {
	var params processParameters
	if err := params.Decode(input); err != nil {
		return "", err
	}

	vm := MakeSecureVM()

	for _, it := range params.ExtCodes {
		vm.ExtCode(it.Key, it.Value)
	}
	for _, it := range params.ExtVars {
		vm.ExtVar(it.Key, it.Value)
	}
	for _, it := range params.TLACodes {
		vm.TLACode(it.Key, it.Value)
	}
	for _, it := range params.TLAVars {
		vm.TLAVar(it.Key, it.Value)
	}

	return vm.EvaluateAnonymousSnippet(params.Filename, params.Snippet)
}
