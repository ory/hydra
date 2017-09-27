package health

// A list of clients.
// swagger:response healthStatus
type swaggerListClientsResult struct {
	// in: body
	Body struct {
		// Status always contains "ok"
		Status string `json:"status"`
	}
}
