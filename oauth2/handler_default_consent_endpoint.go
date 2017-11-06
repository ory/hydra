// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
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

package oauth2

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) DefaultConsentHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h.L.Warnln("It looks like no consent endpoint was set. All OAuth2 flows except client credentials will fail.")

	w.Write([]byte(`
<html>
<head>
	<title>Misconfigured consent endpoint</title>
</head>
<body>
<p>
	It looks like you forgot to set the consent endpoint url, which can be set using the <code>CONSENT_URL</code>
	environment variable.
</p>
<p>
	If you are an administrator, please read <a href="https://ory-am.gitbooks.io/hydra/content/oauth2.html">
	the guide</a> to understand what you need to do. If you are a user, please contact the administrator.
</p>
</body>
</html>
`))
}
