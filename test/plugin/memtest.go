package main

import (
	"github.com/ory/hydra/config"
)

type MemTestPlugin struct {
	config.MemoryBackend
}

func (m *MemTestPlugin) Prefixes() []string {
	return []string{"memtest"}
}

var BackendConnector config.BackendConnector = &MemTestPlugin{}
