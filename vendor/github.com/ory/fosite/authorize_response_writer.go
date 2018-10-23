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

package fosite

import (
	"context"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

func (f *Fosite) NewAuthorizeResponse(ctx context.Context, ar AuthorizeRequester, session Session) (AuthorizeResponder, error) {
	var resp = &AuthorizeResponse{
		Header:   http.Header{},
		Query:    url.Values{},
		Fragment: url.Values{},
	}

	ar.SetSession(session)
	for _, h := range f.AuthorizeEndpointHandlers {
		if err := h.HandleAuthorizeEndpointRequest(ctx, ar, resp); err != nil {
			return nil, err
		}
	}

	if !ar.DidHandleAllResponseTypes() {
		return nil, errors.WithStack(ErrUnsupportedResponseType)
	}

	return resp, nil
}
