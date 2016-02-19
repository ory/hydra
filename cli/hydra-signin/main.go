package main

import (
	"encoding/json"
	"fmt"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/ory-am/common/env"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/pborman/uuid"
	"html/template"
	"net/http"
)

type data struct {
	User     string `json:"subject"`
	Redirect string
}

func main() {
	http.HandleFunc("/sign-in", func(w http.ResponseWriter, r *http.Request) {
		/**
		 * Normally, you would authenticate the user through the OAuth2 Password Grant as described here: https://aaronparecki.com/articles/2012/07/29/1/oauth2-simplified#others
		 * If you receive a token, the user authenticated successfully, if you receive an error the request was malformed or the user credentials where wrong.
		 *
		 * This file is just for demonstrations sake.
		 */

		code := uuid.New()
		state := r.URL.Query().Get("state")

		d := data{
			User:     env.Getenv("ACCOUNT_ID", "not-provided"),
			Redirect: fmt.Sprintf("%s?code=%s&state=%s", r.URL.Query().Get("redirect_uri"), code, state),
		}

		if verify := r.URL.Query().Get("verify"); verify != "" {
			out, _ := json.Marshal(&d)
			w.Write(out)
		} else {
			renderTemplate(w, d)
		}
	})

	http.HandleFunc("/authenticate/callback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Congratulations! You have received the authorization code: %s", r.URL.Query().Get("code"))
	})

	http.ListenAndServe(":3000", nil)
}

func renderTemplate(w http.ResponseWriter, p data) {
	tmpl := template.New("page")
	var err error
	tmpl, err = tmpl.Parse(`<h1>Sign in</h1>Hello {{.User}}!<br><a href="{{.Redirect}}">Press this link to log in</a>`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	tmpl.Execute(w, p)
}
