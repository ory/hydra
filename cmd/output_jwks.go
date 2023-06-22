// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	hydra "github.com/ory/hydra-client-go/v2"
)

type (
	outputJsonWebKey struct {
		Set string `json:"set"`
		hydra.JsonWebKey
	}
	outputJSONWebKeyCollection struct {
		Set  string             `json:"set"`
		Keys []hydra.JsonWebKey `json:"keys"`
	}
)

func (outputJsonWebKey) Header() []string {
	return []string{"SET ID", "KEY ID", "ALGORITHM", "USE"}
}

func (i outputJsonWebKey) Columns() []string {
	data := [7]string{
		i.Set,
		i.Kid,
		i.Alg,
		i.Use,
	}
	return data[:]
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
		rows[i] = outputJsonWebKey{Set: c.Set, JsonWebKey: key}.Columns()
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
	for i, client := range c.Keys {
		ids[i] = fmt.Sprintf("%v", client.Kid)
	}
	return ids
}
