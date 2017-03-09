package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoinURL(t *testing.T) {
	for k, c := range []struct {
		give []string
		get  string
	}{
		{
			give: []string{"http://localhost/", "/home"},
			get:  "http://localhost/home",
		},
		{
			give: []string{"http://localhost", "/home"},
			get:  "http://localhost/home",
		},
		{
			give: []string{"https://localhost/", "/home"},
			get:  "https://localhost/home",
		},
		{
			give: []string{"http://localhost/", "/home", "home/", "/home/"},
			get:  "http://localhost/home/home/home/",
		},
	} {
		assert.Equal(t, c.get, JoinURL(c.give[0], c.give[1:]...), "Case %d", k)
	}
}
