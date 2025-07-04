// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmdx

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConditionalPrinter(t *testing.T) {
	const (
		msgAlwaysOut = "always out"
		msgAlwaysErr = "always err"
		msgQuietOut  = "quiet out"
		msgQuietErr  = "quiet err"
		msgLoudOut   = "loud out"
		msgLoudErr   = "loud err"
		msgArgsSet   = "args were set"
	)
	setup := func() *cobra.Command {
		cmd := &cobra.Command{
			Use: "test cmd",
			Run: func(cmd *cobra.Command, args []string) {
				_, _ = fmt.Fprint(cmd.OutOrStdout(), msgAlwaysOut)
				_, _ = fmt.Fprint(cmd.ErrOrStderr(), msgAlwaysErr)
				_, _ = NewQuietOutPrinter(cmd).Print(msgQuietOut)
				_, _ = NewQuietErrPrinter(cmd).Print(msgQuietErr)
				_, _ = NewLoudOutPrinter(cmd).Print(msgLoudOut)
				_, _ = NewLoudErrPrinter(cmd).Print(msgLoudErr)
				_, _ = NewConditionalPrinter(cmd.OutOrStdout(), len(args) > 0).Print(msgArgsSet)
			},
		}
		RegisterNoiseFlags(cmd.Flags())
		return cmd
	}

	for _, tc := range []struct {
		stdErrMsg, stdOutMsg, args []string
		setQuiet                   bool
	}{
		{
			stdOutMsg: []string{msgLoudOut},
			stdErrMsg: []string{msgLoudErr},
			setQuiet:  false,
			args:      []string{},
		},
		{
			stdOutMsg: []string{msgQuietOut},
			stdErrMsg: []string{msgQuietErr},
			setQuiet:  true,
			args:      []string{},
		},
		{
			stdOutMsg: []string{msgQuietOut, msgArgsSet},
			stdErrMsg: []string{msgQuietErr},
			setQuiet:  true,
			args:      []string{"foo"},
		},
	} {
		t.Run(fmt.Sprintf("case=quiet:%v", tc.setQuiet), func(t *testing.T) {
			cmd := setup()
			if tc.setQuiet {
				require.NoError(t, cmd.Flags().Set(FlagQuiet, "true"))
			}
			out, err := &bytes.Buffer{}, &bytes.Buffer{}
			cmd.SetOut(out)
			cmd.SetErr(err)
			cmd.SetArgs(tc.args)

			require.NoError(t, cmd.Execute())
			assert.Equal(t, strings.Join(append([]string{msgAlwaysOut}, tc.stdOutMsg...), ""), out.String())
			assert.Equal(t, strings.Join(append([]string{msgAlwaysErr}, tc.stdErrMsg...), ""), err.String())
		})
	}
}
