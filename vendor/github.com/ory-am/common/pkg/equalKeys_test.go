package pkg_test

import "testing"
import . "github.com/ory-am/common/pkg"

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
