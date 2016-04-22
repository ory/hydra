package integration_test

import (
	"testing"
	"os"
	"github.com/ory-am/hydra/handler"
)

var h handler.OAuth2Handler

func TestMain(m *testing.M) {
	retCode := m.Run()
	os.Exit(retCode)
}

