package cmd

import (
	"github.com/ory/x/pagination/tokenpagination"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestFatal(t *testing.T) {
	oldOsExit := osExit
	defer func() { osExit = oldOsExit }()

	var got int
	myExit := func(code int) {
		got = code
	}

	osExit = myExit
	fatal("Fatal message")
	if exp := 1; got != exp {
		t.Errorf("Expected exit code: %d, got: %d", exp, got)
	}
}

func TestGetPageToken(t *testing.T) {
	u, _ := url.Parse("https://example.com/foobar")
	rec := httptest.NewRecorder()
	tokenpagination.PaginationHeader(rec, u, 100, 3, 10)
	assert.Equal(t, `eyJwYWdlIjoiNDAiLCJ2IjoxfQ`, getPageToken(rec.Result()), rec.Result().Header.Get("Link"))
}
