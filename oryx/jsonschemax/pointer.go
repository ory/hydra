// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonschemax

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// JSONPointerToDotNotation converts JSON Pointer "#/foo/bar" to dot-notation "foo.bar".
func JSONPointerToDotNotation(pointer string) (string, error) {
	if !strings.HasPrefix(pointer, "#/") {
		return pointer, errors.Errorf("remote JSON pointers are not supported: %s", pointer)
	}

	var path []string
	for _, item := range strings.Split(strings.TrimPrefix(pointer, "#/"), "/") {
		item = strings.Replace(item, "~1", "/", -1)
		item = strings.Replace(item, "~0", "~", -1)
		item, err := url.PathUnescape(item)
		if err != nil {
			return "", err
		}
		path = append(path, strings.ReplaceAll(item, ".", "\\."))
	}

	return strings.Join(path, "."), nil
}
