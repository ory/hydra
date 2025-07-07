// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keysetpagination

import (
	"cmp"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/ssoready/hyrumtoken"
)

// Pagination Request Parameters
//
// For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination).
//
// swagger:model keysetPaginationRequestParameters
type RequestParameters struct {
	// Items per Page
	//
	// This is the number of items per page to return.
	// For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination).
	//
	// required: false
	// in: query
	// default: 250
	// min: 1
	// max: 1000
	PageSize int `json:"page_size"`

	// Next Page Token
	//
	// The next page token.
	// For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination).
	//
	// required: false
	// in: query
	PageToken string `json:"page_token"`
}

// Pagination Response Header
//
// The `Link` HTTP header contains multiple links (`first`, `next`) formatted as:
// `<https://{project-slug}.projects.oryapis.com/admin/sessions?page_size=250&page_token=>; rel="first"`
//
// For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination).
//
// swagger:model keysetPaginationResponseHeaders
type ResponseHeaders struct {
	// The Link HTTP Header
	//
	// The `Link` header contains a comma-delimited list of links to the following pages:
	//
	// - first: The first page of results.
	// - next: The next page of results.
	//
	// Pages are omitted if they do not exist. For example, if there is no next page, the `next` link is omitted. Examples:
	//
	//	</admin/sessions?page_size=250&page_token={last_item_uuid}; rel="first",/admin/sessions?page_size=250&page_token=>; rel="next"
	//
	Link string `json:"link"`
}

// SetLinkHeader adds the Link header for the page encoded by the paginator.
// It contains links to the first and next page, if one exists.
func SetLinkHeader(w http.ResponseWriter, keys [][32]byte, u *url.URL, p *Paginator) {
	size := p.Size()
	link := []string{linkPart(u, "first", p.DefaultToken().Encrypt(keys), size)}
	if !p.isLast {
		link = append(link, linkPart(u, "next", p.PageToken().Encrypt(keys), size))
	}
	w.Header().Set("Link", strings.Join(link, ","))
}

func linkPart(u *url.URL, rel, token string, size int) string {
	q := u.Query()
	q.Set("page_token", token)
	q.Set("page_size", strconv.Itoa(size))
	u.RawQuery = q.Encode()
	return fmt.Sprintf("<%s>; rel=%q", u.String(), rel)
}

// ParseQueryParams extracts the pagination options from the URL query.
func ParseQueryParams(keys [][32]byte, q url.Values) ([]Option, error) {
	var opts []Option
	if t := cmp.Or(q["page_token"]...); t != "" {
		raw, err := url.QueryUnescape(t)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		token, err := ParsePageToken(keys, raw)
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithToken(token))
	}
	if s := cmp.Or(q["page_size"]...); s != "" {
		size, err := strconv.Atoi(s)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		opts = append(opts, WithSize(size))
	}
	return opts, nil
}

// ParsePageToken parses a page token from the given raw string using the provided keys.
// It panics if no keys are provided.
func ParsePageToken(keys [][32]byte, raw string) (t PageToken, err error) {
	for i := range keys {
		err = errors.WithStack(hyrumtoken.Unmarshal(&keys[i], raw, &t))
		if err == nil {
			return
		}
	}
	// as a last resort, try the fallback key
	err = hyrumtoken.Unmarshal(fallbackEncryptionKey, raw, &t)
	return t, errors.WithStack(err)
}
