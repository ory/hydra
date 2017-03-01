package cli

import (
	"fmt"

	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/hydra/warden/group"
	"github.com/spf13/cobra"
)

type GroupHandler struct {
	Config *config.Config
	M      *group.HTTPManager
}

func newGroupHandler(c *config.Config) *GroupHandler {
	return &GroupHandler{
		Config: c,
		M:      &group.HTTPManager{},
	}
}

func (h *GroupHandler) CreateGroup(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Print(cmd.UsageString())
		return
	}

	var err error
	h.M.Dry, _ = cmd.Flags().GetBool("dry")
	h.M.Endpoint = h.Config.Resolve("/warden/groups")
	h.M.Client = h.Config.OAuth2Client(cmd)

	cc := &group.Group{ID: args[0]}
	err = h.M.CreateGroup(cc)
	if h.M.Dry {
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

	var err error
	h.M.Dry, _ = cmd.Flags().GetBool("dry")
	h.M.Endpoint = h.Config.Resolve("/warden/groups")
	h.M.Client = h.Config.OAuth2Client(cmd)

	err = h.M.DeleteGroup(args[0])
	if h.M.Dry {
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

	var err error
	h.M.Dry, _ = cmd.Flags().GetBool("dry")
	h.M.Endpoint = h.Config.Resolve("/warden/groups")
	h.M.Client = h.Config.OAuth2Client(cmd)

	err = h.M.AddGroupMembers(args[0], args[1:])
	if h.M.Dry {
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

	var err error
	h.M.Dry, _ = cmd.Flags().GetBool("dry")
	h.M.Endpoint = h.Config.Resolve("/warden/groups")
	h.M.Client = h.Config.OAuth2Client(cmd)

	err = h.M.RemoveGroupMembers(args[0], args[1:])
	if h.M.Dry {
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

	h.M.Dry, _ = cmd.Flags().GetBool("dry")
	h.M.Endpoint = h.Config.Resolve("/warden/groups")
	h.M.Client = h.Config.OAuth2Client(cmd)

	gn, err := h.M.FindGroupNames(args[0])
	if h.M.Dry {
		fmt.Printf("%s\n", err)
		return
	}

	pkg.Must(err, "Could not find groups: %s", err)
	fmt.Printf("Subject %s belongs to groups %v.\n", args[0], gn)
}
