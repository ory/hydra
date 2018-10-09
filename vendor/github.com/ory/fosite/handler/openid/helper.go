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
 *
 */

package openid

import (
	"context"

	"github.com/ory/fosite"
)

type IDTokenHandleHelper struct {
	IDTokenStrategy OpenIDConnectTokenStrategy
}

func (i *IDTokenHandleHelper) generateIDToken(ctx context.Context, fosr fosite.Requester) (token string, err error) {
	token, err = i.IDTokenStrategy.GenerateIDToken(ctx, fosr)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (i *IDTokenHandleHelper) IssueImplicitIDToken(ctx context.Context, ar fosite.Requester, resp fosite.AuthorizeResponder) error {
	token, err := i.generateIDToken(ctx, ar)
	if err != nil {
		return err
	}

	resp.AddFragment("id_token", token)
	return nil
}

func (i *IDTokenHandleHelper) IssueExplicitIDToken(ctx context.Context, ar fosite.Requester, resp fosite.AccessResponder) error {
	token, err := i.generateIDToken(ctx, ar)
	if err != nil {
		return err
	}

	resp.SetExtra("id_token", token)
	return nil
}
