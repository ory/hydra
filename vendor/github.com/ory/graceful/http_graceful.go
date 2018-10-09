package graceful

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"
)

// StarFunc is the type of the function invoked by Graceful to start the server
type StartFunc func() error

// ShutdownFunc is the type of the function invoked by Graceful to shutdown the server
type ShutdownFunc func(context.Context) error

// DefaultShutdownTimeout defines how long Graceful will wait before forcibly shutting down
var DefaultShutdownTimeout = 5 * time.Second

// Graceful sets up graceful handling of SIGINT, typically for an HTTP server. When SIGINT is trapped,
// the shutdown handler will be invoked with a context that expires after DefaultShutdownTimeout (5s).
//
//   server := graceful.WithDefaults(http.Server{})
//
//   if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
//	   log.Fatal("Failed to gracefully shut down")
//   }
func Graceful(start StartFunc, shutdown ShutdownFunc) error {
	var (
		stopChan = make(chan os.Signal)
		errChan  = make(chan error)
	)

	// Setup the graceful shutdown handler (traps SIGINT)
	go func() {
		signal.Notify(stopChan, os.Interrupt)

		<-stopChan

		timer, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
		defer cancel()

		if err := shutdown(timer); err != nil {
			errChan <- errors.WithStack(err)
			return
		}

		errChan <- nil
	}()

	// Start the server
	if err := start(); err != http.ErrServerClosed {
		return err
	}

	return <-errChan
}
