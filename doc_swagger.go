package main

// A genericError is the error format of Hydra's RESTful APIs.
// swagger:response genericError
type genericError struct {
	// in: body
	body struct {
		     // code
		     code    int                      `json:"code,omitempty"`

		     // status
		     status  string                   `json:"status,omitempty"`
		     // request
		     request string                   `json:"request,omitempty"`

		     // reason
		     reason  string                   `json:"reason,omitempty"`

		     // details
		     details []map[string]interface{} `json:"details,omitempty"`


		     // message
		     message string                   `json:"message"`
	     }
}
