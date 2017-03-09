package cli

import (
	"fmt"
	"os"

	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/hydra/policy"
	"github.com/ory-am/ladon"
	"github.com/spf13/cobra"
	"github.com/square/go-jose/json"
)

type PolicyHandler struct {
	Config *config.Config
	M      *policy.HTTPManager
}

func newPolicyHandler(c *config.Config) *PolicyHandler {
	return &PolicyHandler{
		Config: c,
		M:      &policy.HTTPManager{},
	}
}

func (h *PolicyHandler) CreatePolicy(cmd *cobra.Command, args []string) {
	h.M.Dry, _ = cmd.Flags().GetBool("dry")
	h.M.Endpoint = h.Config.Resolve("/policies")
	h.M.Client = h.Config.OAuth2Client(cmd)

	files, _ := cmd.Flags().GetStringSlice("files")
	if len(files) > 0 {
		for _, path := range files {
			reader, err := os.Open(path)
			pkg.Must(err, "Could not open file %s: %s", path, err)
			var p ladon.DefaultPolicy
			err = json.NewDecoder(reader).Decode(&p)
			pkg.Must(err, "Could not parse JSON: %s", err)
			err = h.M.Create(&p)
			pkg.Must(err, "Could not create policy: %s", err)
			fmt.Printf("Imported policy %s from %s.\n", p.ID, path)
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

	p := &ladon.DefaultPolicy{
		ID:          id,
		Description: description,
		Subjects:    subjects,
		Resources:   resources,
		Actions:     actions,
		Effect:      effect,
	}
	err := h.M.Create(p)
	if h.M.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not create policy: %s", err)
	fmt.Printf("Created policy %s.\n", p.ID)

}

func (h *PolicyHandler) AddResourceToPolicy(cmd *cobra.Command, args []string) {
	h.M.Dry, _ = cmd.Flags().GetBool("dry")
	h.M.Endpoint = h.Config.Resolve("/policies")
	h.M.Client = h.Config.OAuth2Client(cmd)

	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	pp, err := h.M.Get(args[0])
	pkg.Must(err, "Could not get policy: %s", err)

	p := pp.(*ladon.DefaultPolicy)
	p.Resources = append(p.Resources, args[1:]...)

	err = h.M.Update(p)
	if h.M.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not update policy: %s", err)
	fmt.Printf("Added resources to policy %s", p.ID)
}

func (h *PolicyHandler) RemoveResourceFromPolicy(cmd *cobra.Command, args []string) {
	h.M.Dry, _ = cmd.Flags().GetBool("dry")
	h.M.Endpoint = h.Config.Resolve("/policies")
	h.M.Client = h.Config.OAuth2Client(cmd)

	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	pp, err := h.M.Get(args[0])
	pkg.Must(err, "Could not get policy: %s", err)

	p := pp.(*ladon.DefaultPolicy)
	resources := []string{}
	for _, r := range p.Resources {
		var filter bool
		for _, a := range args[1:] {
			if r == a {
				filter = true
			}
		}
		if !filter {
			resources = append(resources, r)
		}
	}
	p.Resources = resources

	err = h.M.Update(p)
	if h.M.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not update policy: %s", err)
	fmt.Printf("Removed resources from policy %s", p.ID)
}

func (h *PolicyHandler) AddSubjectToPolicy(cmd *cobra.Command, args []string) {
	h.M.Dry, _ = cmd.Flags().GetBool("dry")
	h.M.Endpoint = h.Config.Resolve("/policies")
	h.M.Client = h.Config.OAuth2Client(cmd)

	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	pp, err := h.M.Get(args[0])
	if h.M.Dry {
		fmt.Printf("%s\n", err)
	} else {
		pkg.Must(err, "Could not get policy: %s", err)
	}

	p := pp.(*ladon.DefaultPolicy)
	p.Subjects = append(p.Subjects, args[1:]...)

	err = h.M.Update(p)
	if h.M.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not update policy: %s", err)
	fmt.Printf("Added subjects to policy %s", p.ID)
}

func (h *PolicyHandler) RemoveSubjectFromPolicy(cmd *cobra.Command, args []string) {
	h.M.Dry, _ = cmd.Flags().GetBool("dry")
	h.M.Endpoint = h.Config.Resolve("/policies")
	h.M.Client = h.Config.OAuth2Client(cmd)

	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	pp, err := h.M.Get(args[0])
	pkg.Must(err, "Could not get policy: %s", err)

	p := pp.(*ladon.DefaultPolicy)
	subjects := []string{}
	for _, r := range p.Subjects {
		var filter bool
		for _, a := range args[1:] {
			if r == a {
				filter = true
			}
		}
		if !filter {
			subjects = append(subjects, r)
		}
	}
	p.Subjects = subjects

	err = h.M.Update(p)
	if h.M.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not update policy: %s", err)
	fmt.Printf("Removed subjects from policy %s", p.ID)
}

func (h *PolicyHandler) AddActionToPolicy(cmd *cobra.Command, args []string) {
	h.M.Dry, _ = cmd.Flags().GetBool("dry")
	h.M.Endpoint = h.Config.Resolve("/policies")
	h.M.Client = h.Config.OAuth2Client(cmd)

	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	pp, err := h.M.Get(args[0])
	pkg.Must(err, "Could not get policy: %s", err)

	p := pp.(*ladon.DefaultPolicy)
	p.Actions = append(p.Actions, args[1:]...)

	err = h.M.Update(p)
	if h.M.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not update policy: %s", err)
	fmt.Printf("Added actions to policy %s", p.ID)
}

func (h *PolicyHandler) RemoveActionFromPolicy(cmd *cobra.Command, args []string) {
	h.M.Dry, _ = cmd.Flags().GetBool("dry")
	h.M.Endpoint = h.Config.Resolve("/policies")
	h.M.Client = h.Config.OAuth2Client(cmd)

	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	pp, err := h.M.Get(args[0])
	pkg.Must(err, "Could not get policy: %s", err)

	p := pp.(*ladon.DefaultPolicy)
	actions := []string{}
	for _, r := range p.Actions {
		var filter bool
		for _, a := range args[1:] {
			if r == a {
				filter = true
			}
		}
		if !filter {
			actions = append(actions, r)
		}
	}
	p.Actions = actions

	err = h.M.Update(p)
	if h.M.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not update policy: %s", err)
	fmt.Printf("Removed actions from policy %s", p.ID)
}

func (h *PolicyHandler) GetPolicy(cmd *cobra.Command, args []string) {
	h.M.Dry, _ = cmd.Flags().GetBool("dry")
	h.M.Endpoint = h.Config.Resolve("/policies")
	h.M.Client = h.Config.OAuth2Client(cmd)

	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	p, err := h.M.Get(args[0])
	if h.M.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not retrieve policy: %s", err)

	out, err := json.MarshalIndent(p, "", "\t")
	pkg.Must(err, "Could not convert policy to JSON: %s", err)

	fmt.Printf("%s\n", out)
}

func (h *PolicyHandler) DeletePolicy(cmd *cobra.Command, args []string) {
	h.M.Dry, _ = cmd.Flags().GetBool("dry")
	h.M.Endpoint = h.Config.Resolve("/policies")
	h.M.Client = h.Config.OAuth2Client(cmd)

	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	for _, arg := range args {
		err := h.M.Delete(arg)
		if h.M.Dry {
			fmt.Printf("%s\n", err)
			continue
		}
		pkg.Must(err, "Could not delete policy: %s", err)
		fmt.Printf("Connection %s deleted.\n", arg)
	}
}
