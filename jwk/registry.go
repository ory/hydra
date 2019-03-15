package jwk

import "github.com/ory/herodot"

type Registry interface {
	JWKManager() Manager
	JWKGenerators()    map[string]KeyGenerator
	Writer() herodot.Writer
	Cipher() *AEAD
}

type Configuration interface {
	WellKnownKeys(defaults ...string) []string
}
