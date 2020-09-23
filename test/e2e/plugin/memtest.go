package main

import (
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/viper"
	"github.com/ory/x/logrusx"
)

type MemTestPlugin struct {
	*driver.RegistrySQL
}

func NewRegistry() driver.Registry {
	viper.Set(configuration.ViperKeyDSN, configuration.DefaultSQLiteMemoryDSN)
	r := driver.NewRegistrySQL()
	r.WithLogger(logrusx.New("Hydra plugin registry", "test"))
	return &MemTestPlugin{RegistrySQL: r}
}

func main() {}
