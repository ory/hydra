package driver

import (
	"github.com/ory/fosite"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/x/dbal"
)

type RegistryMemory struct {
	*RegistryBase
}

var _ Registry = new(RegistryMemory)

func init() {
	dbal.RegisterDriver(func() dbal.Driver {
		return NewRegistryMemory()
	})
}

func NewRegistryMemory() *RegistryMemory {
	r := &RegistryMemory{
		RegistryBase: new(RegistryBase),
	}
	r.RegistryBase.with(r)
	return r
}

// WithOAuth2Provider forces an oauth2 provider which is only used for testing.
func (m *RegistryMemory) WithOAuth2Provider(f fosite.OAuth2Provider) *RegistryMemory {
	m.RegistryBase.fop = f
	return m
}

// WithConsentStrategy forces a consent strategy which is only used for testing.
func (m *RegistryMemory) WithConsentStrategy(c consent.Strategy) *RegistryMemory {
	m.RegistryBase.cos = c
	return m
}

func (m *RegistryMemory) Init() error {
	return nil
}

func (m *RegistryMemory) CanHandle(dsn string) bool {
	return dsn == "memory"
}

func (m *RegistryMemory) Ping() error {
	return nil
}

func (m *RegistryMemory) ClientManager() client.Manager {
	if m.cm == nil {
		m.cm = client.NewMemoryManager(m)
	}
	return m.cm
}

func (m *RegistryMemory) ConsentManager() consent.Manager {
	if m.com == nil {
		m.com = consent.NewMemoryManager(m)
	}
	return m.com
}

func (m *RegistryMemory) OAuth2Storage() x.FositeStorer {
	if m.fs == nil {
		m.fs = oauth2.NewFositeMemoryStore(m.r, m.c)
	}
	return m.fs
}

func (m *RegistryMemory) KeyManager() jwk.Manager {
	if m.km == nil {
		m.km = jwk.NewMemoryManager()
	}
	return m.km
}
