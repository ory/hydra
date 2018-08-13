//
// Copyright © 2015-2018 Gorka Lerchundi Osa <glertxundi@gmail.com>
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
// @author		Gorka Lerchundi Osa <glertxundi@gmail.com>
// @copyright 	2015-2018 Gorka Lerchundi Osa <glertxundi@gmail.com>
// @license 	Apache-2.0
//

// +build noplugin

package config

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type PluginConnection struct {
	Config *Config
	Logger logrus.FieldLogger
}

func (c *PluginConnection) Load() error {
	return errors.New("config: unable to load plugin connection because 'noplugin' tag was declared")
}
