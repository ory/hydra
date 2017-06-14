package config

import (
	"plugin"
	"github.com/ory/ladon"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/warden/group"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
	"os"
)

type ClientFactory func(*Config) client.Manager
type GroupFactory func(*Config) group.Manager
type JWKFactory func(*Config) jwk.Manager
type OAuth2Factory func(client.Manager, *Config) pkg.FositeStorer
type PolicyFactory func(*Config) ladon.Manager
type SchemaCreator func() error
type Connector func(url string) error

type PluginConnection struct {
	Config     *Config
	plugin     *plugin.Plugin
	didConnect bool
}

func (c *PluginConnection) load() error {
	cf := c.Config
	if c.plugin != nil {
		return nil
	}

	p, err := plugin.Open(cf.DatabasePlugin)
	if err != nil {
		return errors.WithStack(err)
	}

	c.plugin = p
	return nil
}

func (c *PluginConnection) Connect() error {
	cf := c.Config
	if c.didConnect {
		return nil
	}

	if err := c.load(); err != nil {
		return errors.WithStack(err)
	}

	if l, err := c.plugin.Lookup("Connect"); err != nil {
		return errors.Wrap(err, "Unable to look up `Connect`")
	} else if c, ok := l.(Connector); !ok {
		return errors.Wrap(err, "Unable to type assert `Connect`")
	} else {
		if err := c(os.Getenv(cf.DatabaseURL)); err != nil {
			return errors.Wrap(err, "Could not Connect to database")
		}
	}
	return nil
}

func (c *PluginConnection) NewPolicyManager() (ladon.Manager, error) {
	if err := c.load(); err != nil {
		return nil, errors.WithStack(err)
	}

	if l, err := c.plugin.Lookup("NewPolicyManager"); err != nil {
		return nil, errors.Wrap(err, "Unable to look up `NewPolicyManager`")
	} else if m, ok := l.(PolicyFactory); !ok {
		return nil, errors.Wrap(err, "Unable to type assert `NewPolicyManager`")
	} else {
		return m(c.Config), nil
	}
}
