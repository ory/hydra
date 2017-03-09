package fosite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgumentsExact(t *testing.T) {
	for k, c := range []struct {
		args   Arguments
		exact  string
		expect bool
	}{
		{
			args:   Arguments{"foo"},
			exact:  "foo",
			expect: true,
		},
		{
			args:   Arguments{"foo", "bar"},
			exact:  "foo",
			expect: false,
		},
		{
			args:   Arguments{"foo", "bar"},
			exact:  "bar",
			expect: false,
		},
		{
			args:   Arguments{"foo", "bar"},
			exact:  "baz",
			expect: false,
		},
		{
			args:   Arguments{},
			exact:  "baz",
			expect: false,
		},
	} {
		assert.Equal(t, c.expect, c.args.Exact(c.exact), "%d", k)
		t.Logf("Passed test case %d", k)
	}
}

func TestArgumentsHas(t *testing.T) {
	for k, c := range []struct {
		args   Arguments
		has    []string
		expect bool
	}{
		{
			args:   Arguments{"foo", "bar"},
			has:    []string{"foo", "bar"},
			expect: true,
		},
		{
			args:   Arguments{"foo", "bar"},
			has:    []string{"bar", "foo"},
			expect: true,
		},
		{
			args:   Arguments{"bar", "foo"},
			has:    []string{"foo"},
			expect: true,
		},
		{
			args:   Arguments{"foo", "bar"},
			has:    []string{"bar", "foo", "baz"},
			expect: false,
		},
		{
			args:   Arguments{"foo", "bar"},
			has:    []string{"foo"},
			expect: true,
		},
		{
			args:   Arguments{"foo", "bar"},
			has:    []string{"bar"},
			expect: true,
		},
		{
			args:   Arguments{"foo", "bar"},
			has:    []string{"baz"},
			expect: false,
		},
		{
			args:   Arguments{},
			has:    []string{"baz"},
			expect: false,
		},
	} {
		assert.Equal(t, c.expect, c.args.Has(c.has...), "%d", k)
		t.Logf("Passed test case %d", k)
	}
}

func TestArgumentsMatches(t *testing.T) {
	for k, c := range []struct {
		args   Arguments
		is     []string
		expect bool
	}{
		{
			args:   Arguments{"foo", "bar"},
			is:     []string{"foo", "bar"},
			expect: true,
		},
		{
			args:   Arguments{"foo", "foo"},
			is:     []string{"foo"},
			expect: false,
		},
		{
			args:   Arguments{"foo", "foo"},
			is:     []string{"bar", "foo"},
			expect: false,
		},
		{
			args:   Arguments{"foo", "bar"},
			is:     []string{"bar", "foo", "baz"},
			expect: false,
		},
		{
			args:   Arguments{"foo", "bar"},
			is:     []string{"foo"},
			expect: false,
		},
		{
			args:   Arguments{"foo", "bar"},
			is:     []string{"bar", "bar"},
			expect: false,
		},
		{
			args:   Arguments{"foo", "bar"},
			is:     []string{"baz"},
			expect: false,
		},
		{
			args:   Arguments{},
			is:     []string{"baz"},
			expect: false,
		},
	} {
		assert.Equal(t, c.expect, c.args.Matches(c.is...), "%d", k)
		t.Logf("Passed test case %d", k)
	}
}
