package webbrowser

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func init() {
	// Register a generic browser, if any, for current OS.
	if os, ok := osCommand[runtime.GOOS]; ok {
		Candidates = append(Candidates, GenericBrowser{os.cmd, os.args})
	}
}

type args struct {
	cmd  string
	args []string
}

var (
	osCommand = map[string]*args{
		"darwin":  &args{"open", nil},
		"freebsd": &args{"xdg-open", nil},
		"linux":   &args{"xdg-open", nil},
		"netbsd":  &args{"xdg-open", nil},
		"openbsd": &args{"xdg-open", nil}, // It may be open instead
		"windows": &args{"cmd", []string{"/c", "start"}},
	}
	ErrCantOpen     = errors.New("webbrowser.Open: can't open webpage")
	ErrNoCandidates = errors.New("webbrowser.Open: no browser candidate found for your OS.")
)

// List of registered `Browser`s that will be tried with Open.
var Candidates []Browser

type Browser interface {
	Open(string) error
}

// GenericBrowser just holds a command name and its arguments; the url will be
// appended as last arg. If you need to use string replacement for url define
// your own implementation.
type GenericBrowser struct {
	Cmd  string
	Args []string
}

func (gb GenericBrowser) Open(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}

	// Enforce a scheme (windows requires scheme to be set to work properly).
	if u.Scheme != "https" {
		u.Scheme = "http"
	}
	s = u.String()

	// Escape characters not allowed by cmd/bash
	switch runtime.GOOS {
	case "windows":
		s = strings.Replace(s, "&", `^&`, -1)
	}

	var cmd *exec.Cmd
	if gb.Args != nil {
		cmd = exec.Command(gb.Cmd, append(gb.Args, s)...)
	} else {
		cmd = exec.Command(gb.Cmd, s)
	}
	return cmd.Run()
}

// Open opens an URL on the first available candidate found.
func Open(s string) error {
	if len(Candidates) == 0 {
		return ErrNoCandidates
	}

	for _, b := range Candidates {
		err := b.Open(s)
		if err == nil {
			return nil
		}
	}

	// Try to determine if there's a display available (only linux) and we
	// aren't on a terminal (all but windows).
	switch runtime.GOOS {
	case "linux":
		// No display, no need to open a browser. Lynx users **MAY** have
		// something to say about this.
		if os.Getenv("DISPLAY") == "" {
			return fmt.Errorf("Tried to open %q on default webbrowser, no screen found.\n", s)
		}
		fallthrough
	case "darwin":
		// Check SSH env vars.
		if os.Getenv("SSH_CLIENT") != "" || os.Getenv("SSH_TTY") != "" {
			return fmt.Errorf("Tried to open %q on default webbrowser, but you are running a shell session.\n", s)
		}
	}

	return ErrCantOpen
}
