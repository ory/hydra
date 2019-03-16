package configuration

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func setEnv(key, value string) func(t *testing.T) {
	return func(t *testing.T) {
		require.NoError(t, os.Setenv(key, value))
	}
}

func TestSubjectTypesSupported(t *testing.T) {
	p := new(ViperProvider)
	for k, tc := range []struct {
		d string
		p func(t *testing.T)
		e []string
		c func(t *testing.T)
	}{
		{
			d: "Load legacy environment variable in legacy format",
			p: setEnv("OIDC_SUPPORTED_SUBJECT_TYPES", "pairwise,public,foobar"),
			c: setEnv("OIDC_SUPPORTED_SUBJECT_TYPES", ""),
			e: []string{"public", "pairwise"},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s",k,tc.d), func(t *testing.T) {
			assert.EqualValues(t, tc.e, p.SubjectTypesSupported())
		})
	}
}
