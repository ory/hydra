package pkg

import (
	"net/url"
	"path"
)

func JoinURL(u *url.URL, args ...string) (ep *url.URL) {
	ep = new(url.URL)
	*ep = *u
	ep.Path = path.Join(append([]string{ep.Path}, args...)...)
	return ep
}
