// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package osx

import (
	"cmp"
	"os"
)

// GetenvDefault returns an environment variable or the default value if it is empty.
func GetenvDefault(key string, def string) string {
	return cmp.Or(os.Getenv(key), def)
}
