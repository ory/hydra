package driver

import (
	"github.com/sirupsen/logrus"

	"github.com/ory/hydra/driver/configuration"
)

type DefaultDriver struct {
	l logrus.FieldLogger
	c configuration.Provider
	r Registry
}

func NewDefaultDriver( r Registry, c configuration.Provider, l logrus.FieldLogger) Driver {
	configuration.MustValidate(c)
	return NewDefaultDriverWithoutValidation( r, c, l)
}

func NewDefaultDriverWithoutValidation( r Registry, c configuration.Provider, l logrus.FieldLogger) Driver {
	return &DefaultDriver{
		r: r.WithConfig(c),
		l: l,
		c: c,
	}
}

func (r *DefaultDriver) Logger() logrus.FieldLogger {
	return r.l
}

func (r *DefaultDriver) Configuration() configuration.Provider {
	return r.c
}

func (r *DefaultDriver) Registry() Registry {
	return r.r
}
