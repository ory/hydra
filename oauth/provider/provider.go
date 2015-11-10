package provider

import (
	"github.com/RangelReale/osin"
	"golang.org/x/oauth2"
	"net/url"
)

const (
	ProviderQueryParam = "opr"
	RedirectQueryParam = "ord"
	ClientQueryParam   = "ocl"
	ScopeQueryParam    = "osc"
	StateQueryParam    = "ost"
	TypeQueryParam     = "otp"
)

type Provider interface {
	GetAuthCodeURL(state string) (string)
	Exchange(code string) (Session, error)
	GetID() string
}

func GetAuthCodeURL(conf oauth2.Config, state string) (string) {
	return conf.AuthCodeURL(state)
}
