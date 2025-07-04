// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/sjson"

	"github.com/ory/jsonschema/v3"

	"github.com/spf13/cast"
	"github.com/tidwall/gjson"

	"github.com/ory/x/castx"
	"github.com/ory/x/jsonschemax"
)

var isNumRegex = regexp.MustCompile("^[0-9]+$")

func NewKoanfEnv(prefix string, rawSchema []byte, schema *jsonschema.Schema) (*Env, error) {
	paths, err := getSchemaPaths(rawSchema, schema)
	if err != nil {
		return nil, err
	}

	return &Env{
		paths:  paths,
		prefix: prefix,
	}, nil
}

// Env implements an environment variables provider.
type Env struct {
	prefix string
	paths  []jsonschemax.Path
}

// ReadBytes is not supported by the env provider.
func (e *Env) ReadBytes() ([]byte, error) {
	return nil, errors.New("env provider does not support this method")
}

// Read reads all available environment variables into a key:value map
// and returns it.
func (e *Env) Read() (map[string]interface{}, error) {
	// Collect the environment variable keys.
	var keys []string
	for _, k := range os.Environ() {
		if e.prefix != "" {
			if strings.HasPrefix(k, e.prefix) {
				keys = append(keys, k)
			}
		} else {
			keys = append(keys, k)
		}
	}

	raw := "{}"
	var err error
	for _, k := range keys {
		parts := strings.SplitN(k, "=", 2)

		key, value := e.extract(parts[0], parts[1])
		// If the callback blanked the key, it should be omitted
		if key == "" {
			continue
		}

		raw, err = sjson.Set(raw, key, value)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil, errors.WithStack(err)
	}

	return m, nil
}

// Watch is not supported.
func (e *Env) Watch(cb func(event interface{}, err error)) error {
	return errors.New("env provider does not support this method")
}

func (e *Env) extract(key string, value string) (string, interface{}) {
	key = strings.Replace(strings.ToLower(strings.TrimPrefix(key, e.prefix)), "_", ".", -1)

	for _, path := range e.paths {
		normalized := strings.Replace(path.Name, "_", ".", -1)
		name := path.Name

		// Crazy hack to get arrays working.
		var indices []string
		searchParts := strings.Split(normalized, ".")
		keyParts := strings.Split(key, ".")
		if len(searchParts) == len(keyParts) {
			for k, search := range searchParts {
				if search != keyParts[k] {
					indices = nil
				}

				if search != "#" {
					continue
				}

				if !isNumRegex.MatchString(keyParts[k]) {
					continue
				}

				searchParts[k] = keyParts[k]
				indices = append(indices, keyParts[k])
			}
		}

		if len(indices) > 0 {
			normalized = strings.Join(searchParts, ".")
			for _, index := range indices {
				name = strings.Replace(name, "#", index, 1)
			}
		}

		if normalized == key {
			switch path.TypeHint {
			case jsonschemax.String:
				return name, cast.ToString(value)
			case jsonschemax.Float:
				return name, cast.ToFloat64(value)
			case jsonschemax.Int:
				return name, cast.ToInt64(value)
			case jsonschemax.Bool:
				return name, cast.ToBool(value)
			case jsonschemax.Nil:
				return name, nil
			case jsonschemax.BoolSlice:
				if !gjson.Valid(value) {
					return name, cast.ToBoolSlice(value)
				}
				fallthrough
			case jsonschemax.StringSlice:
				if !gjson.Valid(value) {
					return name, castx.ToStringSlice(value)
				}
				fallthrough
			case jsonschemax.IntSlice:
				if !gjson.Valid(value) {
					return name, cast.ToIntSlice(value)
				}
				fallthrough
			case jsonschemax.FloatSlice:
				if !gjson.Valid(value) {
					return name, castx.ToFloatSlice(value)
				}
				fallthrough
			case jsonschemax.JSON:
				return name, decode(value)
			default:
				return name, value
			}
		}
	}

	return "", nil
}

func decode(value string) (v interface{}) {
	b := []byte(value)
	var arr []interface{}
	if err := json.Unmarshal(b, &arr); err == nil {
		return &arr
	}
	h := map[string]interface{}{}
	if err := json.Unmarshal(b, &h); err == nil {
		return &h
	}
	return nil
}
