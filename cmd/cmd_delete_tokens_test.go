package cmd_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/cmd"
	"github.com/ory/x/cmdx"
)

func TestDeleteAccessTokensCmd(t *testing.T) {
	c := cmd.NewDeleteAccessTokensCmd(new(cobra.Command))

	reg := setup(t, c)
	expected := createClientCredentialsClient(t, reg)
	t.Run("case=deletes tokens", func(t *testing.T) {
		stdout := cmdx.ExecNoErr(t, c, expected.GetID())
		assert.Equal(t, fmt.Sprintf(`"%s"`, expected.GetID()), strings.TrimSpace(stdout))
	})
}
