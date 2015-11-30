// Package storage defines an interface, which all osin-storage implementations are going to support.
package storage

import "github.com/RangelReale/osin"

// Storage extends github.com/RangelReale/osin.Storage with create, update and delete methods for clients.
type Storage interface {
	osin.Storage

	// CreateClient stores the client in the database and returns an error, if something went wrong.
	CreateClient(client osin.Client) error

	// UpdateClient updates the client (identified by it's id) and replaces the values with the values of client.
	// Returns an error if something went wrong.
	UpdateClient(client osin.Client) error

	// RemoveClient removes a client (identified by id) from the database. Returns an error if something went wrong.
	RemoveClient(id string) error
}
