// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	"strings"

	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/ory/jsonschema/v3"
	"github.com/ory/x/jsonschemax"
)

type PFlagProvider struct {
	p     *posflag.Posflag
	paths []jsonschemax.Path
}

func NewPFlagProvider(rawSchema []byte, schema *jsonschema.Schema, f *pflag.FlagSet, k *koanf.Koanf) (*PFlagProvider, error) {
	paths, err := getSchemaPaths(rawSchema, schema)
	if err != nil {
		return nil, err
	}
	return &PFlagProvider{
		p:     posflag.Provider(f, ".", k),
		paths: paths,
	}, nil
}

func (p *PFlagProvider) ReadBytes() ([]byte, error) {
	return nil, errors.New("pflag provider does not support this method")
}

func (p *PFlagProvider) Read() (map[string]interface{}, error) {
	all, err := p.p.Read()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	knownFlags := make(map[string]interface{}, len(all))
	for k, v := range all {
		k = strings.ReplaceAll(k, ".", "-")
		for _, path := range p.paths {
			normalized := strings.ReplaceAll(path.Name, ".", "-")
			if k == normalized {
				knownFlags[k] = v
				break
			}
		}
	}
	return knownFlags, nil
}

var _ koanf.Provider = (*PFlagProvider)(nil)
