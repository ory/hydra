// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"context"
	"net/http"
)

func (f *Fosite) WriteAuthorizeResponse(ctx context.Context, rw http.ResponseWriter, ar AuthorizeRequester, resp AuthorizeResponder) {
	// Set custom headers, e.g. "X-MySuperCoolCustomHeader" or "X-DONT-CACHE-ME"...
	wh := rw.Header()
	rh := resp.GetHeader()
	for k := range rh {
		wh.Set(k, rh.Get(k))
	}

	wh.Set("Cache-Control", "no-store")
	wh.Set("Pragma", "no-cache")

	redir := ar.GetRedirectURI()
	switch rm := ar.GetResponseMode(); rm {
	case ResponseModeFormPost:
		//form_post
		rw.Header().Add("Content-Type", "text/html;charset=UTF-8")
		WriteAuthorizeFormPostResponse(redir.String(), resp.GetParameters(), GetPostFormHTMLTemplate(ctx, f), rw)
		return
	case ResponseModeQuery, ResponseModeDefault:
		// Explicit grants
		q := redir.Query()
		rq := resp.GetParameters()
		for k := range rq {
			q.Set(k, rq.Get(k))
		}
		redir.RawQuery = q.Encode()
		sendRedirect(redir.String(), rw)
		return
	case ResponseModeFragment:
		// Implicit grants
		// The endpoint URI MUST NOT include a fragment component.
		redir.Fragment = ""

		u := redir.String()
		fr := resp.GetParameters()
		if len(fr) > 0 {
			u = u + "#" + fr.Encode()
		}
		sendRedirect(u, rw)
		return
	default:
		if f.ResponseModeHandler(ctx).ResponseModes().Has(rm) {
			f.ResponseModeHandler(ctx).WriteAuthorizeResponse(ctx, rw, ar, resp)
			return
		}
	}
}

// https://tools.ietf.org/html/rfc6749#section-4.1.1
// When a decision is established, the authorization server directs the
// user-agent to the provided client redirection URI using an HTTP
// redirection response, or by other means available to it via the
// user-agent.
func sendRedirect(url string, rw http.ResponseWriter) {
	rw.Header().Set("Location", url)
	rw.WriteHeader(http.StatusSeeOther)
}
