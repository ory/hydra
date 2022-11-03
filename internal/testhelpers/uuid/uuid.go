// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package uuid

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

// AssertUUID helper requires that a UUID is non-zero, common version/variant used in Hydra.
func AssertUUID(t *testing.T, id *uuid.UUID) {
	require.Equal(t, id.Version(), uuid.V4)
	require.Equal(t, id.Variant(), uuid.VariantRFC4122)
}
