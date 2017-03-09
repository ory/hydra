package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/negroni"
	"github.com/meatballhat/negroni-logrus"
)

func main() {
	r := http.NewServeMux()
	r.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "success!\n")
	})

	n := negroni.New()

	nl := negronilogrus.NewMiddleware()
	// override the default Before
	nl.Before = customBefore
	// override the default After
	nl.After = customAfter

	n.Use(nl)
	n.UseHandler(r)

	n.Run(":9999")
}

func customBefore(entry *logrus.Entry, _ *http.Request, remoteAddr string) *logrus.Entry {
	return entry.WithFields(logrus.Fields{
		"REMOTE_ADDR": remoteAddr,
		"YELLING":     true,
	})
}

func customAfter(entry *logrus.Entry, res negroni.ResponseWriter, latency time.Duration, name string) *logrus.Entry {
	fields := logrus.Fields{
		"ALL_DONE":        true,
		"RESPONSE_STATUS": res.Status(),

		fmt.Sprintf("%s_LATENCY", strings.ToUpper(name)): latency,
	}

	// one way to replace an existing entry key
	if requestId, ok := entry.Data["request_id"]; ok {
		fields["REQUEST_ID"] = requestId
		delete(entry.Data, "request_id")
	}

	return entry.WithFields(fields)
}
