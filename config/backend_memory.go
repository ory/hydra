/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package config

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
)

type MemoryBackend struct {
	l logrus.FieldLogger
}

func init() {
	RegisterBackend(&MemoryBackend{})
}

func (m *MemoryBackend) Init(url string, l logrus.FieldLogger, _ ...ConnectorOptions) error {
	m.l = l
	return nil
}

func (m *MemoryBackend) NewConsentManager(_ client.Manager, fs pkg.FositeStorer) consent.Manager {
	expectDependency(m.l, fs)
	return consent.NewMemoryManager(fs)
}

func (m *MemoryBackend) NewOAuth2Manager(clientManager client.Manager, accessTokenLifespan time.Duration, _ string) pkg.FositeStorer {
	expectDependency(m.l, clientManager)
	return oauth2.NewFositeMemoryStore(clientManager, accessTokenLifespan)
}

func (m *MemoryBackend) NewClientManager(hasher fosite.Hasher) client.Manager {
	expectDependency(m.l, hasher)
	return client.NewMemoryManager(hasher)
}

func (m *MemoryBackend) NewJWKManager(_ *jwk.AEAD) jwk.Manager {
	return &jwk.MemoryManager{}
}

func (m *MemoryBackend) Prefixes() []string {
	return []string{"memory"}
}

func (m *MemoryBackend) Ping() error {
	return nil
}
