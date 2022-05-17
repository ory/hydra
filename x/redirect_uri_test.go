package x

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockrc struct{}

func (m *mockrc) InsecureRedirects(ctx context.Context) []string {
	return []string{
		"http://foo.com/bar",
		"http://baz.com/bar",
	}
}

func TestIsRedirectURISecure(t *testing.T) {
	for d, c := range []struct {
		u   string
		err bool
	}{
		{u: "http://google.com", err: true},
		{u: "https://google.com", err: false},
		{u: "http://localhost", err: false},
		{u: "http://test.localhost", err: false},
		{u: "wta://auth", err: false},
		{u: "http://foo.com/bar", err: false},
		{u: "http://baz.com/bar", err: false},
	} {
		uu, err := url.Parse(c.u)
		require.NoError(t, err)
		assert.Equal(t, !c.err, IsRedirectURISecure(new(mockrc))(context.Background(), uu), "case %d", d)
	}
}
