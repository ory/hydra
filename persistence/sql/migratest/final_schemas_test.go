package migratest

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/viper"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/sqlcon/dockertest"
)

var comments = regexp.MustCompile("(--[^\n]*\n)|(?s:/\\*.+\\*/)")
var migrationTableStatements = regexp.MustCompile("[^;]*(hydra_[a-zA-Z0-9_]*_migration|schema_migration)[^;]*;")

func stripDump(d string) string {
	d = comments.ReplaceAllLiteralString(d, "")
	d = migrationTableStatements.ReplaceAllLiteralString(d, "")
	return strings.ReplaceAll(d, "\r\n", "")
}

func getContainerID(t *testing.T, containerPort string) string {
	cid, err := exec.Command("docker", "ps", "-f", fmt.Sprintf("expose=%s", containerPort), "-q").CombinedOutput()
	require.NoError(t, err)
	containerID := strings.TrimSuffix(string(cid), "\n")
	require.False(t, strings.Contains(containerID, "\n"), "there is more than one docker container running with port %s, I am confused: %s", containerPort, containerID)
	return containerID
}

func dump(t *testing.T) string {
	dump, err := exec.Command("docker", "exec", "-t", getContainerID(t, "26257"), "./cockroach", "dump", "defaultdb", "--insecure", "--dump-mode=schema").CombinedOutput()
	require.NoError(t, err, "%s", dump)
	return stripDump(string(dump))
}

func TestGetFinalSchemas(t *testing.T) {
	c := dockertest.ConnectToTestCockroachDBPop(t)

	viper.Set(configuration.ViperKeyDSN, "cockroach"+strings.TrimLeft(c.URL(), "postgres"))
	d := driver.NewDefaultDriver(logrusx.New("", ""), true, []string{}, "", "", "", false)

	require.NoError(t, d.Registry().Persister().MigrateUp(context.Background()))

	fmt.Println(dump(t))
}
