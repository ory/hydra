// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flagx

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestStringToStringCommand(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringToString("map-value", nil, "test string to string map usage")

	cmd.SetArgs([]string{"--map-value", "foo=bar,key=val"})
	cmd.Execute()

	mapped := MustGetStringToStringMap(cmd, "map-value")

	if len(mapped) != 2 {
		t.Errorf("expected 2 values in map and got %d", len(mapped))
	}
	val, ok := mapped["foo"]
	if !ok {
		t.Errorf("failed to get value 'foo' from flags")
	}
	if val != "bar" {
		t.Errorf("failed to get expected value from map, got %s", val)
	}
}
