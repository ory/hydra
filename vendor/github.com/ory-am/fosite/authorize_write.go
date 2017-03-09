package fosite

import (
	"net/http"
	"regexp"
)

var (
	// scopeMatch = regexp.MustCompile("scope=[^\\&]+.*$")
	plusMatch = regexp.MustCompile("\\+")
)

func (c *Fosite) WriteAuthorizeResponse(rw http.ResponseWriter, ar AuthorizeRequester, resp AuthorizeResponder) {
	redir := ar.GetRedirectURI()

	// Explicit grants
	q := redir.Query()
	rq := resp.GetQuery()
	for k := range rq {
		q.Set(k, rq.Get(k))
	}
	redir.RawQuery = q.Encode()

	// Set custom headers, e.g. "X-MySuperCoolCustomHeader" or "X-DONT-CACHE-ME"...
	wh := rw.Header()
	rh := resp.GetHeader()
	for k := range rh {
		wh.Set(k, rh.Get(k))
	}

	// Implicit grants
	redir.Fragment = resp.GetFragment().Encode()

	u := redir.String()
	u = plusMatch.ReplaceAllString(u, "%20")

	// https://tools.ietf.org/html/rfc6749#section-4.1.1
	// When a decision is established, the authorization server directs the
	// user-agent to the provided client redirection URI using an HTTP
	// redirection response, or by other means available to it via the
	// user-agent.
	wh.Set("Location", u)
	rw.WriteHeader(http.StatusFound)
}
