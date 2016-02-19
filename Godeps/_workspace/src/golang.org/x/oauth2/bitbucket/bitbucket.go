// Copyright 2015 The oauth2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bitbucket provides constants for using OAuth2 to access Bitbucket.
package bitbucket

import (
	"github.com/ory-am/hydra/Godeps/_workspace/src/golang.org/x/oauth2"
)

// Endpoint is Bitbucket's OAuth 2.0 endpoint.
var Endpoint = oauth2.Endpoint{
	AuthURL:  "https://bitbucket.org/site/oauth2/authorize",
	TokenURL: "https://bitbucket.org/site/oauth2/access_token",
}
