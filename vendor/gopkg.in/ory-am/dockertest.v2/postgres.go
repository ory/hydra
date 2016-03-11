package dockertest

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// SetupPostgreSQLContainer sets up a real PostgreSQL instance for testing purposes,
// using a Docker container. It returns the container ID and its IP address,
// or makes the test fail on error.
func SetupPostgreSQLContainer() (c ContainerID, ip string, port int, err error) {
	port = RandomPort()
	forward := fmt.Sprintf("%d:%d", port, 5432)
	if BindDockerToLocalhost != "" {
		forward = "127.0.0.1:" + forward
	}
	c, ip, err = SetupContainer(PostgresImageName, port, 15*time.Second, func() (string, error) {
		return run("--name", GenerateContainerID(), "-d", "-p", forward, "-e", fmt.Sprintf("POSTGRES_PASSWORD=%s", PostgresPassword), PostgresImageName)
	})
	return
}

// ConnectToPostgreSQL starts a PostgreSQL image and passes the database url to the connector callback.
func ConnectToPostgreSQL(tries int, delay time.Duration, connector func(url string) bool) (c ContainerID, err error) {
	c, ip, port, err := SetupPostgreSQLContainer()
	if err != nil {
		return c, fmt.Errorf("Could not set up PostgreSQL container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		url := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=disable", PostgresUsername, PostgresPassword, ip, port)
		if connector(url) {
			return c, nil
		}
		log.Printf("Try %d failed. Retrying.", try)
	}
	return c, errors.New("Could not set up PostgreSQL container.")
}
