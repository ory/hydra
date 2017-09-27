// Package client implements the OAuth 2.0 Client functionality and provides http handlers, http clients and storage adapters.
package client

// swagger:parameters createOAuth2Client
type swaggerCreateClientPayload struct {
	// in: body
	// required: true
	Body Client
}

// swagger:parameters updateOAuth2Client
type swaggerUpdateClientPayload struct {
	// in: path
	// required: true
	ID string `json:"id"`

	// in: body
	// required: true
	Body Client
}

// A list of clients.
// swagger:response oAuth2ClientList
type swaggerListClientsResult struct {
	// in: body
	// type: array
	Body []Client
}

// swagger:parameters getOAuth2Client deleteOAuth2Client
type swaggerQueryClientPayload struct {
	// The id of the OAuth 2.0 Client.
	//
	// unique: true
	// in: path
	ID string `json:"id"`
}
