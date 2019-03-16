package driver

import (
	"github.com/sirupsen/logrus"

	"github.com/ory/hydra/driver/configuration"
)

type Driver interface {
	Logger() logrus.FieldLogger
	Configuration() configuration.Provider
	Registry() Registry
}
