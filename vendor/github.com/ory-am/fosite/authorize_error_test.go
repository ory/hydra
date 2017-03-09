package fosite_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/ory-am/fosite"
	. "github.com/ory-am/fosite/internal"
	"github.com/stretchr/testify/assert"
)

// Test for
// * https://tools.ietf.org/html/rfc6749#section-4.1.2.1
//   If the request fails due to a missing, invalid, or mismatching
//   redirection URI, or if the client identifier is missing or invalid,
//   the authorization server SHOULD inform the resource owner of the
//   error and MUST NOT automatically redirect the user-agent to the
//   invalid redirection URI.
// * https://tools.ietf.org/html/rfc6749#section-3.1.2
//   The redirection endpoint URI MUST be an absolute URI as defined by
//   [RFC3986] Section 4.3.  The endpoint URI MAY include an
//   "application/x-www-form-urlencoded" formatted (per Appendix B) query
//   component ([RFC3986] Section 3.4), which MUST be retained when adding
//   additional query parameters.  The endpoint URI MUST NOT include a
//   fragment component.
func TestWriteAuthorizeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	rw := NewMockResponseWriter(ctrl)
	req := NewMockAuthorizeRequester(ctrl)
	defer ctrl.Finish()

	var urls = []string{
		"https://foobar.com/",
		"https://foobar.com/?foo=bar",
	}
	var purls = []*url.URL{}
	for _, u := range urls {
		purl, _ := url.Parse(u)
		purls = append(purls, purl)
	}

	oauth2 := &Fosite{}
	header := http.Header{}
	for k, c := range []struct {
		err         error
		mock        func()
		checkHeader func(int)
	}{
		{
			err: ErrInvalidGrant,
			mock: func() {
				req.EXPECT().IsRedirectURIValid().Return(false)
				rw.EXPECT().Header().Return(header)
				rw.EXPECT().WriteHeader(http.StatusBadRequest)
				rw.EXPECT().Write(gomock.Any())
			},
			checkHeader: func(k int) {
				assert.Equal(t, "application/json", header.Get("Content-Type"), "%d", k)
			},
		},
		{
			err: ErrInvalidRequest,
			mock: func() {
				req.EXPECT().IsRedirectURIValid().Return(true)
				req.EXPECT().GetRedirectURI().Return(copyUrl(purls[0]))
				req.EXPECT().GetState().Return("foostate")
				req.EXPECT().GetResponseTypes().MaxTimes(2).Return(Arguments([]string{"code"}))
				rw.EXPECT().Header().Return(header)
				rw.EXPECT().WriteHeader(http.StatusFound)
			},
			checkHeader: func(k int) {
				a, _ := url.Parse("https://foobar.com/?error=invalid_request&error_description=The+request+is+missing+a+required+parameter%2C+includes+an+invalid+parameter+value%2C+includes+a+parameter+more+than+once%2C+or+is+otherwise+malformed&state=foostate")
				b, _ := url.Parse(header.Get("Location"))
				assert.Equal(t, a, b, "%d", k)
			},
		},
		{
			err: ErrInvalidRequest,
			mock: func() {
				req.EXPECT().IsRedirectURIValid().Return(true)
				req.EXPECT().GetRedirectURI().Return(copyUrl(purls[1]))
				req.EXPECT().GetState().Return("foostate")
				req.EXPECT().GetResponseTypes().MaxTimes(2).Return(Arguments([]string{"code"}))
				rw.EXPECT().Header().Return(header)
				rw.EXPECT().WriteHeader(http.StatusFound)
			},
			checkHeader: func(k int) {
				a, _ := url.Parse("https://foobar.com/?error=invalid_request&error_description=The+request+is+missing+a+required+parameter%2C+includes+an+invalid+parameter+value%2C+includes+a+parameter+more+than+once%2C+or+is+otherwise+malformed&foo=bar&state=foostate")
				b, _ := url.Parse(header.Get("Location"))
				assert.Equal(t, a, b, "%d", k)
			},
		},
		{
			err: ErrInvalidRequest,
			mock: func() {
				req.EXPECT().IsRedirectURIValid().Return(true)
				req.EXPECT().GetRedirectURI().Return(copyUrl(purls[0]))
				req.EXPECT().GetState().Return("foostate")
				req.EXPECT().GetResponseTypes().MaxTimes(2).Return(Arguments([]string{"token"}))
				rw.EXPECT().Header().Return(header)
				rw.EXPECT().WriteHeader(http.StatusFound)
			},
			checkHeader: func(k int) {
				a, _ := url.Parse("https://foobar.com/")
				a.Fragment = "error=invalid_request&error_description=The+request+is+missing+a+required+parameter%2C+includes+an+invalid+parameter+value%2C+includes+a+parameter+more+than+once%2C+or+is+otherwise+malformed&state=foostate"
				b, _ := url.Parse(header.Get("Location"))
				assert.Equal(t, a, b, "%d", k)
			},
		},
		{
			err: ErrInvalidRequest,
			mock: func() {
				req.EXPECT().IsRedirectURIValid().Return(true)
				req.EXPECT().GetRedirectURI().Return(copyUrl(purls[1]))
				req.EXPECT().GetState().Return("foostate")
				req.EXPECT().GetResponseTypes().MaxTimes(2).Return(Arguments([]string{"token"}))
				rw.EXPECT().Header().Return(header)
				rw.EXPECT().WriteHeader(http.StatusFound)
			},
			checkHeader: func(k int) {
				a, _ := url.Parse("https://foobar.com/?foo=bar")
				a.Fragment = "error=invalid_request&error_description=The+request+is+missing+a+required+parameter%2C+includes+an+invalid+parameter+value%2C+includes+a+parameter+more+than+once%2C+or+is+otherwise+malformed&state=foostate"
				b, _ := url.Parse(header.Get("Location"))
				assert.Equal(t, a.String(), b.String(), "%d", k)
			},
		},
		{
			err: ErrInvalidRequest,
			mock: func() {
				req.EXPECT().IsRedirectURIValid().Return(true)
				req.EXPECT().GetRedirectURI().Return(copyUrl(purls[0]))
				req.EXPECT().GetState().Return("foostate")
				req.EXPECT().GetResponseTypes().MaxTimes(2).Return(Arguments([]string{"code", "token"}))
				rw.EXPECT().Header().Return(header)
				rw.EXPECT().WriteHeader(http.StatusFound)
			},
			checkHeader: func(k int) {
				a, _ := url.Parse("https://foobar.com/")
				a.Fragment = "error=invalid_request&error_description=The+request+is+missing+a+required+parameter%2C+includes+an+invalid+parameter+value%2C+includes+a+parameter+more+than+once%2C+or+is+otherwise+malformed&state=foostate"
				b, _ := url.Parse(header.Get("Location"))
				assert.Equal(t, a, b, "%d", k)
			},
		},
		{
			err: ErrInvalidRequest,
			mock: func() {
				req.EXPECT().IsRedirectURIValid().Return(true)
				req.EXPECT().GetRedirectURI().Return(copyUrl(purls[1]))
				req.EXPECT().GetState().Return("foostate")
				req.EXPECT().GetResponseTypes().MaxTimes(2).Return(Arguments([]string{"code", "token"}))
				rw.EXPECT().Header().Return(header)
				rw.EXPECT().WriteHeader(http.StatusFound)
			},
			checkHeader: func(k int) {
				a, _ := url.Parse("https://foobar.com/?foo=bar")
				a.Fragment = "error=invalid_request&error_description=The+request+is+missing+a+required+parameter%2C+includes+an+invalid+parameter+value%2C+includes+a+parameter+more+than+once%2C+or+is+otherwise+malformed&state=foostate"
				b, _ := url.Parse(header.Get("Location"))
				assert.Equal(t, a.String(), b.String(), "%d", k)
			},
		},
	} {
		c.mock()
		oauth2.WriteAuthorizeError(rw, req, c.err)
		c.checkHeader(k)
		header = http.Header{}
		t.Logf("Passed test case %d", k)
	}
}

func copyUrl(u *url.URL) *url.URL {
	url, _ := url.Parse(u.String())
	return url
}
