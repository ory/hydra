package driver

import (
	"github.com/go-errors/errors"
	"github.com/sirupsen/logrus"

	"github.com/ory/herodot"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/x/dbal"
)

type Registry interface {
	dbal.Driver

	WithConfig(c configuration.Provider) Registry

	Logger() logrus.FieldLogger
	Writer() herodot.Writer
}

func NewRegistry(c configuration.Provider) (Registry, error) {
	driver, err := dbal.GetDriverFor(c.DSN())
	if err != nil {
		return nil, err
	}

	registry, ok := driver.(Registry)
	if !ok {
		return nil, errors.Errorf("driver of type %T does not implement interface Registry", driver)
	}

	return registry, nil
}
