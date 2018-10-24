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

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/ory/x/sqlcon"
)

type SQLBackend struct {
	db *sqlx.DB
	l  logrus.FieldLogger
	Options
}

func init() {
	RegisterBackend(&SQLBackend{})
}

func (s *SQLBackend) Init(url string, l logrus.FieldLogger, opts ...ConnectorOptions) error {
	for _, opt := range opts {
		opt(&s.Options)
	}

	sqlconOptions := []sqlcon.Opt{}

	if s.UseTracing {
		sqlconOptions = append(sqlconOptions, sqlcon.WithDistributedTracing())
	}

	if s.useRandomDriverName {
		sqlconOptions = append(sqlconOptions, sqlcon.WithRandomDriverName())
	}

	if s.omitSQLArgsFromSpans {
		sqlconOptions = append(sqlconOptions, sqlcon.WithOmitArgsFromTraceSpans())
	}

	if s.allowRootTracingSpans {
		sqlconOptions = append(sqlconOptions, sqlcon.WithAllowRoot())
	}

	connection, err := sqlcon.NewSQLConnection(url, l, sqlconOptions...)
	if err != nil {
		return err
	}
	s.l = l
	s.db = connection.GetDatabase()
	return nil
}

func (s *SQLBackend) NewConsentManager(clientManager client.Manager, fs pkg.FositeStorer) consent.Manager {
	expectDependency(s.l, clientManager, s.db, fs)
	return consent.NewSQLManager(
		s.db,
		clientManager,
		fs,
	)
}

func (s *SQLBackend) NewOAuth2Manager(clientManager client.Manager, accessTokenLifespan time.Duration, tokenStrategy string) pkg.FositeStorer {
	expectDependency(s.l, clientManager, s.db)
	return oauth2.NewFositeSQLStore(clientManager, s.db, s.l, accessTokenLifespan, tokenStrategy == "jwt")
}

func (s *SQLBackend) NewClientManager(hasher fosite.Hasher) client.Manager {
	expectDependency(s.l, hasher, s.db)
	return &client.SQLManager{
		DB:     s.db,
		Hasher: hasher,
	}
}

func (s *SQLBackend) NewJWKManager(cipher *jwk.AEAD) jwk.Manager {
	expectDependency(s.l, cipher, s.db)
	return &jwk.SQLManager{
		DB:     s.db,
		Cipher: cipher,
	}
}

func (s *SQLBackend) Prefixes() []string {
	return []string{"mysql", "postgres"}
}

func (s *SQLBackend) Ping() error {
	expectDependency(s.l, s.db)
	return s.db.Ping()
}
