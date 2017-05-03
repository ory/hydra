package cli

import (
	"testing"

	"github.com/ory/hydra/config"
)

func TestNewHandler(t *testing.T) {
	_ = NewHandler(&config.Config{})
}
