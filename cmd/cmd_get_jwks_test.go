package cmd_test

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/hydra/cmd"
	"github.com/ory/x/cmdx"
)

func TestGetJwks(t *testing.T) {
	ctx := context.Background()
	c := cmd.NewGetJWKSCmd()
	reg := setup(t, c)

	set := uuid.Must(uuid.NewV4()).String()
	_ = createJWK(t, reg, set, "RS256")

	t.Run("case=gets jwks", func(t *testing.T) {
		actual := gjson.Parse(cmdx.ExecNoErr(t, c, set))
		assert.NotEmpty(t, actual.Get("kid").String(), actual.Raw)

		expected, err := reg.KeyManager().GetKeySet(ctx, set)
		require.NoError(t, err)

		assert.Equal(t, expected.Keys[0].KeyID, actual.Get("kid").String())
	})
}
