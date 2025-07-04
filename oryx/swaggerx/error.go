// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package swaggerx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime"
)

func FormatSwaggerError(err error) string {
	var e *runtime.APIError
	if errors.As(err, &e) {
		body, err := json.MarshalIndent(e, "\t", "  ")
		if err != nil {
			body = []byte(fmt.Sprintf("%+v", e.Response))
		}

		switch e.Code {
		case http.StatusForbidden:
			return fmt.Sprintf("The service responded with status code 403 indicating that you lack permission to access the resource. The full error details are:\n\n\t%s\n\n", body)
		case http.StatusUnauthorized:
			return fmt.Sprintf("The service responded with status code 401 indicating that you forgot to include credentials (e.g. token, TLS certificate, ...) in the HTTP request.  The full error details are:\n\n\t%s\n\n", body)
		case http.StatusNotFound:
			return fmt.Sprintf("The service responded with status code 404 indicating that the resource does not exist. Check that the URL is correct (are you using the correct admin/public/... endpoint?) and that the resource exists. The full error details are:\n\n\t%s\n\n", body)
		default:
			return fmt.Sprintf("Unable to complete operation %s because the server responded with status code %d:\n\n\t%s\n", e.OperationName, e.Code, body)
		}
	}
	return fmt.Sprintf("%+v", err)
}
