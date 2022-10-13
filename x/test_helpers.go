// Copyright © 2022 Ory Corp

package x

import (
	"github.com/ory/fosite/storage"
)

func FositeStore() *storage.MemoryStore {
	return storage.NewMemoryStore()
}
