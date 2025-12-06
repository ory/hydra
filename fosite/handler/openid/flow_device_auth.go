// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openid

import (
	"context"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/rfc8628"
)

// OpenIDConnectDeviceHandler a response handler for the Device Authorization Grant with OpenID Connect identity layer
type OpenIDConnectDeviceHandler struct {
	Storage  OpenIDConnectRequestStorageProvider
	Strategy rfc8628.DeviceCodeStrategyProvider
	Config   interface {
		fosite.IDTokenLifespanProvider
	}
	*IDTokenHandleHelper
}

func (c *OpenIDConnectDeviceHandler) HandleDeviceEndpointRequest(ctx context.Context, dar fosite.DeviceRequester, resp fosite.DeviceResponder) error {
	// We don't want to create the openid session on this call, because we don't know if the user
	// will actually complete the flow and give consent. The implementer MUST call the CreateOpenIDConnectSession
	// methods when the user logs in to instantiate the session.
	return nil
}
