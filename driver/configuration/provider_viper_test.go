package configuration

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func setEnv(key, value string) func(t *testing.T) {
	return func(t *testing.T) {
		require.NoError(t, os.Setenv(key, value))
	}
}

func TestSubjectTypesSupported(t *testing.T) {
	p := NewViperProvider(logrus.New(), false)
	viper.Set(ViperKeySubjectIdentifierAlgorithmSalt, "00000000")
	for k, tc := range []struct {
		d string
		p func(t *testing.T)
		e []string
		c func(t *testing.T)
	}{
		{
			d: "Load legacy environment variable in legacy format",
			p: setEnv(strings.ToUpper(strings.Replace(ViperKeySubjectTypesSupported, ".", "_", -1)), "public,pairwise,foobar"),
			c: setEnv(strings.ToUpper(strings.Replace(ViperKeySubjectTypesSupported, ".", "_", -1)), ""),
			e: []string{"public", "pairwise"},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
			tc.p(t)
			assert.EqualValues(t, tc.e, p.SubjectTypesSupported())
			tc.c(t)
		})
	}
}
