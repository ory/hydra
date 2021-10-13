package uuid

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

func AssertUUID(t *testing.T, id *uuid.UUID) {
	require.Equal(t, id.Version(), uuid.V4)
	require.Equal(t, id.Variant(), uuid.VariantRFC4122)
}
