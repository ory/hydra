// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package networkx

import (
	"time"

	"github.com/gofrs/uuid"
)

type Network struct {
	ID uuid.UUID `json:"id" db:"id"`

	// CreatedAt is a helper struct field for gobuffalo.pop.
	CreatedAt time.Time `json:"-" db:"created_at"`

	// UpdatedAt is a helper struct field for gobuffalo.pop.
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

func (p Network) TableName() string {
	return "networks"
}

func NewNetwork() *Network {
	return &Network{
		ID: uuid.Must(uuid.NewV4()),
	}
}
