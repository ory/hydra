package connection

import (
	"github.com/go-errors/errors"
)

var ErrNotFound = errors.New("not found")

// Connection connects an subject S with a token T issued by provider P
type Connection interface {
	// GetID returns the connection's unique identifier.
	GetID() string

	// GetProvider returns the connection's provider, for example "Google".
	GetProvider() string

	// GetLocalSubject returns the connection's local subject, for example "peter".
	GetLocalSubject() string

	// GetRemoteSubject returns the connection's remote subject, for example "peter@gmail.com".
	GetRemoteSubject() string
}

// DefaultConnection is a default implementation of the Connection interface
type DefaultConnection struct {
	ID            string
	Provider      string
	LocalSubject  string
	RemoteSubject string
}

func (c *DefaultConnection) GetID() string {
	return c.ID
}

func (c *DefaultConnection) GetProvider() string {
	return c.Provider
}

func (c *DefaultConnection) GetLocalSubject() string {
	return c.LocalSubject
}

func (c *DefaultConnection) GetRemoteSubject() string {
	return c.RemoteSubject
}
