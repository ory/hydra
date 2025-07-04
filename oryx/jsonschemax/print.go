// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonschemax

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/ory/jsonschema/v3"
)

func FormatValidationErrorForCLI(w io.Writer, conf []byte, err error) {
	if err == nil {
		return
	}

	if e := new(jsonschema.ValidationError); errors.As(err, &e) {
		_, _ = fmt.Fprintln(w, "The configuration contains values or keys which are invalid:")
		pointer, validation := FormatError(e)

		if pointer == "#" {
			if len(e.Causes) == 0 {
				_, _ = fmt.Fprintln(w, "(root)")
				_, _ = fmt.Fprintln(w, "^-- "+validation)
				_, _ = fmt.Fprintln(w, "")
			}
		} else {
			spaces := make([]string, len(pointer)+3)
			_, _ = fmt.Fprintf(w, "%s: %+v", pointer, gjson.GetBytes(conf, pointer).Value())
			_, _ = fmt.Fprintln(w, "")
			_, _ = fmt.Fprintf(w, "%s^-- %s", strings.Join(spaces, " "), validation)
			_, _ = fmt.Fprintln(w, "")
			_, _ = fmt.Fprintln(w, "")
		}

		for _, cause := range e.Causes {
			FormatValidationErrorForCLI(w, conf, cause)
		}
		return
	}
}

func FormatError(e *jsonschema.ValidationError) (string, string) {
	var (
		err     error
		pointer string
		message string
	)

	pointer = e.InstancePtr
	message = e.Message
	switch ctx := e.Context.(type) {
	case *jsonschema.ValidationErrorContextRequired:
		if len(ctx.Missing) > 0 {
			message = "one or more required properties are missing"
			pointer = ctx.Missing[0]
		}
	}

	// We can ignore the error as it will simply echo the pointer.
	pointer, err = JSONPointerToDotNotation(pointer)
	if err != nil {
		pointer = e.InstancePtr
	}

	return pointer, message
}
