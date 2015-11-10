package provider

import (
	"golang.org/x/oauth2"
)

type Session interface {
	GetRemoteSubject() string
	GetExtra() interface{}
	GetToken() *oauth2.Token
}

type DefaultSession struct {
	RemoteSubject string
	Extra         interface{}
	Token         *oauth2.Token
}

func (s *DefaultSession) GetRemoteSubject() string {
	return s.RemoteSubject
}

func (s *DefaultSession) GetExtra() interface{} {
	return s.Extra
}

func (s *DefaultSession) GetToken() *oauth2.Token {
	return s.Token
}
