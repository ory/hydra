package main

import (
	"fmt"
	"net/http"

	"flag"

	"github.com/Sirupsen/logrus"

	"github.com/urfave/negroni"
	"github.com/meatballhat/negroni-logrus"
)

var quiet bool

func main() {
	flag.BoolVar(&quiet, "quiet", false, "if set, hide messages from the logger")
	flag.Parse()
	r := http.NewServeMux()
	r.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "success!\n")
	})

	loglevel := logrus.InfoLevel

	if quiet {
		loglevel = logrus.ErrorLevel
	}

	n := negroni.New()
	n.Use(negronilogrus.NewCustomMiddleware(loglevel, &logrus.JSONFormatter{}, "web"))
	n.UseHandler(r)

	n.Run(":9999")
}
