package configuration

import (
	"github.com/ory/x/stringslice"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

type ViperProvider struct {
	l  logrus.FieldLogger
	ss [][]byte
}

func NewViperProvider(l logrus.FieldLogger) Provider {
	return &ViperProvider{
		l: l,
	}
}

func (v *ViperProvider) GetSubjectTypesSupported() []string {
	types := viper.GetStringSlice("OIDC_SUBJECT_TYPES_SUPPORTED")
	if len(types) == 0 {
		return []string{"public"}
	}

	return stringslice.Filter(types, func(s string) bool {
		return !(s == "public" || s == "pairwise")
	})
}

func (v *ViperProvider) DefaultClientScope() []string {
	return viper.GetStringSlice("OIDC_DYNAMIC_CLIENT_REGISTRATION_DEFAULT_SCOPE")
}
