// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

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
		var userObj *user.User
		userObj, e = user.Lookup(p.Owner)
		if e != nil {
			return e
		}
		uid, e = strconv.Atoi(userObj.Uid)
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
