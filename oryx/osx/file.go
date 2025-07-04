// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package osx

import (
	"encoding/base64"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/pkg/errors"

	"github.com/ory/x/httpx"
)

type options struct {
	disableFileLoader            bool
	disableHTTPLoader            bool
	disableBase64Loader          bool
	base64enc                    *base64.Encoding
	disableResilientBase64Loader bool
	hc                           *retryablehttp.Client
}

type Option func(o *options)

func (o *options) apply(opts []Option) *options {
	for _, f := range opts {
		f(o)
	}
	return o
}

func newOptions() *options {
	return &options{
		disableFileLoader:   false,
		disableHTTPLoader:   false,
		disableBase64Loader: false,
		base64enc:           base64.RawURLEncoding,
		hc:                  httpx.NewResilientClient(),
	}
}

// WithDisabledFileLoader disables the file loader.
func WithDisabledFileLoader() Option {
	return func(o *options) {
		o.disableFileLoader = true
	}
}

// WithEnabledFileLoader enables the file loader.
func WithEnabledFileLoader() Option {
	return func(o *options) {
		o.disableFileLoader = false
	}
}

// WithDisabledHTTPLoader disables the HTTP loader.
func WithDisabledHTTPLoader() Option {
	return func(o *options) {
		o.disableHTTPLoader = true
	}
}

// WithEnabledHTTPLoader enables the HTTP loader.
func WithEnabledHTTPLoader() Option {
	return func(o *options) {
		o.disableHTTPLoader = false
	}
}

// WithDisabledBase64Loader disables the base64 loader.
func WithDisabledBase64Loader() Option {
	return func(o *options) {
		o.disableBase64Loader = true
	}
}

// WithEnabledBase64Loader disables the base64 loader.
func WithEnabledBase64Loader() Option {
	return func(o *options) {
		o.disableBase64Loader = false
	}
}

// WithBase64Encoding sets the base64 encoding.
func WithBase64Encoding(enc *base64.Encoding) Option {
	return func(o *options) {
		o.base64enc = enc
	}
}

// WithoutResilientBase64Encoding sets the base64 encoding.
func WithoutResilientBase64Encoding() Option {
	return func(o *options) {
		o.disableResilientBase64Loader = true
	}
}

// WithHTTPClient sets the HTTP client.
func WithHTTPClient(hc *retryablehttp.Client) Option {
	return func(o *options) {
		o.hc = hc
	}
}

// RestrictedReadFile works similar to ReadFileFromAllSources but has all
// sources disabled per default. You need to enable the loaders you wish to use
// explicitly.
func RestrictedReadFile(source string, opts ...Option) (bytes []byte, err error) {
	o := newOptions()
	o.disableFileLoader = true
	o.disableBase64Loader = true
	o.disableHTTPLoader = true
	return readFile(source, o.apply(opts))
}

// ReadFileFromAllSources reads a file from base64, http, https, and file sources.
//
// Using options, you can disable individual loaders. For example, the following will
// return an error:
//
//	ReadFileFromAllSources("https://foo.bar/baz.txt", WithDisabledHTTPLoader())
//
// Possible formats are:
//
// - /path/to/file
// - file:///path/to/file
// - https://host.com/path/to/file
// - http://host.com/path/to/file
// - base64://<base64 encoded string>
//
// For more options, check:
//
// - WithDisabledFileLoader
// - WithDisabledHTTPLoader
// - WithDisabledBase64Loader
// - WithBase64Encoding
// - WithHTTPClient
func ReadFileFromAllSources(source string, opts ...Option) (bytes []byte, err error) {
	return readFile(source, newOptions().apply(opts))
}

func readFile(source string, o *options) (bytes []byte, err error) {
	parsed, err := url.Parse(source)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse URL")
	}

	switch parsed.Scheme {
	case "":
		if o.disableFileLoader {
			return nil, errors.New("file loader disabled")
		}

		//#nosec G304 -- false positive
		bytes, err = os.ReadFile(source)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read the file")
		}
	case "file":
		if o.disableFileLoader {
			return nil, errors.New("file loader disabled")
		}

		//#nosec G304 -- false positive
		bytes, err = os.ReadFile(parsed.Host + parsed.Path)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read the file")
		}
	case "http", "https":
		if o.disableHTTPLoader {
			return nil, errors.New("http(s) loader disabled")
		}
		resp, err := o.hc.Get(parsed.String())
		if err != nil {
			return nil, errors.Wrap(err, "unable to load remote file")
		}
		defer resp.Body.Close()

		bytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read the HTTP response body")
		}
	case "base64":
		if o.disableBase64Loader {
			return nil, errors.New("base64 loader disabled")
		}

		if o.disableResilientBase64Loader {
			bytes, err = o.base64enc.DecodeString(strings.TrimPrefix(source, "base64://"))
			if err != nil {
				return nil, errors.Wrap(err, "unable to base64 decode the location")
			}
			return bytes, nil
		}

		for _, enc := range []*base64.Encoding{
			base64.StdEncoding,
			base64.URLEncoding,
			base64.RawURLEncoding,
			base64.RawStdEncoding,
		} {
			bytes, err = enc.DecodeString(strings.TrimPrefix(source, "base64://"))
			if err == nil {
				return bytes, nil
			}
		}

		return nil, errors.Wrap(err, "unable to base64 decode the location")
	default:
		return nil, errors.Errorf("unsupported source `%s`", parsed.Scheme)
	}

	return bytes, nil

}
