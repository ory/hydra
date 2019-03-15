package client

import (
	"github.com/ory/fosite"
	"github.com/ory/herodot"
)

type Registry interface {
	ClientValidator() *Validator
	ClientManager() Manager
	Writer() herodot.Writer
	ClientHasher() fosite.Hasher
}
