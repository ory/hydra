//
// Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
// @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
// @license 	Apache-2.0
//

// +build !noplugin

package config

import (
	"plugin"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type PluginConnection struct {
	Config     *Config
	plugin     *plugin.Plugin
	didConnect bool
	connector  BackendConnector
	Logger     logrus.FieldLogger
}

func (c *PluginConnection) load() error {
	if c.plugin != nil {
		return nil
	}

	cf := c.Config
	p, err := plugin.Open(cf.DatabasePlugin)
	if err != nil {
		return errors.WithStack(err)
	}

	c.plugin = p
	return nil
}

func (c *PluginConnection) Load() error {
	cf := c.Config
	if c.didConnect {
		return nil
	}

	if err := c.load(); err != nil {
		return errors.WithStack(err)
	}

	if l, err := c.plugin.Lookup("BackendConnector"); err != nil {
		return errors.Wrap(err, "Unable to look up `BackendConnector`")
	} else if connector, ok := l.(*BackendConnector); !ok {
		return errors.New("Unable to type assert `BackendConnector`")
	} else {
		cf.GetLogger().Info("Successfully loaded database plugin")
		c.connector = *connector
		cf.GetLogger().Debugf("Address of database plugin is: %p", connector)
		RegisterBackend(c.connector)
	}
	return nil
}
