package cli

import (
	"testing"
	"github.com/ory-am/hydra/config"
)

func TestNewHandler(t *testing.T) {
	var c = new(config.Config)
	_ = NewHandler(c)
}
