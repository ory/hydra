package connection

// Connection connects an subject S with a token T issued by provider P
type Connection struct {
	ID            string `json:"id,omitempty" gorethink:"id"`
	Provider      string `json:"provider" valid:"required" gorethink:"provider"`
	LocalSubject  string `json:"localSubject" valid:"required" gorethink:"localsubject"`
	RemoteSubject string `json:"remoteSubject" valid:"required" gorethink:"remotesubject"`
}

// GetID returns the connection's unique identifier.
func (c *Connection) GetID() string {
	return c.ID
}

// GetProvider returns the connection's provider, for example "Google".
func (c *Connection) GetProvider() string {
	return c.Provider
}

// GetLocalSubject returns the connection's local subject, for example "peter".
func (c *Connection) GetLocalSubject() string {
	return c.LocalSubject
}

// GetRemoteSubject returns the connection's remote subject, for example "peter@gmail.com".
func (c *Connection) GetRemoteSubject() string {
	return c.RemoteSubject
}
