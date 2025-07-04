// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package uuidx

import "github.com/gofrs/uuid"

// NewV4 returns a new randomly generated UUID or panics.
func NewV4() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}
