package cli

import (
	"fmt"

	"net/http"

	"github.com/ory/hydra/config"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/spf13/cobra"
)

type GroupHandler struct {
	Config *config.Config
}

func (h *GroupHandler) newGroupManager(cmd *cobra.Command) *hydra.WardenApi {
	client := hydra.NewWardenApiWithBasePath(h.Config.ClusterURL)
	client.Configuration.Transport = h.Config.OAuth2Client(cmd).Transport
	if term, _ := cmd.Flags().GetBool("fake-tls-termination"); term {
		client.Configuration.DefaultHeader["X-Forwarded-Proto"] = "https"
	}

	return client
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

	_, response, err := m.CreateGroup(hydra.Group{Id: args[0]})
	checkResponse(response, err, http.StatusCreated)
	fmt.Printf("Group %s created.\n", args[0])
}

func (h *GroupHandler) DeleteGroup(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Print(cmd.UsageString())
		return
	}

	m := h.newGroupManager(cmd)
	response, err := m.DeleteGroup(args[0])
	checkResponse(response, err, http.StatusNoContent)
	fmt.Printf("Group %s deleted.\n", args[0])
}

func (h *GroupHandler) AddMembers(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	m := h.newGroupManager(cmd)
	response, err := m.AddMembersToGroup(args[0], hydra.GroupMembers{Members: args[1:]})
	checkResponse(response, err, http.StatusNoContent)
	fmt.Printf("Members %v added to group %s.\n", args[1:], args[0])
}

func (h *GroupHandler) RemoveMembers(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	m := h.newGroupManager(cmd)
	response, err := m.RemoveMembersFromGroup(args[0], hydra.GroupMembers{Members: args[1:]})
	checkResponse(response, err, http.StatusNoContent)
	fmt.Printf("Members %v removed from group %s.\n", args[1:], args[0])
}

func (h *GroupHandler) FindGroups(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Print(cmd.UsageString())
		return
	}

	m := h.newGroupManager(cmd)
	groups, response, err := m.FindGroupsByMember(args[0])
	checkResponse(response, err, http.StatusOK)
	formatResponse(groups)
}
