// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"encoding/json"
	"fmt"

	hydra "github.com/ory/hydra-client-go/v2"
)

type (
	outputJsonWebKey struct {
		Set string `json:"set"`
		hydra.JsonWebKey
	}
	outputJSONWebKeyCollection struct {
		// Set is empty when the collection holds keys from more than one set.
		Set  string             `json:"set,omitempty"`
		Keys []outputJsonWebKey `json:"keys"`
	}
)

func newOutputJsonWebKeys(set string, keys []hydra.JsonWebKey) []outputJsonWebKey {
	out := make([]outputJsonWebKey, len(keys))
	for i, key := range keys {
		out[i] = outputJsonWebKey{Set: set, JsonWebKey: key}
	}
	return out
}

// MarshalJSON shadows the marshaler promoted from the embedded
// hydra.JsonWebKey, which would otherwise drop the Set field.
func (i outputJsonWebKey) MarshalJSON() ([]byte, error) {
	key, err := i.JsonWebKey.ToMap()
	if err != nil {
		return nil, err
	}
	key["set"] = i.Set
	return json.Marshal(key)
}

func (outputJsonWebKey) Header() []string {
	return []string{"SET ID", "KEY ID", "ALGORITHM", "USE"}
}

func (i outputJsonWebKey) Columns() []string {
	return []string{i.Set, i.Kid, i.Alg, i.Use}
}

func (i outputJsonWebKey) Interface() interface{} {
	return i
}

func (outputJSONWebKeyCollection) Header() []string {
	return outputJsonWebKey{}.Header()
}

func (c outputJSONWebKeyCollection) Table() [][]string {
	rows := make([][]string, len(c.Keys))
	for i, key := range c.Keys {
		rows[i] = key.Columns()
	}
	return rows
}

func (c outputJSONWebKeyCollection) Interface() interface{} {
	return c
}

func (c outputJSONWebKeyCollection) Len() int {
	return len(c.Keys)
}

func (c outputJSONWebKeyCollection) IDs() []string {
	ids := make([]string, len(c.Keys))
	for i, key := range c.Keys {
		ids[i] = fmt.Sprintf("%v", key.Kid)
	}
	return ids
}
