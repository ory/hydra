package provider_test

import (
	. "github.com/ory-am/hydra/endpoint/connector"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultSession(t *testing.T) {
	s := DefaultSession{RemoteSubject: "subject", ForceLocalSubject: "subject", Extra: map[string]interface{}{"extra": "foo"}}
	assert.Equal(t, "subject", s.GetRemoteSubject())
	assert.Equal(t, "subject", s.GetForcedLocalSubject())
	assert.Equal(t, "foo", s.GetExtra()["extra"])
}
