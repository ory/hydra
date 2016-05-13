package config

import (
	"github.com/ory-am/hydra/warden"
	"github.com/ory-am/fosite/hash"
)

type Context struct {
	Connection interface{}

	Hasher     hash.Hasher
	Warden     *warden.LocalWarden
}
