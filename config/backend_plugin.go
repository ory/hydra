package config

import (
	"plugin"
	"github.com/ory/ladon"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/warden/group"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
	"github.com/Sirupsen/logrus"
	"github.com/ory/fosite"
)

type PluginConnection struct {
	Config     *Config
	plugin     *plugin.Plugin
	didConnect bool
	Logger     logrus.FieldLogger
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
	} else if c, ok := l.(func(url string) error); !ok {
		return errors.New("Unable to type assert `Connect`")
	} else {
		if err := c(cf.DatabaseURL); err != nil {
			return errors.Wrap(err, "Could not Connect to database")
		}
	}
	return nil
}

func (c *PluginConnection) NewClientManager() (client.Manager, error) {
	if err := c.load(); err != nil {
		return nil, errors.WithStack(err)
	}

	ctx := c.Config.Context()
	if l, err := c.plugin.Lookup("NewClientManager"); err != nil {
		return nil, errors.Wrap(err, "Unable to look up `NewClientManager`")
	} else if m, ok := l.(func(fosite.Hasher) client.Manager); !ok {
		return nil, errors.New("Unable to type assert `NewClientManager`")
	} else {
		return m(ctx.Hasher), nil
	}
}

func (c *PluginConnection) NewGroupManager() (group.Manager, error) {
	if err := c.load(); err != nil {
		return nil, errors.WithStack(err)
	}

	if l, err := c.plugin.Lookup("NewGroupManager"); err != nil {
		return nil, errors.Wrap(err, "Unable to look up `NewGroupManager`")
	} else if m, ok := l.(func() group.Manager); !ok {
		return nil, errors.New("Unable to type assert `NewGroupManager`")
	} else {
		return m(), nil
	}
}

func (c *PluginConnection) NewJWKManager() (jwk.Manager, error) {
	if err := c.load(); err != nil {
		return nil, errors.WithStack(err)
	}

	if l, err := c.plugin.Lookup("NewJWKManager"); err != nil {
		return nil, errors.Wrap(err, "Unable to look up `NewJWKManager`")
	} else if m, ok := l.(func(*jwk.AEAD) jwk.Manager); !ok {
		return nil, errors.New("Unable to type assert `NewJWKManager`")
	} else {
		return m(&jwk.AEAD{
			Key: c.Config.GetSystemSecret(),
		}), nil
	}
}

func (c *PluginConnection) NewOAuth2Manager(clientManager client.Manager) (pkg.FositeStorer, error) {
	if err := c.load(); err != nil {
		return nil, errors.WithStack(err)
	}

	if l, err := c.plugin.Lookup("NewOAuth2Manager"); err != nil {
		return nil, errors.Wrap(err, "Unable to look up `NewOAuth2Manager`")
	} else if m, ok := l.(func(client.Manager, logrus.FieldLogger) pkg.FositeStorer); !ok {
		return nil, errors.New("Unable to type assert `NewOAuth2Manager`")
	} else {
		return m(clientManager, c.Config.GetLogger()), nil
	}
}

func (c *PluginConnection) NewPolicyManager() (ladon.Manager, error) {
	if err := c.load(); err != nil {
		return nil, errors.WithStack(err)
	}

	if l, err := c.plugin.Lookup("NewPolicyManager"); err != nil {
		return nil, errors.Wrap(err, "Unable to look up `NewPolicyManager`")
	} else if m, ok := l.(func() ladon.Manager); !ok {
		return nil, errors.Errorf("Unable to type assert `NewPolicyManager`, got %v", l)
	} else {
		return m(), nil
	}
}
