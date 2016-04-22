package group

type Group struct {
	ID string `json:"id"`
	Members []string `json:"members"`
	Policies []string `json:"policies"`
}

type Groups map[string]Group

type Manager interface {
	GetGroup(id string) (Groups, error)

	JoinGroup(user, group string) error

	LeaveGroup(user, group string) error

	CreateGroup(group Group) error

	AddGroupPolicy(policy, group string)	 error

	RemoveGroupPolicy(policy, group string) error

	DeleteGroup(id string) error
}
