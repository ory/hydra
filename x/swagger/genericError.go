// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package swagger

import "github.com/ory/herodot"

// swagger:model genericError
//
//nolint:deadcode,unused
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type GenericError struct{ herodot.DefaultError }
