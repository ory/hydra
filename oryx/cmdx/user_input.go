// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmdx

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// asks for confirmation with the question string s and reads the answer
// pass nil to use os.Stdin and os.Stdout
func AskForConfirmation(s string, stdin io.Reader, stdout io.Writer) bool {
	if stdin == nil {
		stdin = os.Stdin
	}
	if stdout == nil {
		stdout = os.Stdout
	}

	ok, err := AskScannerForConfirmation(s, bufio.NewReader(stdin), stdout)
	if err != nil {
		Must(err, "Unable to confirm: %s", err)
	}

	return ok
}

func AskScannerForConfirmation(s string, reader *bufio.Reader, stdout io.Writer) (bool, error) {
	if stdout == nil {
		stdout = os.Stdout
	}

	for {
		_, err := fmt.Fprintf(stdout, "%s [y/n]: ", s)
		if err != nil {
			return false, errors.Wrap(err, "unable to print to stdout")
		}

		response, err := reader.ReadString('\n')
		if err != nil {
			return false, errors.Wrap(err, "unable to read from stdin")
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true, nil
		} else if response == "n" || response == "no" {
			return false, nil
		}
	}
}
