package health

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/ory/hydra/config"
)

var (
	serverStartTime time.Time
)

const (
	OK   = "OK"
	WARN = "WARN"
	CRIT = "CRIT"
)

type dependentInfo struct {
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	State        dependentState `json:"state"`
	ResponseTime float64        `json:"responseTime"`
}

type dependentState struct {
	Status   string `json:"status"`
	Details  string `json:"details,omitempty"`
	Version  string `json:"version,omitempty"`
	Revision string `json:"revision,omitempty"`
}

type projectInfo struct {
	Repo   string   `json:"repo"`
	Home   string   `json:"home"`
	Owners []string `json:"owners"`
	Logs   []string `json:"logs"`
	Stats  []string `json:"stats"`
}

func init() {
	serverStartTime = time.Now()
}

func status(status string, c *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"status":   status,
		"version":  c.BuildVersion,
		"revision": c.BuildHash,
	}
}

func simpleStatus(c *config.Config) []byte {
	content := status(OK, c)
	data, _ := json.Marshal(content)
	return data
}

func detailedStatus(c *config.Config) []byte {
	dependent := dbCheck(c.Context().Connection)
	status := status(dependent.State.Status, c)
	status["project"] = getProject()
	status["host"] = os.Getenv("ISSUER")
	status["description"] = "Sand authentication service for service to service communications."
	status["name"] = "Sand"
	status["uptime"] = int64(time.Since(serverStartTime).Seconds())
	status["dependencies"] = []interface{}{dependent}
	data, _ := json.Marshal(status)
	return data
}

func getProject() projectInfo {
	logsStr := os.Getenv("APPLICATION_LOG_LINKS")
	logs := strings.Split(logsStr, " ")

	statsStr := os.Getenv("APPLICATION_STATS_LINKS")
	stats := strings.Split(statsStr, " ")

	return projectInfo{
		Repo:   "https://github.com/coupa/hydra-sand",
		Home:   "https://github.com/coupa/hydra-sand",
		Owners: []string{"Technology Platform"},
		Logs:   logs,
		Stats:  stats,
	}
}

func dbCheck(connection interface{}) dependentInfo {
	var err error
	name := ""
	var t float64
	sTime := time.Now()

	switch conn := connection.(type) {
	case *config.MemoryConnection:
		name = "memory"
	case *config.SQLConnection:
		err = conn.Ping()
		t = time.Since(sTime).Seconds()
		name = conn.URL.Scheme
	case *config.PluginConnection:
		name = "plugin"
		err = conn.Connect()
		t = time.Since(sTime).Seconds()
	default:
		err = errors.New("No DB connection")
	}

	state := dependentState{Status: OK}
	if err != nil {
		state.Status = CRIT
		state.Details = err.Error()
	}
	info := dependentInfo{
		Name:         "Database (" + name + ")",
		Type:         "internal",
		State:        state,
		ResponseTime: t,
	}
	return info
}
