package cmd

import (
	"testing"

	"github.com/ory/x/cmdx"
)

func TestUsageStrings(t *testing.T) {
	cmdx.AssertUsageTemplates(t, NewRootCmd(nil, nil, nil))
}
