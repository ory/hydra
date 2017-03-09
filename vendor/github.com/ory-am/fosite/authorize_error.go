package fosite

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func (c *Fosite) WriteAuthorizeError(rw http.ResponseWriter, ar AuthorizeRequester, err error) {
	rfcerr := ErrorToRFC6749Error(err)

	if !ar.IsRedirectURIValid() {
		js, err := json.MarshalIndent(rfcerr, "", "\t")
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(rfcerr.StatusCode)
		rw.Write(js)
		return
	}

	redirectURI := ar.GetRedirectURI()
	query := url.Values{}
	query.Add("error", rfcerr.Name)
	query.Add("error_description", rfcerr.Description)
	query.Add("state", ar.GetState())

	if ar.GetResponseTypes().Exact("token") || len(ar.GetResponseTypes()) > 1 {
		redirectURI.Fragment = query.Encode()
	} else {
		for key, values := range redirectURI.Query() {
			for _, value := range values {
				query.Add(key, value)
			}
		}
		redirectURI.RawQuery = query.Encode()
	}

	rw.Header().Add("Location", redirectURI.String())
	rw.WriteHeader(http.StatusFound)
}
