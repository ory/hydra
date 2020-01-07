// +build !noplugin

package driver

import (
	"net/url"
	"plugin"
	"strings"

	"github.com/pkg/errors"

	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/x/dbal"
)

type RegistryPlugin struct {
	c configuration.Provider
	Registry
}

var _ Registry = new(RegistryPlugin)

func init() {
	dbal.RegisterDriver(func() dbal.Driver {
		return NewRegistryPlugin()
	})
}

func NewRegistryPlugin() *RegistryPlugin {
	return new(RegistryPlugin)
}

func (m *RegistryPlugin) CanHandle(dsn string) bool {
	u, err := url.Parse(dsn)
	if err != nil {
		return false
	}
	return u.Scheme == "plugin"
}

func (m *RegistryPlugin) WithConfig(c configuration.Provider) Registry {
	m.c = c
	return m
}

func (m *RegistryPlugin) Init() error {
	if m.Registry != nil {
		return nil
	}

	path := strings.Replace(m.c.DSN(), "plugin://", "", 1)
	p, err := plugin.Open(path)
	if err != nil {
		return errors.Wrapf(err, "unable to open plugin path: %s", path)
	}

	l, err := p.Lookup("NewRegistry")
	if err != nil {
		return errors.Wrap(err, "unable to look up `Registry`")
	}

	reg, ok := l.(func() Registry)
	if !ok {
		return errors.Errorf("unable to type assert %T to `func() driver.Registry`", l)
	}

	m.Registry = reg()
	m.Logger().Info("Successfully loaded database plugin")
	m.Logger().Debugf("Memory address of database plugin is: %p", reg)
	m.Registry.WithConfig(m.c)

	return nil
}
