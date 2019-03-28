package main

import (
	"github.com/ory/hydra/driver"
)

type MemTestPlugin struct {
	*driver.RegistryMemory
}

var Registry driver.Registry

func init() {
	Registry = &MemTestPlugin{RegistryMemory: driver.NewRegistryMemory()}
}

func main() {}
