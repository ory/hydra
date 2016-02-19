package provider

import (
	"github.com/go-errors/errors"
	"strings"
)

type Registry interface {
	Find(id string) (Provider, error)
}

type defaultRegistry struct {
	providers map[string]Provider
}

func NewRegistry(providers []Provider) Registry {
	r := &defaultRegistry{map[string]Provider{}}
	for _, v := range providers {
		r.add(v)
	}

	return r
}

func (r *defaultRegistry) add(provider Provider) {
	id := strings.ToLower(provider.GetID())
	r.providers[id] = provider
}

func (r *defaultRegistry) Find(id string) (Provider, error) {
	id = strings.ToLower(id)
	p, ok := r.providers[id]
	if !ok {
		return nil, errors.Errorf("Provider %s not found", id)
	}
	return p, nil
}
