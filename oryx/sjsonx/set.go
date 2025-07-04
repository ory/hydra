// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sjsonx

import (
	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
)

// SetBytes sets multiple key value pairs in the json object using sjson.SetBytes.
func SetBytes(in []byte, vs map[string]interface{}) (out []byte, err error) {
	out = make([]byte, len(in))
	copy(out, in)
	for k, v := range vs {
		out, err = sjson.SetBytes(out, k, v)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return out, nil
}

// Set sets multiple key value pairs in the json object using sjson.Set.
func Set(in string, vs map[string]interface{}) (out string, err error) {
	out = in
	for k, v := range vs {
		out, err = sjson.Set(out, k, v)
		if err != nil {
			return "", errors.WithStack(err)
		}
	}

	return out, nil
}
