package x

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRouterAdminAdmin(t *testing.T) {
	require.NotEmpty(t, NewRouterAdmin())
	require.NotEmpty(t, NewRouterPublic())
}
