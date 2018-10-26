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
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/pkg"
)

var (
	backends         = make(map[string]BackendConnector)
	bmutex           sync.Mutex
	errNilDependency = errors.New("A dependency was expected to be defined but is nil. Please open an issue with the stack trace.")
)

type BackendConnector interface {
	Init(url string, l logrus.FieldLogger, opts ...ConnectorOptions) error
	NewConsentManager(clientManager client.Manager, fs pkg.FositeStorer) consent.Manager
	NewOAuth2Manager(clientManager client.Manager, accessTokenLifespan time.Duration, tokenStrategy string) pkg.FositeStorer
	NewClientManager(hasher fosite.Hasher) client.Manager
	NewJWKManager(cipher *jwk.AEAD) jwk.Manager
	Ping() error
	Prefixes() []string
}

func RegisterBackend(b BackendConnector) {
	bmutex.Lock()
	for _, prefix := range b.Prefixes() {
		backends[prefix] = b
	}
	bmutex.Unlock()
}

func supportedSchemes() []string {
	keys := make([]string, len(backends))
	i := 0
	for k := range backends {
		keys[i] = k
		i++
	}
	return keys
}

func expectDependency(logger logrus.FieldLogger, dependencies ...interface{}) {
	if logger == nil {
		panic("missing logger for dependency check")
	}
	for _, d := range dependencies {
		if d == nil {
			logger.WithError(errors.WithStack(errNilDependency)).Fatalf("A fatal issue occurred.")
		}
	}
}
