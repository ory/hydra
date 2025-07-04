// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build windows
// +build windows

package urlx

import (
	"net/url"
	"path/filepath"
	"strings"
)

// GetURLFilePath returns the path of a URL that is compatible with the runtime os filesystem
func GetURLFilePath(u *url.URL) string {
	if u == nil {
		return ""
	}
	if !(u.Scheme == "file" || u.Scheme == "") {
		return u.Path
	}

	fPath := u.Path
	if u.Host != "" {
		// Make UNC Path
		fPath = "\\\\" + u.Host + filepath.FromSlash(fPath)
		return fPath
	}
	fPathTrimmed := strings.TrimLeft(fPath, "/")
	if winPathRegex.MatchString(fPathTrimmed) {
		// On Windows we should remove the initial path separator in case this
		// is a normal path (for example: "\c:\" -> "c:\"")
		fPath = fPathTrimmed
	}
	return filepath.FromSlash(fPath)
}
