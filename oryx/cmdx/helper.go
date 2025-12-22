// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmdx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"golang.org/x/sync/errgroup"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"

	"github.com/pkg/errors"

	"github.com/ory/x/logrusx"
)

var (
	// ErrNilDependency is returned if a dependency is missing.
	ErrNilDependency = fmt.Errorf("a dependency was expected to be defined but is nil. Please open an issue with the stack trace")
	// ErrNoPrintButFail is returned to detect a failure state that was already reported to the user in some way
	ErrNoPrintButFail = fmt.Errorf("this error should never be printed")

	debugStdout, debugStderr = io.Discard, io.Discard
)

func init() {
	if os.Getenv("DEBUG") != "" {
		debugStdout = os.Stdout
		debugStderr = os.Stderr
	}
}

// FailSilently is supposed to be used within a commands RunE function.
// It silences cobras error handling and returns the ErrNoPrintButFail error.
func FailSilently(cmd *cobra.Command) error {
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	return errors.WithStack(ErrNoPrintButFail)
}

// Must fatals with the optional message if err is not nil.
// Deprecated: do not use this function in commands, as it makes it impossible to test them. Instead, return the error.
func Must(err error, message string, args ...interface{}) {
	if err == nil {
		return
	}

	_, _ = fmt.Fprintf(os.Stderr, message+"\n", args...)
	os.Exit(1)
}

// CheckResponse fatals if err is nil or the response.StatusCode does not match the expectedStatusCode
// Deprecated: do not use this function in commands, as it makes it impossible to test them. Instead, return the error.
func CheckResponse(err error, expectedStatusCode int, response *http.Response) {
	Must(err, "Command failed because error occurred: %s", err)

	if response.StatusCode != expectedStatusCode {
		out, err := io.ReadAll(response.Body)
		if err != nil {
			out = []byte{}
		}
		pretty, err := json.MarshalIndent(json.RawMessage(out), "", "\t")
		if err == nil {
			out = pretty
		}

		Fatalf(
			`Command failed because status code %d was expected but code %d was received.

Response payload:

%s`,
			expectedStatusCode,
			response.StatusCode,
			out,
		)
	}
}

// FormatResponse takes an object and prints a json.MarshalIdent version of it or fatals.
// Deprecated: do not use this function in commands, as it makes it impossible to test them. Instead, return the error.
func FormatResponse(o interface{}) string {
	out, err := json.MarshalIndent(o, "", "\t")
	Must(err, `Command failed because an error occurred while prettifying output: %s`, err)
	return string(out)
}

// Fatalf prints to os.Stderr and exists with code 1.
// Deprecated: do not use this function in commands, as it makes it impossible to test them. Instead, return the error.
func Fatalf(message string, args ...interface{}) {
	if len(args) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, message+"\n", args...)
	} else {
		_, _ = fmt.Fprintln(os.Stderr, message)
	}
	os.Exit(1)
}

// ExpectDependency expects every dependency to be not nil or it fatals.
// Deprecated: do not use this function in commands, as it makes it impossible to test them. Instead, return the error.
func ExpectDependency(logger *logrusx.Logger, dependencies ...interface{}) {
	if logger == nil {
		panic("missing logger for dependency check")
	}
	for _, d := range dependencies {
		if d == nil {
			logger.WithError(errors.WithStack(ErrNilDependency)).Fatalf("A fatal issue occurred.")
		}
	}
}

// CallbackWriter will execute each callback once the message is received.
// The full matched message is passed to the callback. An error returned from the callback is returned by Write.
type CallbackWriter struct {
	Callbacks map[string]func([]byte) error
	buf       bytes.Buffer
}

func (c *CallbackWriter) Write(msg []byte) (int, error) {
	for p, cb := range c.Callbacks {
		if bytes.Contains(msg, []byte(p)) {
			if err := cb(msg); err != nil {
				return 0, err
			}
		}
	}
	return c.buf.Write(msg)
}

func (c *CallbackWriter) String() string {
	return c.buf.String()
}

var _ io.Writer = (*CallbackWriter)(nil)

func prepareCmd(cmd *cobra.Command, stdIn io.Reader, stdOut, stdErr io.Writer, args []string) {
	cmd.SetIn(stdIn)
	outs := []io.Writer{debugStdout}
	if stdOut != nil {
		outs = append(outs, stdOut)
	}
	cmd.SetOut(io.MultiWriter(outs...))
	errs := []io.Writer{debugStderr}
	if stdErr != nil {
		errs = append(errs, stdErr)
	}
	cmd.SetErr(io.MultiWriter(errs...))

	if args == nil {
		args = []string{}
	}
	cmd.SetArgs(args)
}

// ExecBackgroundCtx runs the cobra command in the background.
func ExecBackgroundCtx(ctx context.Context, cmd *cobra.Command, stdIn io.Reader, stdOut, stdErr io.Writer, args ...string) *errgroup.Group {
	prepareCmd(cmd, stdIn, stdOut, stdErr, args)

	eg := &errgroup.Group{}
	eg.Go(func() error {
		defer cmd.SetIn(nil)
		return cmd.ExecuteContext(ctx)
	})

	return eg
}

// Exec runs the provided cobra command with the given reader as STD_IN and the given args.
// Returns STD_OUT, STD_ERR and the error from the execution.
func Exec(t testing.TB, cmd *cobra.Command, stdIn io.Reader, args ...string) (string, string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	return ExecCtx(ctx, cmd, stdIn, args...)
}

func ExecCtx(ctx context.Context, cmd *cobra.Command, stdIn io.Reader, args ...string) (string, string, error) {
	stdOut, stdErr := &bytes.Buffer{}, &bytes.Buffer{}

	prepareCmd(cmd, stdIn, stdOut, stdErr, args)

	// needs to be on a separate line to ensure that the output buffers are read AFTER the command ran
	err := cmd.ExecuteContext(ctx)

	return stdOut.String(), stdErr.String(), err
}

// ExecNoErr is a helper that assumes a successful run from Exec.
// Returns STD_OUT.
func ExecNoErr(t testing.TB, cmd *cobra.Command, args ...string) string {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	return ExecNoErrCtx(ctx, t, cmd, args...)
}

func ExecNoErrCtx(ctx context.Context, t require.TestingT, cmd *cobra.Command, args ...string) string {
	stdOut, stdErr, err := ExecCtx(ctx, cmd, nil, args...)
	if err == nil {
		require.Len(t, stdErr, 0, "std_out: %s\nstd_err: %s", stdOut, stdErr)
	} else {
		require.ErrorIsf(t, err, context.Canceled, "std_out: %s\nstd_err: %s", stdOut, stdErr)
	}
	return stdOut
}

// ExecExpectedErr is a helper that assumes a failing run from Exec returning ErrNoPrintButFail
// Returns STD_ERR.
func ExecExpectedErr(t testing.TB, cmd *cobra.Command, args ...string) string {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	return ExecExpectedErrCtx(ctx, t, cmd, args...)
}

func ExecExpectedErrCtx(ctx context.Context, t require.TestingT, cmd *cobra.Command, args ...string) string {
	stdOut, stdErr, err := ExecCtx(ctx, cmd, nil, args...)
	require.True(t, errors.Is(err, ErrNoPrintButFail), "std_out: %s\nstd_err: %s", stdOut, stdErr)
	require.Len(t, stdOut, 0, stdErr)
	return stdErr
}

type CommandExecuter struct {
	New            func() *cobra.Command
	Ctx            context.Context
	PersistentArgs []string
}

func (c *CommandExecuter) Exec(stdin io.Reader, args ...string) (string, string, error) {
	return ExecCtx(c.Ctx, c.New(), stdin, append(c.PersistentArgs, args...)...)
}

func (c *CommandExecuter) ExecBackground(stdin io.Reader, stdOut, stdErr io.Writer, args ...string) *errgroup.Group {
	return ExecBackgroundCtx(c.Ctx, c.New(), stdin, stdOut, stdErr, append(c.PersistentArgs, args...)...)
}

func (c *CommandExecuter) ExecNoErr(t require.TestingT, args ...string) string {
	return ExecNoErrCtx(c.Ctx, t, c.New(), append(c.PersistentArgs, args...)...)
}

func (c *CommandExecuter) ExecExpectedErr(t require.TestingT, args ...string) string {
	return ExecExpectedErrCtx(c.Ctx, t, c.New(), append(c.PersistentArgs, args...)...)
}

type URL struct {
	url.URL
}

var _ pflag.Value = (*URL)(nil)

func (u *URL) Set(s string) error {
	uu, err := url.Parse(s)
	if err != nil {
		return err
	}
	u.URL = *uu
	return nil
}

func (*URL) Type() string {
	return "url"
}
