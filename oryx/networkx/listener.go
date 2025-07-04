// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package networkx

import (
	"net"
	"strings"

	"github.com/ory/x/configx"
)

func AddressIsUnixSocket(address string) bool {
	return strings.HasPrefix(address, "unix:")
}

func MakeListener(address string, socketPermission *configx.UnixPermission) (net.Listener, error) {
	if AddressIsUnixSocket(address) {
		addr := strings.TrimPrefix(address, "unix:")
		l, err := net.Listen("unix", addr)
		if err != nil {
			return nil, err
		}
		err = socketPermission.SetPermission(addr)
		if err != nil {
			return nil, err
		}
		return l, nil
	}
	return net.Listen("tcp", address)
}
