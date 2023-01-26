// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/ory/hydra/v2/x"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/v2/cmd"
	"github.com/ory/x/assertx"
	"github.com/ory/x/cmdx"
)

func TestDeleteJwks(t *testing.T) {
	ctx := context.Background()
	c := cmd.NewDeleteJWKSCommand()
	reg := setup(t, c)

	t.Run("case=deletes jwks", func(t *testing.T) {
		set := uuid.Must(uuid.NewV4()).String()
		_ = createJWK(t, reg, set, "RS256")
		stdout := cmdx.ExecNoErr(t, c, set)
		assert.Equal(t, fmt.Sprintf(`"%s"`, set), strings.TrimSpace(stdout))

		_, err := reg.KeyManager().GetKeySet(ctx, set)
		assert.ErrorIs(t, err, x.ErrNotFound)
	})

	t.Run("case=deletes multiple jwkss", func(t *testing.T) {
		set1 := uuid.Must(uuid.NewV4()).String()
		set2 := uuid.Must(uuid.NewV4()).String()
		_ = createJWK(t, reg, set1, "RS256")
		_ = createJWK(t, reg, set2, "RS256")
		assertx.EqualAsJSON(t, []string{set1, set2}, json.RawMessage(cmdx.ExecNoErr(t, c, set1, set2)))

		_, err := reg.KeyManager().GetKeySet(ctx, set1)
		assert.ErrorIs(t, err, x.ErrNotFound)

		_, err = reg.KeyManager().GetKeySet(ctx, set2)
		assert.ErrorIs(t, err, x.ErrNotFound)
	})
}
