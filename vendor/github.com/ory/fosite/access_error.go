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
	"encoding/json"
	"fmt"
	"net/http"
)

func (f *Fosite) WriteAccessError(rw http.ResponseWriter, _ AccessRequester, err error) {
	f.writeJsonError(rw, err)
}

func (f *Fosite) writeJsonError(rw http.ResponseWriter, err error) {
	rw.Header().Set("Content-Type", "application/json;charset=UTF-8")

	rfcerr := ErrorToRFC6749Error(err)
	if !f.SendDebugMessagesToClients {
		rfcerr.Debug = ""
	}

	js, err := json.Marshal(rfcerr)
	if err != nil {
		http.Error(rw, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(rfcerr.Code)
	rw.Write(js)
}
