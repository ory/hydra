package flow

import (
	"context"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/oauth2/flowctx"
)

// NewMiddleware returns a new flow-specific middleware.
func NewMiddleware(dependencies flowctx.Dependencies) *flowctx.Middleware {
	return flowctx.NewMiddleware(flowctx.FlowCookie, dependencies)
}

func FromCtx(ctx context.Context, p interface {
	GetConcreteClient(context.Context, string) (*client.Client, error)
}) (*Flow, error) {
	return flowctx.FromCtx[Flow](ctx, flowctx.FlowCookie, flowctx.WithPostDecodeHook(func(ctx context.Context, val any) error {
		f := val.(*Flow)
		c, err := p.GetConcreteClient(ctx, f.ClientID)
		if err != nil {
			return err
		}
		f.Client = c

		return nil
	}))
}
func SetInCtx(ctx context.Context, f *Flow) error {
	return flowctx.SetCtx(ctx, flowctx.FlowCookie, f)
}
