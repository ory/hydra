// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/client"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/ory/hydra/v2/cmd"
	"github.com/ory/x/cmdx"
)

func TestListClient(t *testing.T) {
	c := cmd.NewListClientsCmd()
	reg := setup(t, c)

	expected1 := createClient(t, reg, nil)
	expected2 := createClient(t, reg, nil)

	t.Run("case=lists both clients", func(t *testing.T) {
		actual := gjson.Parse(cmdx.ExecNoErr(t, c))
		assert.Len(t, actual.Get("items").Array(), 2)

		for _, e := range []*client.Client{expected1, expected2} {
			assert.Contains(t, actual.Raw, e.GetID())
		}
	})

	t.Run("case=lists both clients with pagination", func(t *testing.T) {
		actualFirst := gjson.Parse(cmdx.ExecNoErr(t, c, "--format", "json", "--page-size", "1"))
		assert.Len(t, actualFirst.Get("items").Array(), 1)

		require.NotEmpty(t, actualFirst.Get("next_page_token").String(), actualFirst.Raw)
		assert.False(t, actualFirst.Get("is_last_page").Bool(), actualFirst.Raw)

		actualSecond := gjson.Parse(cmdx.ExecNoErr(t, c, "--format", "json", "--page-size", "1", "--page-token", actualFirst.Get("next_page_token").String()))
		assert.Len(t, actualSecond.Array(), 1)

		assert.NotEmpty(t, actualFirst.Get("items.0.client_id").String())
		assert.NotEqualValues(t, actualFirst.Get("items.0.client_id").String(), actualSecond.Get("items.0.client_id").String())
	})
}
