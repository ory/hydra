package dockertest

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const interruptedExitCode = 130

type OnExit struct {
	sync.Mutex
	once     sync.Once
	handlers []func()
}

func NewOnExit() *OnExit {
	return &OnExit{
		handlers: make([]func(), 0),
	}
}

func (at *OnExit) Add(f func()) {
	at.Lock()
	defer at.Unlock()
	at.handlers = append(at.handlers, f)
	at.once.Do(func() {
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
			<-c
			at.Exit(interruptedExitCode)
		}()
	})
}

func (at *OnExit) Exit(status int) {
	at.execute()
	os.Exit(status)
}

func (at *OnExit) execute() {
	at.Lock()
	defer at.Unlock()
	for _, f := range at.handlers {
		f()
	}
	at.handlers = make([]func(), 0)
}
