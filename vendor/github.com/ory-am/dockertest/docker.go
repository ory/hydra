package dockertest

/*
Copyright 2014 The Camlistore Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"database/sql"
	"math/rand"
	"regexp"

	"github.com/mattbaird/elastigo/lib"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2"

	// Import mysql driver
	_ "github.com/go-sql-driver/mysql"

	// Import postgres driver
	_ "github.com/lib/pq"
)

/// runLongTest checks all the conditions for running a docker container
// based on image.
func runLongTest(image string) error {
	DockerMachineAvailable = false
	if haveDockerMachine() {
		DockerMachineAvailable = true
		if !startDockerMachine() {
			log.Printf(`Starting docker machine "%s" failed. This could be because the image is already running or because the image does not exist. Tests will fail if the image does not exist.`, DockerMachineName)
		}
	} else if !haveDocker() {
		return errors.New("Neither 'docker' nor 'docker-machine' available on this system.")
	}
	if ok, err := haveImage(image); !ok || err != nil {
		if err != nil {
			return fmt.Errorf("Error checking for docker image %s: %v", image, err)
		}
		log.Printf("Pulling docker image %s ...", image)
		if err := Pull(image); err != nil {
			return fmt.Errorf("Error pulling %s: %v", image, err)
		}
	}
	return nil
}

func runDockerCommand(command string, args ...string) *exec.Cmd {
	if DockerMachineAvailable {
		command = "/usr/local/bin/" + strings.Join(append([]string{command}, args...), " ")
		cmd := exec.Command("docker-machine", "ssh", DockerMachineName, command)
		return cmd
	}
	return exec.Command(command, args...)
}

// haveDockerMachine returns whether the "docker" command was found.
func haveDockerMachine() bool {
	_, err := exec.LookPath("docker-machine")
	return err == nil
}

// startDockerMachine starts the docker machine and returns false if the command failed to execute
func startDockerMachine() bool {
	_, err := exec.Command("docker-machine", "start", DockerMachineName).Output()
	return err == nil
}

// haveDocker returns whether the "docker" command was found.
func haveDocker() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

func haveImage(name string) (ok bool, err error) {
	out, err := runDockerCommand("docker", "images", "--no-trunc").Output()
	if err != nil {
		return false, err
	}
	return bytes.Contains(out, []byte(name)), nil
}

func run(args ...string) (containerID string, err error) {
	var stdout, stderr bytes.Buffer
	validID := regexp.MustCompile(`^([a-zA-Z0-9]+)$`)
	cmd := runDockerCommand("docker", append([]string{"run"}, args...)...)

	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("Error running docker\nStdOut: %s\nStdErr: %s\nError: %v\n\n", stdout.String(), stderr.String(), err)
		return
	}
	containerID = strings.TrimSpace(string(stdout.String()))
	if !validID.MatchString(containerID) {
		return "", fmt.Errorf("Error running docker: %s", containerID)
	}
	if containerID == "" {
		return "", errors.New("Unexpected empty output from `docker run`")
	}
	return containerID, nil
}

// KillContainer runs docker kill on a container.
func KillContainer(container string) error {
	if container != "" {
		return runDockerCommand("docker", "kill", container).Run()
	}
	return nil
}

// Pull retrieves the docker image with 'docker pull'.
func Pull(image string) error {
	out, err := runDockerCommand("docker", "pull", image).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%v: %s", err, out)
	}
	return err
}

// IP returns the IP address of the container.
func IP(containerID string) (string, error) {
	out, err := runDockerCommand("docker", "inspect", containerID).Output()
	if err != nil {
		return "", err
	}
	type networkSettings struct {
		IPAddress string
	}
	type container struct {
		NetworkSettings networkSettings
	}
	var c []container
	if err := json.NewDecoder(bytes.NewReader(out)).Decode(&c); err != nil {
		return "", err
	}
	if len(c) == 0 {
		return "", errors.New("no output from docker inspect")
	}
	if ip := c[0].NetworkSettings.IPAddress; ip != "" {
		return ip, nil
	}
	return "", errors.New("could not find an IP. Not running?")
}

// setupContainer sets up a container, using the start function to run the given image.
// It also looks up the IP address of the container, and tests this address with the given
// port and timeout. It returns the container ID and its IP address, or makes the test
// fail on error.
func setupContainer(image string, port int, timeout time.Duration, start func() (string, error)) (c ContainerID, ip string, err error) {
	err = runLongTest(image)
	if err != nil {
		return "", "", err
	}

	containerID, err := start()
	if err != nil {
		return "", "", err
	}

	c = ContainerID(containerID)
	ip, err = c.lookup(port, timeout)
	if err != nil {
		c.KillRemove()
		return "", "", err
	}
	return c, ip, nil
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

type pinger interface {
	Ping() error
}

func ping(db pinger, tries int, delay time.Duration) bool {
	for i := 0; i <= tries; i++ {
		time.Sleep(delay)
		if s, ok := db.(*sql.DB); ok {
			if _, err := s.Exec("SELECT 1"); err != nil {
				continue
			}
		} else if s, ok := db.(*mgo.Session); ok {
			if _, err := s.DatabaseNames(); err != nil {
				continue
			}
		}
		if err := db.Ping(); err == nil {
			return true
		}
		log.Printf("Ping try %v failed", i)
	}
	return false
}

// SetupMongoContainer sets up a real MongoDB instance for testing purposes,
// using a Docker container. It returns the container ID and its IP address,
// or makes the test fail on error.
func SetupMongoContainer() (c ContainerID, ip string, port int, err error) {
	port = randInt(1024, 49150)
	forward := fmt.Sprintf("%d:%d", port, 27017)
	if BindDockerToLocalhost != "" {
		forward = "127.0.0.1:" + forward
	}
	c, ip, err = setupContainer(mongoImage, port, 10*time.Second, func() (string, error) {
		res, err := run("--name", uuid.New(), "-d", "-P", "-p", forward, mongoImage)
		return res, err
	})
	return
}

// SetupMySQLContainer sets up a real MySQL instance for testing purposes,
// using a Docker container. It returns the container ID and its IP address,
// or makes the test fail on error.
func SetupMySQLContainer() (c ContainerID, ip string, port int, err error) {
	port = randInt(1024, 49150)
	forward := fmt.Sprintf("%d:%d", port, 3306)
	if BindDockerToLocalhost != "" {
		forward = "127.0.0.1:" + forward
	}
	c, ip, err = setupContainer(mysqlImage, port, 10*time.Second, func() (string, error) {
		return run("-d", "-p", forward, "-e", fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", MySQLPassword), mysqlImage)
	})
	return
}

// SetupPostgreSQLContainer sets up a real PostgreSQL instance for testing purposes,
// using a Docker container. It returns the container ID and its IP address,
// or makes the test fail on error.
func SetupPostgreSQLContainer() (c ContainerID, ip string, port int, err error) {
	port = randInt(1024, 49150)
	forward := fmt.Sprintf("%d:%d", port, 5432)
	if BindDockerToLocalhost != "" {
		forward = "127.0.0.1:" + forward
	}
	c, ip, err = setupContainer(postgresImage, port, 15*time.Second, func() (string, error) {
		return run("-d", "-p", forward, "-e", fmt.Sprintf("POSTGRES_PASSWORD=%s", PostgresPassword), postgresImage)
	})
	return
}

// SetupElasticSearchContainer sets up a real ElasticSearch instance for testing purposes
// using a Docker container. It returns the container ID and its IP address,
// or makes the test fail on error.
func SetupElasticSearchContainer() (c ContainerID, ip string, port int, err error) {
	port = randInt(1024, 49150)
	forward := fmt.Sprintf("%d:%d", port, 9200)
	if BindDockerToLocalhost != "" {
		forward = "127.0.0.1:" + forward
	}
	c, ip, err = setupContainer(elasticsearchImage, port, 15*time.Second, func() (string, error) {
		return run("--name", uuid.New(), "-d", "-P", "-p", forward, elasticsearchImage)
	})
	return
}

// SetupRedisContainer sets up a real Redis instance for testing purposes
// using a Docker container. It returns the container ID and its IP address,
// or makes the test fail on error.
func SetupRedisContainer() (c ContainerID, ip string, port int, err error) {
	port = randInt(1024, 49150)
	forward := fmt.Sprintf("%d:%d", port, 6379)
	if BindDockerToLocalhost != "" {
		forward = "127.0.0.1:" + forward
	}
	c, ip, err = setupContainer(redisImage, port, 15*time.Second, func() (string, error) {
		return run("--name", uuid.New(), "-d", "-P", "-p", forward, redisImage)
	})
	return
}

// OpenPostgreSQLContainerConnection is supported for legacy reasons. Don't use it.
func OpenPostgreSQLContainerConnection(tries int, delay time.Duration) (c ContainerID, db *sql.DB, err error) {
	c, ip, port, err := SetupPostgreSQLContainer()
	if err != nil {
		return c, nil, fmt.Errorf("Could not set up PostgreSQL container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		url := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=disable", PostgresUsername, PostgresPassword, ip, port)
		log.Printf("Try %d: Connecting %s", try, url)
		if db, err = sql.Open("postgres", url); err == nil {
			if ping(db, tries, delay) {
				log.Printf("Try %d: Successfully connected to %v", try, url)
				return c, db, nil
			}
			log.Printf("Try %d: Could not ping database: %v", try, err)
		}
		log.Printf("Try %d: Could not set up PostgreSQL container: %v", try, err)
	}
	return c, nil, errors.New("Could not set up PostgreSQL container.")
}

// OpenMongoDBContainerConnection is supported for legacy reasons. Don't use it.
func OpenMongoDBContainerConnection(tries int, delay time.Duration) (c ContainerID, db *mgo.Session, err error) {
	c, ip, port, err := SetupMongoContainer()
	if err != nil {
		return c, nil, fmt.Errorf("Could not set up MongoDB container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		url := fmt.Sprintf("%s:%d", ip, port)
		log.Printf("Try %d: Connecting %s", try, url)
		if db, err = mgo.Dial(url); err == nil {
			if ping(db, tries, delay) {
				log.Printf("Try %d: Successfully connected to %v", try, url)
				return c, db, nil
			}
			log.Printf("Try %d: Could not ping database: %v", try, err)
		}
		log.Printf("Try %d: Could not set up MongoDB container: %v", try, err)
	}
	return c, nil, errors.New("Could not set up MongoDB container.")
}

// OpenMySQLContainerConnection is supported for legacy reasons. Don't use it.
func OpenMySQLContainerConnection(tries int, delay time.Duration) (c ContainerID, db *sql.DB, err error) {
	c, ip, port, err := SetupMySQLContainer()
	if err != nil {
		return c, nil, fmt.Errorf("Could not set up MySQL container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		url := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql", MySQLUsername, MySQLPassword, ip, port)
		log.Printf("Try %d: Connecting %s", try, url)
		if db, err = sql.Open("mysql", url); err == nil {
			if ping(db, tries, delay) {
				log.Printf("Try %d: Successfully connected to %v", try, url)
				return c, db, nil
			}
			log.Printf("Try %d: Could not ping database: %v", try, err)
		}
		log.Printf("Try %d: Could not set up MySQL container: %v", try, err)
	}
	return c, nil, errors.New("Could not set up MySQL container.")
}

// OpenElasticSearchContainerConnection is supported for legacy reasons. Don't use it.
func OpenElasticSearchContainerConnection(tries int, delay time.Duration) (c ContainerID, con *elastigo.Conn, err error) {
	c, ip, port, err := SetupElasticSearchContainer()
	if err != nil {
		return c, nil, fmt.Errorf("Could not set up ElasticSearch container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		url := fmt.Sprintf("%s:%d", ip, port)
		log.Printf("Try %d: Connecting %s", try, url)

		conn := elastigo.NewConn()
		conn.Domain = ip
		conn.Port = strconv.Itoa(port)

		resp, err := conn.Health("")
		if err == nil && resp.Status != "" {
			log.Printf("Try %d: Successfully connected to %v", try, conn.Domain)
			return c, conn, nil
		}

		log.Printf("Try %d: Could not set up ElasticSearch container: %v", try, err)
	}
	return c, nil, errors.New("Could not set up ElasticSearch container.")
}

// OpenRedisContainerConnection is supported for legacy reasons. Don't use it.
func OpenRedisContainerConnection(tries int, delay time.Duration) (c ContainerID, client *redis.Client, err error) {
	c, ip, port, err := SetupRedisContainer()
	if err != nil {
		return c, nil, fmt.Errorf("Could not set up Redis container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		url := fmt.Sprintf("%s:%d", ip, port)
		log.Printf("Try %d: Connecting %s", try, url)

		client, err := redis.DialTimeout("tcp", url, 10*time.Second)
		if err == nil {
			log.Printf("Try %d: Successfully connected to %v", try, client.Addr)
			return c, client, nil
		}

		log.Printf("Try %d: Could not set up Redis container: %v", try, err)
	}
	return c, nil, errors.New("Could not set up Redis container.")
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

// ConnectToMongoDB starts a MongoDB image and passes the database url to the connector callback.
// The url will match the ip:port pattern (e.g. 123.123.123.123:4241)
func ConnectToMongoDB(tries int, delay time.Duration, connector func(url string) bool) (c ContainerID, err error) {
	c, ip, port, err := SetupMongoContainer()
	if err != nil {
		return c, fmt.Errorf("Could not set up MongoDB container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		url := fmt.Sprintf("%s:%d", ip, port)
		if connector(url) {
			return c, nil
		}
		log.Printf("Try %d failed. Retrying.", try)
	}
	return c, errors.New("Could not set up MongoDB container.")
}

// ConnectToMySQL starts a MySQL image and passes the database url to the connector callback function.
// The url will match the username:password@tcp(ip:port) pattern (e.g. `root:root@tcp(123.123.123.123:3131)`)
func ConnectToMySQL(tries int, delay time.Duration, connector func(url string) bool) (c ContainerID, err error) {
	c, ip, port, err := SetupMySQLContainer()
	if err != nil {
		return c, fmt.Errorf("Could not set up MySQL container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		url := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql", MySQLUsername, MySQLPassword, ip, port)
		if connector(url) {
			return c, nil
		}
		log.Printf("Try %d failed. Retrying.", try)
	}
	return c, errors.New("Could not set up MySQL container.")
}

// ConnectToElasticSearch starts an ElasticSearch image and passes the database url to the connector callback function.
// The url will match the ip:port pattern (e.g. 123.123.123.123:4241)
func ConnectToElasticSearch(tries int, delay time.Duration, connector func(url string) bool) (c ContainerID, err error) {
	c, ip, port, err := SetupElasticSearchContainer()
	if err != nil {
		return c, fmt.Errorf("Could not set up ElasticSearch container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		url := fmt.Sprintf("%s:%d", ip, port)
		if connector(url) {
			return c, nil
		}
		log.Printf("Try %d failed. Retrying.", try)
	}
	return c, errors.New("Could not set up ElasticSearch container.")
}

// ConnectToRedis starts a Redis image and passes the database url to the connector callback function.
// The url will match the ip:port pattern (e.g. 123.123.123.123:6379)
func ConnectToRedis(tries int, delay time.Duration, connector func(url string) bool) (c ContainerID, err error) {
	c, ip, port, err := SetupRedisContainer()
	if err != nil {
		return c, fmt.Errorf("Could not set up Redis container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		url := fmt.Sprintf("%s:%d", ip, port)
		if connector(url) {
			return c, nil
		}
		log.Printf("Try %d failed. Retrying.", try)
	}
	return c, errors.New("Could not set up Redis container.")
}
