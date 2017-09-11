package cli

import (
	"fmt"

	"github.com/ory/hydra/config"
	"github.com/ory/hydra/pkg"
	"github.com/ory/hydra/warden/group"
	"github.com/spf13/cobra"
)

type GroupHandler struct {
	Config *config.Config
}

func (h *GroupHandler) newGroupManager(cmd *cobra.Command) *group.HTTPManager {
	dry, _ := cmd.Flags().GetBool("dry")
	term, _ := cmd.Flags().GetBool("fake-tls-termination")

	return &group.HTTPManager{
		Dry:                dry,
		Endpoint:           h.Config.Resolve("/warden/groups"),
		Client:             h.Config.OAuth2Client(cmd),
		FakeTLSTermination: term,
	}
}

func newGroupHandler(c *config.Config) *GroupHandler {
	return &GroupHandler{
		Config: c,
	}
}

func (h *GroupHandler) CreateGroup(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Print(cmd.UsageString())
		return
	}
	m := h.newGroupManager(cmd)

	var err error
	cc := &group.Group{ID: args[0]}
	err = m.CreateGroup(cc)
	if m.Dry {
		fmt.Printf("%s\n", err)
		return
	}

	pkg.Must(err, "Could not create group: %s", err)
	fmt.Printf("Group %s created.\n", cc.ID)
}

func (h *GroupHandler) DeleteGroup(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Print(cmd.UsageString())
		return
	}
	m := h.newGroupManager(cmd)

	var err error
	err = m.DeleteGroup(args[0])
	if m.Dry {
		fmt.Printf("%s\n", err)
		return
	}

	pkg.Must(err, "Could not create group: %s", err)
	fmt.Printf("Group %s deleted.\n", args[0])
}

func (h *GroupHandler) AddMembers(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}
	m := h.newGroupManager(cmd)

	var err error
	err = m.AddGroupMembers(args[0], args[1:])
	if m.Dry {
		fmt.Printf("%s\n", err)
		return
	}

	pkg.Must(err, "Could not add members to group: %s", err)
	fmt.Printf("Members %v added to group %s.\n", args[1:], args[0])
}

func (h *GroupHandler) RemoveMembers(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}
	m := h.newGroupManager(cmd)

	var err error
	err = m.RemoveGroupMembers(args[0], args[1:])
	if m.Dry {
		fmt.Printf("%s\n", err)
		return
	}

	pkg.Must(err, "Could not remove members to group: %s", err)
	fmt.Printf("Members %v removed from group %s.\n", args[1:], args[0])
}

func (h *GroupHandler) FindGroups(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Print(cmd.UsageString())
		return
	}
	m := h.newGroupManager(cmd)
	gn, err := m.FindGroupNames(args[0])
	if m.Dry {
		fmt.Printf("%s\n", err)
		return
	}

	pkg.Must(err, "Could not find groups: %s", err)
	fmt.Printf("Subject %s belongs to groups %v.\n", args[0], gn)
}
