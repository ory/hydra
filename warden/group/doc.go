// Package group offers capabilities for grouping subjects together, making policy management easier.
package group

// A list of groups the member is belonging to
// swagger:response findGroupsByMemberResponse
type swaggerFindGroupsByMemberResponse struct {
	// in: body
	Body []string
}

// swagger:parameters findGroupsByMember
type swaggerFindGroupsByMemberParameters struct {
	// The id of the member to look up.
	// in: query
	Member int `json:"member"`
}

// swagger:parameters getGroup
type swaggerGetGroupParameters struct {
	// The id of the group to look up.
	// in: path
	ID int `json:"id"`
}

// swagger:parameters removeMembersFromGroup addMembersToGroup
type swaggerModifyMembersParameters struct {
	// The id of the group to modify.
	// in: path
	ID int `json:"id"`

	// in: body
	Body membersRequest
}

// A group
// swagger:response groupResponse
type swaggerGroupResponse struct {
	// in: body
	Body Group
}
