// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"encoding/base64"
	"net/url"
)

func BasicAuth(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(url.QueryEscape(username) + ":" + url.QueryEscape(password)))
}
