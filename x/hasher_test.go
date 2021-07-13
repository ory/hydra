package x

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type hasherConfig struct {
	cost int
}

func (c hasherConfig) BCryptCost() int {
	return c.cost
}

func TestHasher(t *testing.T) {
	for _, cost := range []int{1, 8, 10, 16} {
		result, err := NewBCrypt(&hasherConfig{cost: cost}).Hash(context.Background(), []byte("foobar"))
		require.NoError(t, err)
		require.NotEmpty(t, result)
	}
}

func BenchmarkHasher(b *testing.B) {
	for cost := 1; cost <= 16; cost++ {
		b.Run(fmt.Sprintf("cost=%d", cost), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				result, err := NewBCrypt(&hasherConfig{cost: cost}).Hash(context.Background(), []byte("foobar"))
				require.NoError(b, err)
				require.NotEmpty(b, result)
			}
		})
	}
}
