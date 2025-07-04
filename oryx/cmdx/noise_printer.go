// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmdx

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type ConditionalPrinter struct {
	w     io.Writer
	print bool
}

const (
	FlagQuiet = "quiet"
)

func RegisterNoiseFlags(flags *pflag.FlagSet) {
	flags.BoolP(FlagQuiet, FlagQuiet[:1], false, "Be quiet with output printing.")
}

// NewLoudOutPrinter returns a ConditionalPrinter that
// only prints to cmd.OutOrStdout when --quiet is not set
func NewLoudOutPrinter(cmd *cobra.Command) *ConditionalPrinter {
	quiet, err := cmd.Flags().GetBool(FlagQuiet)
	if err != nil {
		Fatalf(err.Error())
	}

	return &ConditionalPrinter{
		w:     cmd.OutOrStdout(),
		print: !quiet,
	}
}

// NewQuietOutPrinter returns a ConditionalPrinter that
// only prints to cmd.OutOrStdout when --quiet is set
func NewQuietOutPrinter(cmd *cobra.Command) *ConditionalPrinter {
	quiet, err := cmd.Flags().GetBool(FlagQuiet)
	if err != nil {
		Fatalf(err.Error())
	}

	return &ConditionalPrinter{
		w:     cmd.OutOrStdout(),
		print: quiet,
	}
}

// NewLoudErrPrinter returns a ConditionalPrinter that
// only prints to cmd.ErrOrStderr when --quiet is not set
func NewLoudErrPrinter(cmd *cobra.Command) *ConditionalPrinter {
	quiet, err := cmd.Flags().GetBool(FlagQuiet)
	if err != nil {
		Fatalf(err.Error())
	}

	return &ConditionalPrinter{
		w:     cmd.ErrOrStderr(),
		print: !quiet,
	}
}

// NewQuietErrPrinter returns a ConditionalPrinter that
// only prints to cmd.ErrOrStderr when --quiet is set
func NewQuietErrPrinter(cmd *cobra.Command) *ConditionalPrinter {
	quiet, err := cmd.Flags().GetBool(FlagQuiet)
	if err != nil {
		Fatalf(err.Error())
	}

	return &ConditionalPrinter{
		w:     cmd.ErrOrStderr(),
		print: quiet,
	}
}

// NewLoudPrinter returns a ConditionalPrinter that
// only prints to w when --quiet is not set
func NewLoudPrinter(cmd *cobra.Command, w io.Writer) *ConditionalPrinter {
	quiet, err := cmd.Flags().GetBool(FlagQuiet)
	if err != nil {
		Fatalf(err.Error())
	}

	return &ConditionalPrinter{
		w:     w,
		print: !quiet,
	}
}

// NewQuietPrinter returns a ConditionalPrinter that
// only prints to w when --quiet is set
func NewQuietPrinter(cmd *cobra.Command, w io.Writer) *ConditionalPrinter {
	quiet, err := cmd.Flags().GetBool(FlagQuiet)
	if err != nil {
		Fatalf(err.Error())
	}

	return &ConditionalPrinter{
		w:     w,
		print: quiet,
	}
}

func NewConditionalPrinter(w io.Writer, print bool) *ConditionalPrinter {
	return &ConditionalPrinter{
		w:     w,
		print: print,
	}
}

func (p *ConditionalPrinter) Println(a ...interface{}) (n int, err error) {
	if p.print {
		return fmt.Fprintln(p.w, a...)
	}
	return
}

func (p *ConditionalPrinter) Print(a ...interface{}) (n int, err error) {
	if p.print {
		return fmt.Fprint(p.w, a...)
	}
	return
}

func (p *ConditionalPrinter) Printf(format string, a ...interface{}) (n int, err error) {
	if p.print {
		return fmt.Fprintf(p.w, format, a...)
	}
	return
}
