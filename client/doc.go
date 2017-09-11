// Package client implements the OAuth 2.0 Client functionality and provides http handlers, http clients and storage adapters.
package client

// swagger:parameters createOAuthClient
type swaggerCreateClientPayload struct {
	// in: body
	// required: true
	Body Client
}

// swagger:parameters updateOAuthClient
type swaggerUpdateClientPayload struct {
	// in: path
	// required: true
	ID string `json:"id"`

	// in: body
	// required: true
	Body Client
}

// A list of clients.
// swagger:response clientsList
type swaggerListClientsResult struct {
	// in: body
	Body []Client
}

// swagger:parameters getOAuthClient deleteOAuthClient
type swaggerQueryClientPayload struct {
	// The id of the OAuth 2.0 Client.
	//
	// unique: true
	// in: path
	ID string `json:"id"`
}
