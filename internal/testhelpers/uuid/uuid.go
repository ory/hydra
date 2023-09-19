// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package uuid

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

// AssertUUID helper requires that a UUID is non-zero, common version/variant used in Hydra.
func AssertUUID[T string | uuid.UUID](t *testing.T, id T) {
	var uid uuid.UUID
	switch idt := any(id).(type) {
	case uuid.UUID:
		uid = idt
	case string:
		var err error
		uid, err = uuid.FromString(idt)
		require.NoError(t, err)
	}
	require.Equal(t, uid.Version(), uuid.V4)
	require.Equal(t, uid.Variant(), uuid.VariantRFC4122)
}
