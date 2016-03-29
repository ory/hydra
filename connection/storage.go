package connection

// Storage defines an interface for storing connections.
type Storage interface {
	// Create a new connection.
	Create(c Connection) error

	// Delete an existing connection.
	Delete(id string) error

	// Get an existing connection.
	Get(id string) (Connection, error)

	FindAllByLocalSubject(subject string) ([]Connection, error)

	FindByRemoteSubject(connector, subject string) (Connection, error)
}
