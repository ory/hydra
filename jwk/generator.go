// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk

import jose "github.com/go-jose/go-jose/v3"

type KeyGenerator interface {
	Generate(id, use string) (*jose.JSONWebKeySet, error)
}
