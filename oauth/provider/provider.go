package provider

import (
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/RangelReale/osin"
	"github.com/ory-am/hydra/Godeps/_workspace/src/golang.org/x/oauth2"
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
	GetAuthCodeURL(ar *osin.AuthorizeRequest) string
	Exchange(code string) (Session, error)
	GetID() string
}

func GetAuthCodeURL(conf oauth2.Config, ar *osin.AuthorizeRequest, provider string) string {
	redirect, err := url.Parse(conf.RedirectURL)
	if err != nil {
		return ""
	}

	q := redirect.Query()
	q.Set(ProviderQueryParam, provider)
	q.Set(RedirectQueryParam, ar.RedirectUri)
	q.Set(ClientQueryParam, ar.Client.GetId())
	q.Set(ScopeQueryParam, ar.Scope)
	q.Set(StateQueryParam, ar.State)
	q.Set(TypeQueryParam, string(ar.Type))
	redirect.RawQuery = q.Encode()

	conf.RedirectURL = redirect.String()
	return conf.AuthCodeURL(ar.State)
}
