package provider

import (
	"github.com/RangelReale/osin"
	"golang.org/x/oauth2"
)

type Provider interface {
	GetAuthCodeURL(ar *osin.AuthorizeRequest) string
	Exchange(code string) (Session, error)
	GetID() string
}

type Session interface {
	GetSubject() (string, error)
	GetToken() *oauth2.Token
}
