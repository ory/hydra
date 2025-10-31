// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"github.com/ory/hydra/v2/fosite/storage"
)

func FositeStore() *storage.MemoryStore {
	return storage.NewMemoryStore()
}
