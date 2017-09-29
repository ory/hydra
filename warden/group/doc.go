// Package group offers capabilities for grouping subjects together, making policy management easier.
//
// This endpoint is **experimental**, use it at your own risk.
package group

// A list of groups the member is belonging to
// swagger:response findGroupsByMemberResponse
type swaggerFindGroupsByMemberResponse struct {
	// in: body
	// type: array
	Body []Group
}

// swagger:parameters findGroupsByMember
type swaggerFindGroupsByMemberParameters struct {
	// The id of the member to look up.
	// in: query
	// required: true
	Member string `json:"member"`
}

// swagger:parameters createGroup
type swaggerCreateGroupParameters struct {
	// in: body
	Body Group
}

// swagger:parameters getGroup deleteGroup
type swaggerGetGroupParameters struct {
	// The id of the group to look up.
	// in: path
	ID string `json:"id"`
}

// swagger:parameters removeMembersFromGroup addMembersToGroup
type swaggerModifyMembersParameters struct {
	// The id of the group to modify.
	// in: path
	ID string `json:"id"`

	// in: body
	Body membersRequest
}

// A group
// swagger:response groupResponse
type swaggerGroupResponse struct {
	// in: body
	Body Group
}
