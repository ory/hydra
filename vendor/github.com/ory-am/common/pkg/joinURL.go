package pkg

import (
	"fmt"
	"net/url"
	"path"
)

func JoinURL(host string, parts ...string) string {
	var trailing string

	last := parts[len(parts)-1]
	if last[len(last)-1:] == "/" {
		trailing = "/"
	}

	u, err := url.Parse(host)
	if err != nil {
		return fmt.Sprintf("%s%s%s", path.Join(append([]string{u.Path}, parts...)...), trailing)
	}

	if u.Path == "" {
		u.Path = "/"
	}
	return fmt.Sprintf("%s://%s%s%s", u.Scheme, u.Host, path.Join(append([]string{u.Path}, parts...)...), trailing)
}
