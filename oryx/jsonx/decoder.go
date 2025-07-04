// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonx

import (
	"encoding/json"
	"io"
)

// NewStrictDecoder is a shorthand for json.Decoder.DisallowUnknownFields
func NewStrictDecoder(b io.Reader) *json.Decoder {
	d := json.NewDecoder(b)
	d.DisallowUnknownFields()
	return d
}
