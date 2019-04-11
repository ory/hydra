package driver

import (
	"github.com/ory/hydra/driver/configuration"
)

type Driver interface {
	Configuration() configuration.Provider
	Registry() Registry
	CallRegistry() Driver
}
