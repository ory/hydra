package cliclient

import (
	"github.com/pkg/errors"

	"github.com/ory/x/cmdx"

	"github.com/spf13/cobra"

	hydra "github.com/ory/hydra-client-go"
)

type ContextKey int

const (
	ClientContextKey ContextKey = iota + 1
)

func NewClient(cmd *cobra.Command) (*hydra.APIClient, error) {
	if f, ok := cmd.Context().Value(ClientContextKey).(func(cmd *cobra.Command) (*hydra.APIClient, error)); ok {
		return f(cmd)
	} else if f != nil {
		return nil, errors.Errorf("ClientContextKey was expected to be *client.OryHydra but it contained an invalid type %T ", f)
	}

	hc, target, err := cmdx.NewClient(cmd)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	conf := hydra.NewConfiguration()
	conf.HTTPClient = hc
	conf.Servers = hydra.ServerConfigurations{{URL: target.String()}}
	return hydra.NewAPIClient(conf), nil
}
