package hydra

import (
	"testing"
	"github.com/docker/docker/pkg/testutil/assert"
)

func TestInterface(t *testing.T) {
	var sdk SDK
	sdk = new(CodeGenSDK)
	assert.NotNil(t, sdk)
}
