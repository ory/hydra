package cli

import (
	"testing"

	"github.com/ory-am/hydra/config"
)

func TestNewHandler(t *testing.T) {
	_ = NewHandler(&config.Config{})
}
