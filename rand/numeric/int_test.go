package numeric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt64(t *testing.T) {
	seq := Int64()
	assert.NotEmpty(t, seq)
}

func TestInt32(t *testing.T) {
	seq := Int32()
	assert.NotEmpty(t, seq)
}

func TestUInt64(t *testing.T) {
	seq := UInt64()
	assert.NotEmpty(t, seq)
}

func TestUInt32(t *testing.T) {
	seq := UInt32()
	assert.NotEmpty(t, seq)
}

func TestInt64IsUnique(t *testing.T) {
	// Probability of collision is around 1 in a million
	times := 6000000
	s := make(map[int64]bool)

	for i := 0; i < times; i++ {
		k := Int64()
		_, ok := s[k]
		assert.False(t, ok)
		if ok {
			return
		}
		s[k] = true
	}
}

func TestUInt64IsUnique(t *testing.T) {
	// Probability of collision is around 1 in a million
	times := 6000000
	s := make(map[uint64]bool)

	for i := 0; i < times; i++ {
		k := UInt64()
		_, ok := s[k]
		assert.False(t, ok)
		if ok {
			return
		}
		s[k] = true
	}
}

func TestInt32IsUnique(t *testing.T) {
	// Probability of collision is around 1 in 1000
	times := 3000
	s := make(map[int32]bool)

	for i := 0; i < times; i++ {
		k := Int32()
		_, ok := s[k]
		assert.False(t, ok)
		if ok {
			return
		}
		s[k] = true
	}
}

func TestUInt32IsUnique(t *testing.T) {
	// Probability of collision is around 1 in 1000
	times := 3000
	s := make(map[uint32]bool)

	for i := 0; i < times; i++ {
		k := UInt32()
		_, ok := s[k]
		assert.False(t, ok)
		if ok {
			return
		}
		s[k] = true
	}
}

func BenchmarkTestInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Int64()
	}
}
