package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatsdRegex(t *testing.T) {
	regx = statsdSanitizerRegex()
	resource := "rn:hydra.warden-token_allowed"
	sanitizedResource := regx.ReplaceAllString(resource, "_")
	assert.NotEqual(t, resource, sanitizedResource, "regex does not behave as expected")
	assert.Equal(t, "rn_hydra.warden-token_allowed", sanitizedResource, "regex does not behave as expected")
}
