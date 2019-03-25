package driver

import (
	"github.com/sirupsen/logrus"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/jwk"
	"github.com/ory/x/dbal"
)

type RegistryMemory struct {
	*RegistryBase
}

var _ Registry = new(RegistryMemory)

func init() {
	dbal.RegisterDriver(&RegistryMemory{
		RegistryBase: new(RegistryBase),
	})
}

func (m *RegistryMemory) Init(url string, l logrus.FieldLogger, opts ...dbal.DriverOptionModifier) error {
	m.l = l
	return nil
}

func (m *RegistryMemory) WithBuildVersion(bv string) Registry {
	m.buildVersion = bv
	return m
}

func (m *RegistryMemory) WithConfig(c configuration.Provider) Registry {
	m.c = c
	return m
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

func (m *RegistryMemory) KeyManager() jwk.Manager {
	if m.km == nil {
		m.km = jwk.NewMemoryManager()
	}
	return m.km
}
