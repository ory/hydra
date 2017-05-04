package main

// The standard error format
// swagger:response genericError
type genericError struct {
	// in: body
	Body struct {
		Code int `json:"code,omitempty"`

		Status string `json:"status,omitempty"`

		Request string `json:"request,omitempty"`

		Reason string `json:"reason,omitempty"`

		Details []map[string]interface{} `json:"details,omitempty"`

		Message string `json:"message"`
	}
}
