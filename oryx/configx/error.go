// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	"fmt"

	"github.com/pkg/errors"
)

type ImmutableError struct {
	From interface{}
	To   interface{}
	Key  string
	error
}

func NewImmutableError(key string, from, to interface{}) error {
	return &ImmutableError{
		From:  from,
		To:    to,
		Key:   key,
		error: errors.Errorf("immutable configuration key \"%s\" was changed from \"%v\" to \"%v\"", key, from, to),
	}
}

func (e *ImmutableError) Error() string {
	return fmt.Sprintf("immutable configuration key \"%s\" was changed from \"%v\" to \"%v\"", e.Key, e.From, e.To)
}
