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
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package consent

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/ory/x/dbal"
)

var Migrations = map[string]*dbal.PackrMigrationSource{
	"mysql": dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{
		"migrations/sql/shared",
		"migrations/sql/mysql",
	}, true),
	"postgres": dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{
		"migrations/sql/shared",
		"migrations/sql/postgres",
	}, true),
	"cockroach": dbal.NewMustPackerMigrationSource(logrus.New(), AssetNames(), Asset, []string{
		"migrations/sql/cockroach",
	}, true),
}

var sqlParamsAuthenticationRequestHandled = []string{
	"challenge",
	"subject",
	"remember",
	"remember_for",
	"error",
	"requested_at",
	"authenticated_at",
	"acr",
	"was_used",
	"context",
	"forced_subject_identifier",
}

var sqlParamsAuthenticationRequest = []string{
	"challenge",
	"verifier",
	"client_id",
	"subject",
	"request_url",
	"skip",
	"requested_scope",
	"requested_at_audience",
	"authenticated_at",
	"requested_at",
	"csrf",
	"oidc_context",
	"login_session_id",
}

var sqlParamsConsentRequest = append(sqlParamsAuthenticationRequest,
	"forced_subject_identifier",
	"login_challenge",
	"acr",
	"context",
)

var sqlParamsConsentRequestHandled = []string{
	"challenge",
	"granted_scope",
	"granted_at_audience",
	"remember",
	"remember_for",
	"authenticated_at",
	"error",
	"requested_at",
	"session_access_token",
	"session_id_token",
	"was_used",
	"handled_at",
}
var sqlParamsConsentRequestHandledUpdate = func() []string {
	p := make([]string, len(sqlParamsConsentRequestHandled))
	for i, v := range sqlParamsConsentRequestHandled {
		p[i] = fmt.Sprintf("%s=:%s", v, v)
	}
	return p
}()

var sqlParamsAuthSession = []string{
	"id",
	"authenticated_at",
	"subject",
	"remember",
}

var sqlParamsLogoutRequest = []string{
	"challenge",
	"verifier",
	"subject",
	"sid",
	"request_url",
	"redir_url",
	"was_used",
	"accepted",
	"rejected",
	"client_id",
	"rp_initiated",
}
