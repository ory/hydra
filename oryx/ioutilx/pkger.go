// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package ioutilx

import (
	"io"
)

// MustReadAll reads a reader or panics.
func MustReadAll(r io.Reader) []byte {
	all, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return all
}
