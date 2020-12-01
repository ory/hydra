package config

import (
	"net/url"
	"os"
	"os/user"
	"strconv"

	"github.com/ory/x/logrusx"
)

type UnixPermission struct {
	Owner string
	Group string
	Mode  os.FileMode
}

func (p *UnixPermission) SetPermission(file string) error {
	var e error
	e = os.Chmod(file, p.Mode)
	if e != nil {
		return e
	}

	gid := -1
	uid := -1

	if p.Owner != "" {
		var user_obj *user.User
		user_obj, e = user.Lookup(p.Owner)
		if e != nil {
			return e
		}
		uid, e = strconv.Atoi(user_obj.Uid)
		if e != nil {
			return e
		}
	}
	if p.Group != "" {
		var group *user.Group
		group, e := user.LookupGroup(p.Group)
		if e != nil {
			return e
		}
		gid, e = strconv.Atoi(group.Gid)
		if e != nil {
			return e
		}
	}

	e = os.Chown(file, uid, gid)
	if e != nil {
		return e
	}
	return nil
}

func MustValidate(l *logrusx.Logger, p *ViperProvider) {
	if p.ServesHTTPS() {
		if p.IssuerURL().String() == "" {
			l.Fatalf(`Configuration key "%s" must be set unless flag "--dangerous-force-http" is set. To find out more, use "hydra help serve".`, ViperKeyIssuerURL)
		}

		if p.IssuerURL().Scheme != "https" {
			l.Fatalf(`Scheme from configuration key "%s" must be "https" unless --dangerous-force-http is passed but got scheme in value "%s" is "%s". To find out more, use "hydra help serve".`, ViperKeyIssuerURL, p.IssuerURL().String(), p.IssuerURL().Scheme)
		}

		if len(p.InsecureRedirects()) > 0 {
			l.Fatal(`Flag --dangerous-allow-insecure-redirect-urls can only be used in combination with flag --dangerous-force-http`)
		}
	}
}

func urlRoot(u *url.URL) *url.URL {
	if u.Path == "" {
		u.Path = "/"
	}
	return u
}
