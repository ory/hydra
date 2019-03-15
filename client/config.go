package client

type Configuration interface {
	DefaultClientScope() []string
	GetSubjectTypesSupported() []string
}