package hydra

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestInterface(t *testing.T) {
	var sdk SDK
	var err error
	sdk, err = NewSDK(&Configuration{
		EndpointURL:  "http://localhost:4444/",
		ClientID:     "foo",
		ClientSecret: "bar",
	})
	assert.NoError(t, err)
	assert.NotNil(t, sdk)
}

func TestErrorHandlers(t *testing.T) {
	for k, c := range []Configuration{
		{
			EndpointURL:  "http://localhost:4444/",
			ClientSecret: "bar",
			Scopes:       []string{"foo"},
		},
		{
			EndpointURL: "http://localhost:4444/",
			ClientID:    "bar",
			Scopes:      []string{"foo"},
		},
		{
			ClientID:     "foo",
			ClientSecret: "bar",
			Scopes:       []string{"foo"},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			sdk, err := NewSDK(&c)
			assert.Error(t, err)
			assert.Nil(t, sdk)
		})
	}
}
