package driver

import (
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/x"
	"github.com/ory/x/dbal"
	"github.com/sirupsen/logrus"
)

type RegistryMemory struct {
	l      logrus.FieldLogger
	c      configuration.Provider
	cm     client.Manager
	ch *client.Handler
	fh     fosite.Hasher
	cv     *client.Validator
	kg     map[string]jwk.KeyGenerator
	km     jwk.Manager
	kc     *jwk.AEAD
	writer herodot.Writer
}

var _ Registry = new(RegistryMemory)

func init() {
	dbal.RegisterDriver(new(RegistryMemory))
}

func (m *RegistryMemory) Init(url string, l logrus.FieldLogger, opts ...dbal.DriverOptionModifier) error {
	m.l = l
	return nil
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

func (m *RegistryMemory) Writer() herodot.Writer {
	if m.writer == nil {
		m.writer = herodot.NewJSONWriter(m.Logger())
	}
	return m.writer
}

func (m *RegistryMemory) Logger() logrus.FieldLogger {
	if m.l == nil {
		m.l = logrus.New()
	}
	return m.l
}

func (m *RegistryMemory) ClientManager() client.Manager {
	if m.cm == nil {
		m.cm = client.NewMemoryManager(m)
	}
	return m.cm
}

func (m *RegistryMemory) ClientHasher() fosite.Hasher {
	if m.fh == nil {
		m.fh = x.NewBCrypt(m.c)
	}
	return m.fh
}

func (m *RegistryMemory) ClientHandler() *client.Handler {
	if m.ch == nil {
		m.ch = client.NewHandler(m)
	}
	return m.ch
}

func (m *RegistryMemory) ClientValidator() *client.Validator {
	if m.cv == nil {
		m.cv = client.NewValidator(m.c)
	}
	return m.cv
}

func (m *RegistryMemory) KeyManager() jwk.Manager {
	if m.km == nil {
		m.km = jwk.NewMemoryManager()
	}
	return m.km
}

func (m *RegistryMemory) KeyGenerators() map[string]jwk.KeyGenerator {
	if m.kg == nil {
		m.kg = map[string]jwk.KeyGenerator{
			"RS256": &jwk.RS256Generator{},
			"ES512": &jwk.ECDSA512Generator{},
			"HS256": &jwk.HS256Generator{},
			"HS512": &jwk.HS512Generator{},
		}
	}
	return m.kg
}

func (m *RegistryMemory) KeyCipher() *jwk.AEAD {
	if m.kc == nil {
		m.kc = jwk.NewAEAD(m.c.GetSystemSecret())
	}
	return m.kc
}

