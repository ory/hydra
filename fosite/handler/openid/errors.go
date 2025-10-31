// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openid

import "github.com/pkg/errors"

var (
	ErrInvalidSession = errors.New("Session type mismatch")
)
