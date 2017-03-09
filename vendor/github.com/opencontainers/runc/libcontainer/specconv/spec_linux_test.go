// +build linux

package specconv

import (
	"testing"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func TestLinuxCgroupsPathSpecified(t *testing.T) {
	cgroupsPath := "/user/cgroups/path/id"

	spec := &specs.Spec{}
	spec.Linux = &specs.Linux{
		CgroupsPath: &cgroupsPath,
	}

	cgroup, err := createCgroupConfig("ContainerID", false, spec)
	if err != nil {
		t.Errorf("Couldn't create Cgroup config: %v", err)
	}

	if cgroup.Path != cgroupsPath {
		t.Errorf("Wrong cgroupsPath, expected '%s' got '%s'", cgroupsPath, cgroup.Path)
	}
}

func TestLinuxCgroupsPathNotSpecified(t *testing.T) {
	spec := &specs.Spec{}

	cgroup, err := createCgroupConfig("ContainerID", false, spec)
	if err != nil {
		t.Errorf("Couldn't create Cgroup config: %v", err)
	}

	if cgroup.Path != "" {
		t.Errorf("Wrong cgroupsPath, expected it to be empty string, got '%s'", cgroup.Path)
	}
}

func TestDupNamespaces(t *testing.T) {
	spec := &specs.Spec{
		Linux: &specs.Linux{
			Namespaces: []specs.LinuxNamespace{
				{
					Type: "pid",
				},
				{
					Type: "pid",
					Path: "/proc/1/ns/pid",
				},
			},
		},
	}

	_, err := CreateLibcontainerConfig(&CreateOpts{
		Spec: spec,
	})

	if err == nil {
		t.Errorf("Duplicated namespaces should be forbidden")
	}
}
