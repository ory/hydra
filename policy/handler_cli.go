package policy

import (
	"github.com/ory-am/hydra/config"
	"github.com/spf13/cobra"
	"github.com/ory-am/hydra/pkg"
	"fmt"
	"os"
	"github.com/square/go-jose/json"
	"github.com/ory-am/ladon"
)

type CLIHandler struct {
	Config *config.Config
	M      *HTTPManager
}

func NewCLIHandler(c *config.Config) *CLIHandler {
	return &CLIHandler{
		Config: c,
		M: &HTTPManager{},
	}
}

func (h *CLIHandler) CreatePolicy(cmd *cobra.Command, args []string) {
	h.M.Endpoint = h.Config.Resolve("/policies")
	h.M.Client = h.Config.OAuth2Client(cmd)

	files, _ := cmd.Flags().GetStringSlice("files")
	if len(files) > 0 {
		for _, path := range files {
			reader, err := os.Open(path)
			pkg.Must(err, "Could not open file %s: %s", path, err)
			var policy ladon.DefaultPolicy
			err = json.NewDecoder(reader).Decode(&policy)
			pkg.Must(err, "Could not parse JSON: %s", err)
			err = h.M.Create(&policy)
			pkg.Must(err, "Could not create policy: %s", err)
			fmt.Printf("Imported policy %s from %s.\n", policy.ID, path)
		}
		return
	}

	id, _ := cmd.Flags().GetString("id")
	description, _ := cmd.Flags().GetString("description")
	subjects, _ := cmd.Flags().GetStringSlice("subjects")
	resources, _ := cmd.Flags().GetStringSlice("resources")
	actions, _ := cmd.Flags().GetStringSlice("actions")
	isAllow, _ := cmd.Flags().GetBool("allow")
	if len(subjects) == 0 || len(resources) == 0 || len(actions) == 0 {
		fmt.Println(cmd.UsageString())
		fmt.Println("")
		fmt.Println("Got empty subject, resource or action list")
	}

	effect := ladon.DenyAccess
	if isAllow {
		effect = ladon.AllowAccess
	}

	policy := &ladon.DefaultPolicy{
		ID: id,
		Description: description,
		Subjects: subjects,
		Resources: resources,
		Actions: actions,
		Effect: effect,
	}
	err := h.M.Create(policy)
	pkg.Must(err, "Could not create policy: %s", err)
	fmt.Printf("Created policy %s.\n", policy.ID)

}

func (h *CLIHandler) AddResourceToPolicy(cmd *cobra.Command, args []string) {
	fmt.Println("Not yet implemented.")
}

func (h *CLIHandler) RemoveResourceFromPolicy(cmd *cobra.Command, args []string) {
	fmt.Println("Not yet implemented.")
}

func (h *CLIHandler) AddSubjectToPolicy(cmd *cobra.Command, args []string) {
	fmt.Println("Not yet implemented.")
}

func (h *CLIHandler) RemoveSubjectFromPolicy(cmd *cobra.Command, args []string) {
	fmt.Println("Not yet implemented.")
}

func (h *CLIHandler) AddActionToPolicy(cmd *cobra.Command, args []string) {
	fmt.Println("Not yet implemented.")
}

func (h *CLIHandler) RemoveActionFromPolicy(cmd *cobra.Command, args []string) {
	fmt.Println("Not yet implemented.")
}

func (h *CLIHandler) GetPolicy(cmd *cobra.Command, args []string) {
	h.M.Endpoint = h.Config.Resolve("/policies")
	h.M.Client = h.Config.OAuth2Client(cmd)

	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	policy, err := h.M.Get(args[0])
	pkg.Must(err, "Could not delete policy: %s", err)

	out, err := json.MarshalIndent(policy, "", "\t")
	pkg.Must(err, "Could not retrieve policy: %s", err)

	fmt.Printf("%s\n", out)
}


func (h *CLIHandler) DeletePolicy(cmd *cobra.Command, args []string) {
	h.M.Endpoint = h.Config.Resolve("/policies")
	h.M.Client = h.Config.OAuth2Client(cmd)

	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	for _, arg := range args {
		err := h.M.Delete(arg)
		pkg.Must(err, "Could not delete policy: %s", err)
		fmt.Printf("Connection %s deleted.\n", arg)
	}
}
