package contextx

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ory/hydra/driver/config"
)

type (
	Contextualizer interface {
		Network(ctx context.Context, network uuid.UUID) uuid.UUID
		Config(ctx context.Context, config *config.Provider) *config.Provider
	}
	ContextualizerProvider interface {
		Contextualizer() Contextualizer
	}
	DefaultContextualizer struct{}
	StaticContextualizer  struct {
		NID uuid.UUID
		C   *config.Provider
	}
)

var _ Contextualizer = (*DefaultContextualizer)(nil)

func (d *DefaultContextualizer) Network(ctx context.Context, network uuid.UUID) uuid.UUID {
	if network == uuid.Nil {
		panic("NetworkID called before initialized")
	}
	return network
}

func (d *DefaultContextualizer) Config(ctx context.Context, config *config.Provider) *config.Provider {
	return config
}

func (d *StaticContextualizer) Network(ctx context.Context, network uuid.UUID) uuid.UUID {
	return d.NID
}

func (d *StaticContextualizer) Config(ctx context.Context, config *config.Provider) *config.Provider {
	return d.C
}
