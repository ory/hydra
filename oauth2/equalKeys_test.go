// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import "testing"

func TestAssertObjectsAreEqualByKeys(t *testing.T) {
	type foo struct {
		Name string
		Body int
	}
	a := &foo{"foo", 1}
	b := &foo{"bar", 1}
	c := &foo{"baz", 3}

	AssertObjectKeysEqual(t, a, a, "Name", "Body")
	AssertObjectKeysNotEqual(t, a, b, "Name")
	AssertObjectKeysNotEqual(t, a, c, "Name", "Body")
}
