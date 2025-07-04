// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonnetsecure

import "bytes"

func splitNull(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Look for a null byte; if found, return the position after it,
	// the data before it, and no error.
	if i := bytes.IndexByte(data, 0); i >= 0 {
		return i + 1, data[0:i], nil
	}

	// If we're at EOF, we have a final, non-terminated word. Return it.
	if atEOF && len(data) != 0 {
		return len(data), data, nil
	}

	// Request more data.
	return 0, nil, nil
}
