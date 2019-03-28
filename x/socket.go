package x

import "strings"

func AddressIsUnixSocket(address string) bool {
	return strings.HasPrefix(address, "unix:")
}
