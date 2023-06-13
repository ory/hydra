// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package aead

import (
	"context"
	"fmt"
)

func encryptionKey(ctx context.Context, d Dependencies, keySize int) ([]byte, error) {
	keys, err := allKeys(ctx, d)
	if err != nil {
		return nil, err
	}

	key := keys[0]
	if len(key) != keySize {
		return nil, fmt.Errorf("key must be exactly %d bytes long, got %d bytes", keySize, len(key))
	}

	return key, nil
}

func allKeys(ctx context.Context, d Dependencies) ([][]byte, error) {
	global, err := d.GetGlobalSecret(ctx)
	if err != nil {
		return nil, err
	}

	rotated, err := d.GetRotatedGlobalSecrets(ctx)
	if err != nil {
		return nil, err
	}

	keys := append([][]byte{global}, rotated...)
	if len(keys) == 0 {
		return nil, fmt.Errorf("at least one encryption key must be defined but none were")
	}
	return keys, nil
}
