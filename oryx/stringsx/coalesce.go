// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package stringsx

import "cmp"

// Coalesce returns the first non-empty string value
// Deprecated: use cmp.Or instead
func Coalesce(str ...string) string {
	return cmp.Or(str...)
}
