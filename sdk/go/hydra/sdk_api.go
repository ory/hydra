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

package hydra

import (
	"github.com/ory/hydra/sdk/go/hydra/swagger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// SDK helps developers interact with ORY Hydra using a Go API.
type SDK interface {
	GetOAuth2ClientConfig() *clientcredentials.Config
	GetOAuth2Config() *oauth2.Config

	PolicyAPI
	WardenAPI
	JWKApi
	OAuth2API
}

// PolicyAPI offers capabilities for policy management.
type PolicyAPI interface {
	// CreatePolicy creates a policy. If successful, error is nil and the response status code is http.StatusCreated.
	//
	//
	//	import "github.com/ory/hydra/sdk/go/hydra/swagger"
	//
	//  policy, response, err := sdk.CreatePolicy(swagger.Policy{
	//  	Subjects: []string{"foo", "bar"},
	//  	Resources: []string{"foo", "bar"},
	//  	// ...
	// 	})
	//	if err != nil {
	//		// handle network error
	//	} else if response.StatusCode != http.StatusCreated {
	//		// handle application error
	// 	}
	//
	//	fmt.Printf("Policy created: %+v", policy)
	CreatePolicy(body swagger.Policy) (*swagger.Policy, *swagger.APIResponse, error)

	// DeletePolicy deletes a policy. If successful, error is nil and the response status code is http.StatusNoContent.
	//
	//
	//	import "github.com/ory/hydra/sdk/go/hydra/swagger"
	//
	//	id := "1234"
	//  response, err := sdk.DeletePolicy(id)
	//	if err != nil {
	//		// handle network error
	//	} else if response.StatusCode != http.StatusNoContent {
	//		// handle application error
	// 	}
	//
	//	fmt.Printf("Policy created: %s", id)
	DeletePolicy(id string) (*swagger.APIResponse, error)

	// GetPolicy returns a policy. If successful, error is nil and the response status code is http.StatusOK.
	//
	//
	//	import "github.com/ory/hydra/sdk/go/hydra/swagger"
	//
	//	id := "1234"
	//  policy, response, err := sdk.GetPolicy(id)
	//	if err != nil {
	//		// handle network error
	//	} else if response.StatusCode != http.StatusOK {
	//		// handle application error
	// 	}
	//
	//	fmt.Printf("Policy received: %+v", policy)
	GetPolicy(id string) (*swagger.Policy, *swagger.APIResponse, error)

	// ListPolicies returns a policy slice given an offset and a limit. If successful, error is nil and the response status code is http.StatusOK.
	//
	//
	//	import "github.com/ory/hydra/sdk/go/hydra/swagger"
	//
	//  policies, response, err := sdk.ListPolicies(0, 100)
	//	if err != nil {
	//		// handle network error
	//	} else if response.StatusCode != http.StatusOK {
	//		// handle application error
	// 	}
	//
	//	fmt.Printf("Policies received: %+v", policies)
	ListPolicies(offset int64, limit int64) ([]swagger.Policy, *swagger.APIResponse, error)

	// UpdatePolicy updates a policy. If successful, error is nil and the response status code is http.StatusOK.
	//
	//
	//	import "github.com/ory/hydra/sdk/go/hydra/swagger"
	//
	//  policies, response, err := sdk.ListPolicies("1234", swagger.Policy{
	//		ID: "1234",
	//  	Subjects: []string{"foo", "bar"},
	//  	Resources: []string{"foo", "bar"},
	//  	// ...
	// 	})
	//	if err != nil {
	//		// handle network error
	//	} else if response.StatusCode != http.StatusOK {
	//		// handle application error
	// 	}
	//
	//	fmt.Printf("Policy updated: %+v", policies)
	UpdatePolicy(id string, body swagger.Policy) (*swagger.Policy, *swagger.APIResponse, error)
}

type WardenAPI interface {
	AddMembersToGroup(id string, body swagger.GroupMembers) (*swagger.APIResponse, error)
	CreateGroup(body swagger.Group) (*swagger.Group, *swagger.APIResponse, error)
	DeleteGroup(id string) (*swagger.APIResponse, error)
	DoesWardenAllowAccessRequest(body swagger.WardenAccessRequest) (*swagger.WardenAccessRequestResponse, *swagger.APIResponse, error)
	DoesWardenAllowTokenAccessRequest(body swagger.WardenTokenAccessRequest) (*swagger.WardenTokenAccessRequestResponse, *swagger.APIResponse, error)
	FindGroupsByMember(member string) ([]swagger.Group, *swagger.APIResponse, error)
	GetGroup(id string) (*swagger.Group, *swagger.APIResponse, error)
	RemoveMembersFromGroup(id string, body swagger.GroupMembers) (*swagger.APIResponse, error)
}

type JWKApi interface {
	CreateJsonWebKeySet(set string, body swagger.JsonWebKeySetGeneratorRequest) (*swagger.JsonWebKeySet, *swagger.APIResponse, error)
	DeleteJsonWebKey(kid string, set string) (*swagger.APIResponse, error)
	DeleteJsonWebKeySet(set string) (*swagger.APIResponse, error)
	GetJsonWebKey(kid string, set string) (*swagger.JsonWebKeySet, *swagger.APIResponse, error)
	GetJsonWebKeySet(set string) (*swagger.JsonWebKeySet, *swagger.APIResponse, error)
	UpdateJsonWebKey(kid string, set string, body swagger.JsonWebKey) (*swagger.JsonWebKey, *swagger.APIResponse, error)
	UpdateJsonWebKeySet(set string, body swagger.JsonWebKeySet) (*swagger.JsonWebKeySet, *swagger.APIResponse, error)
}

type OAuth2API interface {
	AcceptOAuth2ConsentRequest(id string, body swagger.ConsentRequestAcceptance) (*swagger.APIResponse, error)
	CreateOAuth2Client(body swagger.OAuth2Client) (*swagger.OAuth2Client, *swagger.APIResponse, error)
	DeleteOAuth2Client(id string) (*swagger.APIResponse, error)
	GetOAuth2Client(id string) (*swagger.OAuth2Client, *swagger.APIResponse, error)
	GetOAuth2ConsentRequest(id string) (*swagger.OAuth2ConsentRequest, *swagger.APIResponse, error)
	GetWellKnown() (*swagger.WellKnown, *swagger.APIResponse, error)
	IntrospectOAuth2Token(token string, scope string) (*swagger.OAuth2TokenIntrospection, *swagger.APIResponse, error)
	ListOAuth2Clients() ([]swagger.OAuth2Client, *swagger.APIResponse, error)
	RejectOAuth2ConsentRequest(id string, body swagger.ConsentRequestRejection) (*swagger.APIResponse, error)
	RevokeOAuth2Token(token string) (*swagger.APIResponse, error)
	UpdateOAuth2Client(id string, body swagger.OAuth2Client) (*swagger.OAuth2Client, *swagger.APIResponse, error)
}
