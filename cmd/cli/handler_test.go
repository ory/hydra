package cli

import (
	"github.com/ory-am/hydra/config"
	"testing"
)

func TestNewHandler(t *testing.T) {
	_ = NewHandler(&config.Config{})
}
