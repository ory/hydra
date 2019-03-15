package driver

import (
	"github.com/sirupsen/logrus"

	"github.com/ory/hive-cloud/hive/driver/configuration"
)

type DefaultDriver struct {
	l logrus.FieldLogger
	c configuration.Provider
	r Registry
}

func NewDefaultDriver(bi *BuildInfo, r Registry, c configuration.Provider, l logrus.FieldLogger) Driver {
	configuration.MustValidate(c)
	return NewDefaultDriverWithoutValidation(bi, r, c, l)
}

func NewDefaultDriverWithoutValidation(bi *BuildInfo, r Registry, c configuration.Provider, l logrus.FieldLogger) Driver {
	return &DefaultDriver{
		r: r.WithConfig(c),
		l: l,
		c: c,
	}
}

func (r *DefaultDriver) BuildInfo() *BuildInfo {
	return &BuildInfo{}
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
