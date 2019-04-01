// +build noplugin

package driver

import (
	"github.com/jmoiron/sqlx"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/x"
	"github.com/ory/x/dbal"
	"github.com/ory/x/sqlcon"
	"github.com/ory/x/urlx"
)

type RegistryNoPlugin struct {
	*RegistryBase
	db          *sqlx.DB
	dbalOptions []sqlcon.OptionModifier
}

var _ Registry = new(RegistryNoPlugin)

func init() {
	dbal.RegisterDriver(NewRegistryNoPlugin())
}

func NewRegistryNoPlugin() *RegistryNoPlugin {
	r := &RegistryNoPlugin{
		RegistryBase: new(RegistryBase),
	}
	r.RegistryBase.with(r)
	return r
}

func (m *RegistryNoPlugin) Init() error {
	panic("Unable to load plugin connection because 'noplugin' tag was declared")
}

func (m *RegistryNoPlugin) CanHandle(dsn string) bool {
	u := urlx.ParseOrFatal(m.Logger(), dsn)
	return u.Scheme == "plugin"
}

func (m *RegistryNoPlugin) Ping() error {
	panic("Unable to load plugin connection because 'noplugin' tag was declared")
}

func (m *RegistryNoPlugin) ClientManager() client.Manager {
	panic("Unable to load plugin connection because 'noplugin' tag was declared")
}

func (m *RegistryNoPlugin) ConsentManager() consent.Manager {
	panic("Unable to load plugin connection because 'noplugin' tag was declared")
}

func (m *RegistryNoPlugin) OAuth2Storage() x.FositeStorer {
	panic("Unable to load plugin connection because 'noplugin' tag was declared")
}

func (m *RegistryNoPlugin) KeyManager() jwk.Manager {
	panic("Unable to load plugin connection because 'noplugin' tag was declared")
}
