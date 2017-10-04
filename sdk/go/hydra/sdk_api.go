package hydra

import "github.com/ory/hydra/sdk/go/hydra/swagger"

type SDK interface {
	PolicyAPI
	WardenAPI
	JWKApi
	OAuth2API
}

type PolicyAPI interface {
	CreatePolicy(body swagger.Policy) (*swagger.Policy, *swagger.APIResponse, error)
	DeletePolicy(id string) (*swagger.APIResponse, error)
	GetPolicy(id string) (*swagger.Policy, *swagger.APIResponse, error)
	ListPolicies(offset int64, limit int64) ([]swagger.Policy, *swagger.APIResponse, error)
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
	GetOAuth2ConsentRequest(id string) (*swagger.OAuth2consentRequest, *swagger.APIResponse, error)
	GetWellKnown() (*swagger.WellKnown, *swagger.APIResponse, error)
	IntrospectOAuth2Token(token string, scope string) (*swagger.OAuth2TokenIntrospection, *swagger.APIResponse, error)
	ListOAuth2Clients() ([]swagger.OAuth2Client, *swagger.APIResponse, error)
	RejectOAuth2ConsentRequest(id string, body swagger.ConsentRequestRejection) (*swagger.APIResponse, error)
	RevokeOAuth2Token(token string) (*swagger.APIResponse, error)
	UpdateOAuth2Client(id string, body swagger.OAuth2Client) (*swagger.OAuth2Client, *swagger.APIResponse, error)
}
