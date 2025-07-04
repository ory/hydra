// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package urlx

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/ory/x/logrusx"
)

// winPathRegex is a regex for [DRIVE-LETTER]:
var winPathRegex = regexp.MustCompile("^[A-Za-z]:.*")

// Parse parses rawURL into a URL structure with special handling for file:// URLs
//
// File URLs with relative paths (file://../file, ../file) will be returned as a
// url.URL object without the Scheme set to "file". This is because the file
// scheme does not support relative paths. Make sure to check for
// both "file" or "" (an empty string) in URL.Scheme if you are looking for
// a file path.
//
// Use the companion function GetURLFilePath() to get a file path suitable
// for the current operating system.
func Parse(rawURL string) (*url.URL, error) {
	lcRawURL := strings.ToLower(rawURL)
	if strings.HasPrefix(lcRawURL, "file:///") {
		return url.Parse(rawURL)
	}

	// Normally the first part after file:// is a hostname, but since
	// this is often misused we interpret the URL like a normal path
	// by removing the "file://" from the beginning (if it exists)
	rawURL = trimPrefixIC(rawURL, "file://")

	if winPathRegex.MatchString(rawURL) {
		// Windows path
		return url.Parse("file:///" + rawURL)
	}

	if strings.HasPrefix(lcRawURL, "\\\\") {
		// Windows UNC path
		// We extract the hostname and create an appropriate file:// URL
		// based on the hostname and the path
		host, path := extractUNCPathParts(rawURL)
		// It is safe to replace the \ with / here because this is POSIX style path
		return url.Parse("file://" + host + strings.ReplaceAll(path, "\\", "/"))
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	// Since go1.19:
	//
	// > The URL type now distinguishes between URLs with no authority and URLs with an empty authority.
	// > For example, http:///path has an empty authority (host), while http:/path has none.
	//
	// See https://golang.org/doc/go1.19#net/url for more details.
	parsed.OmitHost = false
	return parsed, nil
}

// ParseOrPanic parses a url or panics.
func ParseOrPanic(in string) *url.URL {
	out, err := url.Parse(in)
	if err != nil {
		panic(err.Error())
	}
	return out
}

// ParseOrFatal parses a url or fatals.
func ParseOrFatal(l *logrusx.Logger, in string) *url.URL {
	out, err := url.Parse(in)
	if err != nil {
		l.WithError(err).Fatalf("Unable to parse url: %s", in)
	}
	return out
}

// ParseRequestURIOrPanic parses a request uri or panics.
func ParseRequestURIOrPanic(in string) *url.URL {
	out, err := url.ParseRequestURI(in)
	if err != nil {
		panic(err.Error())
	}
	return out
}

// ParseRequestURIOrFatal parses a request uri or fatals.
func ParseRequestURIOrFatal(l *logrusx.Logger, in string) *url.URL {
	out, err := url.ParseRequestURI(in)
	if err != nil {
		l.WithError(err).Fatalf("Unable to parse url: %s", in)
	}
	return out
}

func extractUNCPathParts(uncPath string) (host, path string) {
	parts := strings.Split(strings.TrimPrefix(uncPath, "\\\\"), "\\")
	host = parts[0]
	if len(parts) > 0 {
		path = "\\" + strings.Join(parts[1:], "\\")
	}
	return host, path
}

// trimPrefixIC returns s without the provided leading prefix string using
// case insensitive matching.
// If s doesn't start with prefix, s is returned unchanged.
func trimPrefixIC(s, prefix string) string {
	if strings.HasPrefix(strings.ToLower(s), prefix) {
		return s[len(prefix):]
	}
	return s
}
