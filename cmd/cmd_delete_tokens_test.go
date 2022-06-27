package cmd_test

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/ory/hydra/client"
	"strings"
	"testing"

	"github.com/ory/hydra/cmd"
	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestDeleteAccessTokensCmd(t *testing.T) {
	c := cmd.NewDeleteAccessTokensCmd(new(cobra.Command))

	reg := setup(t, c)

	expected := createClient(t, reg, &client.Client{
		GrantTypes:              []string{"client_credentials"},
		TokenEndpointAuthMethod: "client_secret_post",
		Secret:                  uuid.Must(uuid.NewV4()).String()},
	)

	t.Run("case=deletes tokens", func(t *testing.T) {
		stdout := cmdx.ExecNoErr(t, c, expected.GetID())
		assert.Equal(t, fmt.Sprintf(`"%s"`, expected.GetID()), strings.TrimSpace(stdout))
	})
}
