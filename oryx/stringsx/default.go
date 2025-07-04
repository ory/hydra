// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package stringsx

import "cmp"

// Deprecated: use cmp.Or instead
func DefaultIfEmpty(s string, defaultValue string) string {
	return cmp.Or(s, defaultValue)
}
