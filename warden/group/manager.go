package group

type Group struct {
	ID      string   `json:"id"`
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
