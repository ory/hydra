package fosite

// Client represents a client or an app.
type Client interface {
	// GetID returns the client ID.
	GetID() string

	// GetHashedSecret returns the hashed secret as it is stored in the store.
	GetHashedSecret() []byte

	// Returns the client's allowed redirect URIs.
	GetRedirectURIs() []string

	// Returns the client's allowed grant types.
	GetGrantTypes() Arguments

	// Returns the client's allowed response types.
	GetResponseTypes() Arguments

	// Returns the scopes this client is allowed to request.
	GetScopes() Arguments

	// IsPublic returns true, if this client is marked as public.
	IsPublic() bool
}

// DefaultClient is a simple default implementation of the Client interface.
type DefaultClient struct {
	ID            string   `json:"id"`
	Secret        []byte   `json:"client_secret,omitempty"`
	RedirectURIs  []string `json:"redirect_uris"`
	GrantTypes    []string `json:"grant_types"`
	ResponseTypes []string `json:"response_types"`
	Scopes        []string `json:"scopes"`
	Public        bool     `json:"public"`
}

func (c *DefaultClient) GetID() string {
	return c.ID
}

func (c *DefaultClient) IsPublic() bool {
	return c.Public
}

func (c *DefaultClient) GetRedirectURIs() []string {
	return c.RedirectURIs
}

func (c *DefaultClient) GetHashedSecret() []byte {
	return c.Secret
}

func (c *DefaultClient) GetScopes() Arguments {
	return c.Scopes
}

func (c *DefaultClient) GetGrantTypes() Arguments {
	// https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata
	//
	// JSON array containing a list of the OAuth 2.0 Grant Types that the Client is declaring
	// that it will restrict itself to using.
	// If omitted, the default is that the Client will use only the authorization_code Grant Type.
	if len(c.GrantTypes) == 0 {
		return Arguments{"authorization_code"}
	}
	return Arguments(c.GrantTypes)
}

func (c *DefaultClient) GetResponseTypes() Arguments {
	// https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata
	//
	// <JSON array containing a list of the OAuth 2.0 response_type values that the Client is declaring
	// that it will restrict itself to using. If omitted, the default is that the Client will use
	// only the code Response Type.
	if len(c.ResponseTypes) == 0 {
		return Arguments{"code"}
	}
	return Arguments(c.ResponseTypes)
}
