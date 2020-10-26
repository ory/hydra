package clients

import (
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFullCmdPath(t *testing.T) {
	t.Run("gets path from nested commands", func(t *testing.T) {
		expectedNames := []string{"a", "b", "c", "d", "e"}
		cmds := make([]*cobra.Command, len(expectedNames))

		for i, n := range expectedNames {
			cmds[i] = &cobra.Command{Use: n}
			if i > 0 {
				cmds[i-1].AddCommand(cmds[i])
			}
		}

		actualNames := getFullCmdPath(cmds[len(cmds)-1])

		assert.Equal(t, expectedNames, actualNames)
	})
}
