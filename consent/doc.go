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

import "time"

// swagger:parameters getLogoutRequest
type swaggerGetLogoutRequestByChallenge struct {
	// in: query
	// required: true
	Challenge string `json:"logout_challenge"`
}

// swagger:parameters rejectLogoutRequest
type swaggerRejectLogoutRequest struct {
	// in: query
	// required: true
	Challenge string `json:"logout_challenge"`

	// in: body
	Body RequestDeniedError
}

// swagger:model flushLoginConsentRequest
type FlushLoginConsentRequest struct {
	// NotAfter sets after which point tokens should not be flushed. This is useful when you want to keep a history
	// of recent login and consent database entries for auditing.
	NotAfter time.Time `json:"notAfter"`
}
