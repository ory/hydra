// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"net/http"
	"net/url"

	"github.com/tomnomnom/linkheader"
)

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
