package path

import (
	"path"
	"github.com/Sirupsen/logrus"
	"net/url"
)

func Join(u string, paths ...string) string {
	ur, err := url.Parse(u)
	if err != nil {
		logrus.WithError(err).WithField("url", u).Panic("Could not parse url")
	}
	ur.Path = path.Join(append([]string{ur.Path}, paths...)...)
	return ur.String()
}
