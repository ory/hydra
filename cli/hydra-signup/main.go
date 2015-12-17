package main

import (
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/sign-up", func(w http.ResponseWriter, r *http.Request) {
		/**
		 * This file is just for demonstrations sake.
		 */
		renderTemplate(w)
	})
	http.ListenAndServe(":3001", nil)
}

func renderTemplate(w http.ResponseWriter) {
	tmpl := template.New("page")
	var err error
	tmpl, err = tmpl.Parse(`<h1>Sign up</h1><p>This would be the page where you show your sign up form. You could then create the account using Hydra's /accounts endpoint.</p>`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	tmpl.Execute(w, struct{}{})
}
