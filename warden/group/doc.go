// Package group offers capabilities for grouping subjects together, making policy management easier.
//
// This endpoint is **experimental**, use it at your own risk.
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
