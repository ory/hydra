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

	FindGroupNames(subject string) ([]string, error)
}
