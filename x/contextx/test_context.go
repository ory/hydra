package contextx

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ory/hydra/driver/config"
)

// TestContextualizer is a mock implementation of the Contextualizer interface.
type TestContextualizer struct{}

// fakeNIDContext is a test key for NID.
const fakeNIDContext = "fake nid context"

// SetNIDContext sets the nid context for the given context.
func SetNIDContext(ctx context.Context, nid uuid.UUID) context.Context {
	return context.WithValue(ctx, fakeNIDContext, nid)
}

func (d *TestContextualizer) Network(ctx context.Context, network uuid.UUID) uuid.UUID {
	nid, ok := ctx.Value(fakeNIDContext).(uuid.UUID)
	if !ok {
		return network
	}
	return nid
}

func (d *TestContextualizer) Config(ctx context.Context, config *config.Provider) *config.Provider {
	return config
}
