package main

import (
	"github.com/ory/hydra/driver"
)

type MemTestPlugin struct {
	*driver.RegistryMemory
}

func NewRegistry() driver.Registry {
	return &MemTestPlugin{RegistryMemory: driver.NewRegistryMemory()}
}

func main() {}
