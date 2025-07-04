// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package assertx

import (
	"testing"
	"time"
)

func TestEqualAsJSONExcept(t *testing.T) {
	a := map[string]interface{}{"foo": "bar", "baz": "bar", "bar": "baz"}
	b := map[string]interface{}{"foo": "bar", "baz": "bar", "bar": "not-baz"}

	EqualAsJSONExcept(t, a, b, []string{"bar"})
}

func TestTimeDifferenceLess(t *testing.T) {
	TimeDifferenceLess(t, time.Now(), time.Now().Add(time.Second), 2)
}
