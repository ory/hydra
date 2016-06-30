package oauth2

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/ory-am/fosite"
)

func TestToStringSlice(t *testing.T) {
	assert.Equal(t, []string{"foo"}, toStringSlice((map[string]interface{}{
		"scp": fosite.Arguments{"foo"},
	})["scp"]))
	assert.Equal(t, []string{"foo"}, toStringSlice((map[string]interface{}{
		"scp": []string{"foo"},
	})["scp"]))
	assert.Equal(t, []string{"foo"}, toStringSlice((map[string]interface{}{
		"scp": []interface{}{"foo", 123},
	})["scp"]))
}