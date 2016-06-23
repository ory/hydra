package main

// Use golang.org/x/oauth2 client to test
// Open url in browser:
// http://localhost:14000/app

import (
	"fmt"
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/RangelReale/osin/example"
	"golang.org/x/oauth2"
)

func main() {
	config := osin.NewServerConfig()
	// goauth2 checks errors using status codes
	config.ErrorStatusCode = 401

	server := osin.NewServer(config, example.NewTestStorage())

	client := &oauth2.Config{
		ClientID:     "1234",
		ClientSecret: "aabbccdd",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:14000/authorize",
			TokenURL: "http://localhost:14000/token",
		},
		RedirectURL: "http://localhost:14000/appauth/code",
	}

	// Authorization code endpoint
	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		resp := server.NewResponse()
		defer resp.Close()

		if ar := server.HandleAuthorizeRequest(resp, r); ar != nil {
			if !example.HandleLoginPage(ar, w, r) {
				return
			}
			ar.Authorized = true
			server.FinishAuthorizeRequest(resp, r, ar)
		}
		if resp.IsError && resp.InternalError != nil {
			fmt.Printf("ERROR: %s\n", resp.InternalError)
		}
		osin.OutputJSON(resp, w, r)
	})

	// Access token endpoint
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		resp := server.NewResponse()
		defer resp.Close()

		if ar := server.HandleAccessRequest(resp, r); ar != nil {
			ar.Authorized = true
			server.FinishAccessRequest(resp, r, ar)
		}
		if resp.IsError && resp.InternalError != nil {
			fmt.Printf("ERROR: %s\n", resp.InternalError)
		}
		osin.OutputJSON(resp, w, r)
	})

	// Information endpoint
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		resp := server.NewResponse()
		defer resp.Close()

		if ir := server.HandleInfoRequest(resp, r); ir != nil {
			server.FinishInfoRequest(resp, r, ir)
		}
		osin.OutputJSON(resp, w, r)
	})

	// Application home endpoint
	http.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body>"))
		//w.Write([]byte(fmt.Sprintf("<a href=\"/authorize?response_type=code&client_id=1234&state=xyz&scope=everything&redirect_uri=%s\">Login</a><br/>", url.QueryEscape("http://localhost:14000/appauth/code"))))
		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Login</a><br/>", client.AuthCodeURL(""))))
		w.Write([]byte("</body></html>"))
	})

	// Application destination - CODE
	http.HandleFunc("/appauth/code", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		code := r.Form.Get("code")

		w.Write([]byte("<html><body>"))
		w.Write([]byte("APP AUTH - CODE<br/>"))
		defer w.Write([]byte("</body></html>"))

		if code == "" {
			w.Write([]byte("Nothing to do"))
			return
		}

		var jr *oauth2.Token
		var err error

		// if parse, download and parse json
		if r.Form.Get("doparse") == "1" {
			jr, err = client.Exchange(oauth2.NoContext, code)
			if err != nil {
				jr = nil
				w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", err)))
			}
		}

		// show json access token
		if jr != nil {
			w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", jr.AccessToken)))
			if jr.RefreshToken != "" {
				w.Write([]byte(fmt.Sprintf("REFRESH TOKEN: %s<br/>\n", jr.RefreshToken)))
			}
		}

		w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))

		cururl := *r.URL
		curq := cururl.Query()
		curq.Add("doparse", "1")
		cururl.RawQuery = curq.Encode()
		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Download Token</a><br/>", cururl.String())))
	})

	http.ListenAndServe(":14000", nil)
}
