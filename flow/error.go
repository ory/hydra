package flow

import "github.com/ory/fosite"

var ErrorLogoutFlowExpired = fosite.ErrRequestUnauthorized.WithHint("The logout request has expired, please try the flow again.")
