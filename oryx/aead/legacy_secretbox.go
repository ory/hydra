// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package aead

import "golang.org/x/crypto/nacl/secretbox"

// OpenLegacySecretbox opens a payload sealed by the NaCl secretbox scheme
// that was used before this package existed: a 24-byte nonce prefixed to the
// ciphertext, without additional data. Remove it once all such payloads have
// expired.
func OpenLegacySecretbox(key [32]byte, raw []byte) ([]byte, bool) {
	if len(raw) < 24 {
		return nil, false
	}
	var nonce [24]byte
	copy(nonce[:], raw[:24])
	return secretbox.Open(nil, raw[24:], &nonce, &key)
}
