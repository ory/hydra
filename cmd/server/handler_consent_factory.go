// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/oauth2"
)

func injectConsentManager(c *config.Config) {
	var ctx = c.Context()
	var manager oauth2.ConsentRequestManager

	switch con := ctx.Connection.(type) {
	case *config.MemoryConnection:
		manager = oauth2.NewConsentRequestMemoryManager()
		break
	case *config.SQLConnection:
		manager = oauth2.NewConsentRequestSQLManager(con.GetDatabase())
		break
	case *config.PluginConnection:
		var err error
		if manager, err = con.NewConsentRequestManager(); err != nil {
			c.GetLogger().Fatalf("Could not load client manager plugin %s", err)
		}
		break
	default:
		panic("Unknown connection type.")
	}

	ctx.ConsentManager = manager

}

func newConsentHanlder(c *config.Config, router *httprouter.Router) *oauth2.ConsentSessionHandler {
	ctx := c.Context()
	h := &oauth2.ConsentSessionHandler{
		H: herodot.NewJSONWriter(c.GetLogger()),
		W: ctx.Warden, M: ctx.ConsentManager,
		ResourcePrefix: c.AccessControlResourcePrefix,
	}

	h.SetRoutes(router)
	return h
}
