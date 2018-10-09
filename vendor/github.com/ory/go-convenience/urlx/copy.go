package urlx

import "net/url"

// Copy creates returns a copy of the provided url.URL pointer.
func Copy(u *url.URL) *url.URL {
	a := new(url.URL)
	*a = *u
	return a
}
