package account

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountCases(t *testing.T) {
	for _, c := range [][]string{
		{"1", "foo@bar", "secret"},
		{"2", "baz@foo", "top secret"},
	} {
		a := &DefaultAccount{c[0], c[1], c[2], `{"foo": "bar"}`}
		assert.Equal(t, c[0], a.GetID())
		assert.Equal(t, c[1], a.GetEmail())
		assert.Equal(t, c[2], a.GetPassword())
		assert.Equal(t, `{"foo": "bar"}`, a.GetData())
	}
}
