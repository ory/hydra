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

package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/tomnomnom/linkheader"
)

var osExit = os.Exit

func fatal(message string, args ...interface{}) {
	fmt.Printf(message+"\n", args...)
	osExit(1)
}

func getPageToken(resp *http.Response) string {
	for _, link := range linkheader.Parse(resp.Header.Get("Link")) {
		if link.Rel != "next" {
			continue
		}

		parsed, err := url.Parse(link.URL)
		if err != nil {
			continue
		}

		if pageToken := parsed.Query().Get("page_token"); len(pageToken) > 0 {
			return pageToken
		}
	}

	return ""
}
