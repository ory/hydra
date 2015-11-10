package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStateData(t *testing.T) {
	sd := &StateData{
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	assert.True(t, sd.IsExpired())

	sd = &StateData{
		ExpiresAt: time.Now().Add(time.Hour),
	}
	assert.False(t, sd.IsExpired())
}
