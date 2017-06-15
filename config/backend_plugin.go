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
	"github.com/Sirupsen/logrus"
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
	Logger logrus.FieldLogger
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

func (c *PluginConnection) NewClientManager() (client.Manager, error) {
	if err := c.load(); err != nil {
		return nil, errors.WithStack(err)
	}

	if l, err := c.plugin.Lookup("NewClientManager"); err != nil {
		return nil, errors.Wrap(err, "Unable to look up `NewClientManager`")
	} else if m, ok := l.(ClientFactory); !ok {
		return nil, errors.Wrap(err, "Unable to type assert `NewClientManager`")
	} else {
		return m(c.Config), nil
	}
}

func (c *PluginConnection) NewGroupManager() (group.Manager, error) {
	if err := c.load(); err != nil {
		return nil, errors.WithStack(err)
	}

	if l, err := c.plugin.Lookup("NewGroupManager"); err != nil {
		return nil, errors.Wrap(err, "Unable to look up `NewGroupManager`")
	} else if m, ok := l.(GroupFactory); !ok {
		return nil, errors.Wrap(err, "Unable to type assert `NewGroupManager`")
	} else {
		return m(c.Config), nil
	}
}

func (c *PluginConnection) NewJWKManager() (jwk.Manager, error) {
	if err := c.load(); err != nil {
		return nil, errors.WithStack(err)
	}

	if l, err := c.plugin.Lookup("NewJWKManager"); err != nil {
		return nil, errors.Wrap(err, "Unable to look up `NewJWKManager`")
	} else if m, ok := l.(JWKFactory); !ok {
		return nil, errors.Wrap(err, "Unable to type assert `NewJWKManager`")
	} else {
		return m(c.Config), nil
	}
}

func (c *PluginConnection) NewOAuth2Manager(clientManager client.Manager) (pkg.FositeStorer, error) {
	if err := c.load(); err != nil {
		return nil, errors.WithStack(err)
	}

	if l, err := c.plugin.Lookup("NewOAuth2Manager"); err != nil {
		return nil, errors.Wrap(err, "Unable to look up `NewOAuth2Manager`")
	} else if m, ok := l.(OAuth2Factory); !ok {
		return nil, errors.Wrap(err, "Unable to type assert `NewOAuth2Manager`")
	} else {
		return m(clientManager, c.Config), nil
	}
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
