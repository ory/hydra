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

// Group represents a warden group
//
// swagger:model group
type Group struct {
	// ID is the groups id.
	ID string `json:"id"`

	// Members is who belongs to the group.
	Members []string `json:"members"`
}

type Manager interface {
	CreateGroup(*Group) error
	GetGroup(id string) (*Group, error)
	DeleteGroup(id string) error

	AddGroupMembers(group string, members []string) error
	RemoveGroupMembers(group string, members []string) error

	FindGroupsByMember(subject string) ([]Group, error)
}
