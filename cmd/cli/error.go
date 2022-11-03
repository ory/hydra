// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"bytes"
	"encoding/json"
)

func FormatSwaggerError(err error) string {
	if err == nil {
		return ""
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(err); err != nil {
		panic(err)
	}

	var e struct {
		Payload json.RawMessage
	}
	if err := json.NewDecoder(&b).Decode(&e); err != nil {
		panic(err)
	}

	if len(e.Payload) == 0 {
		return err.Error()
	}

	dec := json.NewEncoder(&b)
	dec.SetIndent("", "  ")
	if err := dec.Encode(e.Payload); err != nil {
		panic(err)
	}

	return b.String()
}
