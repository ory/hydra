// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonx

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/ory/x/osx"
)

type options struct {
	ignoreKeys  []string
	onlySchemes []string
}

type OptionsModifier func(*options)

func newOptions(o []OptionsModifier) *options {
	opt := &options{}
	for _, f := range o {
		f(opt)
	}
	return opt
}

func WithIgnoreKeys(keys ...string) OptionsModifier {
	return func(o *options) {
		o.ignoreKeys = keys
	}
}

func WithOnlySchemes(scheme ...string) OptionsModifier {
	return func(o *options) {
		o.onlySchemes = scheme
	}
}

func EmbedSources(in json.RawMessage, opts ...OptionsModifier) (out json.RawMessage, err error) {
	out = make([]byte, len(in))
	copy(out, in)
	if err := embed(gjson.ParseBytes(in), nil, &out, newOptions(opts)); err != nil {
		return nil, err
	}
	return out, nil
}

func embed(parsed gjson.Result, parents []string, result *json.RawMessage, o *options) (err error) {
	if parsed.IsObject() {
		parsed.ForEach(func(k, v gjson.Result) bool {
			err = embed(v, append(parents, strings.ReplaceAll(k.String(), ".", "\\.")), result, o)
			return err == nil
		})
		if err != nil {
			return err
		}
	} else if parsed.IsArray() {
		for kk, vv := range parsed.Array() {
			if err = embed(vv, append(parents, strconv.Itoa(kk)), result, o); err != nil {
				return err
			}
		}
	} else if parsed.Type != gjson.String {
		return nil
	}

	if len(parents) > 0 && slices.Contains(o.ignoreKeys, parents[len(parents)-1]) {
		return nil
	}

	loc, err := url.ParseRequestURI(parsed.String())
	if err != nil {
		// Not a URL, return
		return nil
	}

	if len(o.onlySchemes) == 0 {
		if loc.Scheme != "file" && loc.Scheme != "http" && loc.Scheme != "https" && loc.Scheme != "base64" {
			// Not a known pattern, ignore
			return nil
		}
	} else if !slices.Contains(o.onlySchemes, loc.Scheme) {
		// Not a known pattern, ignore
		return nil
	}

	contents, err := osx.ReadFileFromAllSources(loc.String())
	if err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(contents)
	key := strings.Join(parents, ".")
	if key == "" {
		key = "@"
	}

	interim, err := sjson.SetBytes(*result, key, "base64://"+encoded)
	if err != nil {
		return err
	}

	*result = interim
	return
}
