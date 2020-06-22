package configuration

import (
	"os"
	"os/user"
	"strconv"
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
