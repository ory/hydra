// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"encoding/binary"

	"github.com/pkg/errors"
)

// IntToBytes converts an int64 to a byte slice. It is the inverse of BytesToInt.
func IntToBytes(i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i)) //nolint:gosec

	return b
}

// BytesToInt converts a byte slice to an int64. It is the inverse of IntToBytes.
func BytesToInt(b []byte) (int64, error) {
	if len(b) != 8 {
		return 0, errors.New("byte slice must be 8 bytes long")
	}
	return int64(binary.LittleEndian.Uint64(b)), nil //nolint:gosec
}
