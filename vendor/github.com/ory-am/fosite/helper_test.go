package fosite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	for k, c := range []struct {
		needle   string
		haystack []string
		ok       bool
	}{
		{needle: "foo", haystack: []string{"foo", "bar"}, ok: true},
		{needle: "bar", haystack: []string{"foo", "bar"}, ok: true},
		{needle: "baz", haystack: []string{"foo", "bar"}, ok: false},
		{needle: "foo", haystack: []string{"bar"}, ok: false},
		{needle: "bar", haystack: []string{"bar"}, ok: true},
		{needle: "foo", haystack: []string{}, ok: false},
	} {
		assert.Equal(t, c.ok, StringInSlice(c.needle, c.haystack), "%d", k)
		t.Logf("Passed test case %d", k)
	}
}
