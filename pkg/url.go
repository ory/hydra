package pkg

import (
	"net/url"
	"path"

	"github.com/ory-am/common/pkg"
)

func CopyURL(u *url.URL) *url.URL {
	a := new(url.URL)
	*a = *u
	return a
}

func JoinURLStrings(host string, args ...string) string {
	return pkg.JoinURL(host, args...)
}

func JoinURL(u *url.URL, args ...string) (ep *url.URL) {
	ep = CopyURL(u)
	ep.Path = path.Join(append([]string{ep.Path}, args...)...)
	return ep
}
