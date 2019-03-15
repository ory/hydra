package configuration

import (
	"net/url"
)

type Provider interface {
	AdminListenOn() string
	PublicListenOn() string
	DSN() string

	PublicURL() *url.URL
	AdminURL() *url.URL

	BCryptCostFactor() int


	DefaultClientScope() []string
	GetSubjectTypesSupported() []string
}

func MustValidate(p Provider) {
}
