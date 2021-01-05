package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	"camlistore.org/pkg/netutil"
	"github.com/spf13/viper"
)

const (
	mongoImage    = "mongo:4.2@sha256:1c2243a5e21884ffa532ca9d20c221b170d7b40774c235619f98e2f6eaec520a"
	mongoURI      = "mongodb://localhost:27017" // TODO: switch to 27016
	mongoDatabase = "signauxfaibles"
)

func TestMain(t *testing.T) {
	// SetupMongoContainer may skip or fatal the test if docker isn't found or something goes
	// wrong when setting up the container. Thus, no error is returned
	_ /*containerID*/, ip := setupMongoContainer(t)
	fmt.Println(ip)
	defer KillRemove(t)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		viper.Set("DB_DIAL", mongoURI)
		viper.Set("DB", mongoDatabase)
		os.Args = []string{"sfdata", "etablissements"}
		main()
		wg.Done()
	}()
	wg.Wait()
}

// Code from https://github.com/niilo/golib/blob/master/test/dockertest/docker.go

type ContainerID string

func (c ContainerID) IP() (string, error) {
	return IP(string(c))
}

func (c ContainerID) Kill() error {
	return KillContainer(string(c))
}

// Remove runs "docker rm" on the container
func (c ContainerID) Remove() error {
	return exec.Command("docker", "rm", string(c)).Run()
}

// KillRemove calls Kill on the container, and then Remove if there was
// no error. It logs any error to t.
func /*(c ContainerID)*/ KillRemove(t *testing.T) {
	if err := exec.Command("docker", "stop", "dockertestmongodb").Run(); err != nil {
		t.Log(err)
		return
	}
	// if err := c.Kill(); err != nil {
	// 	t.Log(err)
	// 	return
	// }
	// if err := c.Remove(); err != nil {
	// 	t.Log(err)
	// }
}

// lookup retrieves the ip address of the container, and tries to reach
// before timeout the tcp address at this ip and given port.
func (c ContainerID) lookup(port int, timeout time.Duration) (ip string, err error) {
	ip, err = c.IP()
	if err != nil {
		err = fmt.Errorf("error getting IP: %v", err)
		return
	}
	addr := fmt.Sprintf("%s:%d", ip, port)
	err = netutil.AwaitReachable(addr, timeout)
	return
}

func KillContainer(container string) error {
	return exec.Command("docker", "kill", container).Run()
}

// IP returns the IP address of the container.
func IP(containerID string) (string, error) {
	out, err := exec.Command("docker", "inspect", containerID).Output()
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

// setupMongoContainer sets up a real MongoDB instance for testing purposes,
// using a Docker container. It returns the container ID and its IP address,
// or makes the test fail on error.
// Currently using https://hub.docker.com/_/mongo/
func setupMongoContainer(t *testing.T) (c ContainerID, ip string) {
	return setupContainer(t, mongoImage, 27017, 10*time.Second, func() (string, error) {
		log.Println("Starting mongodb container...")
		return run("--rm", "-d", "-p", "27017:27017", "--name", "dockertestmongodb", mongoImage)
	})
}

// setupContainer sets up a container, using the start function to run the given image.
// It also looks up the IP address of the container, and tests this address with the given
// port and timeout. It returns the container ID and its IP address, or makes the test
// fail on error.
func setupContainer(t *testing.T, image string, port int, timeout time.Duration, start func() (string, error)) (c ContainerID, ip string) {
	runLongTest(t, image)

	/*containerID*/
	_, err := start()
	if err != nil {
		t.Fatalf("docker run: %v", err)
	}
	// c = ContainerID(containerID)
	// ip, err = c.lookup(port, timeout)
	// if err != nil {
	// 	c.KillRemove(t)
	// 	t.Errorf("Container %v setup failed: %v", c, err)
	// }
	return
}

// haveDocker returns whether the "docker" command was found.
func haveDocker() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

func haveImage(name string) (ok bool, err error) {
	out, err := exec.Command("docker", "images", "--no-trunc").Output()
	if err != nil {
		return
	}
	return bytes.Contains(out, []byte(name)), nil
}

// Pull retrieves the docker image with 'docker pull'.
func Pull(image string) error {
	out, err := exec.Command("docker", "pull", image).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%v: %s", err, out)
	}
	return err
}

/// runLongTest checks all the conditions for running a docker container
// based on image.
func runLongTest(t *testing.T, image string) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	if !haveDocker() {
		t.Error("'docker' command not found")
	}
	if ok, err := haveImage(image); !ok || err != nil {
		if err != nil {
			t.Errorf("Error running docker to check for %s: %v", image, err)
		}
		log.Printf("Pulling docker image %s ...", image)
		if err := Pull(image); err != nil {
			t.Errorf("Error pulling %s: %v", image, err)
		}
	}
}

func run(args ...string) (containerID string, err error) {
	log.Println(append([]string{"docker", "run"}, args...))
	cmd := exec.Command("docker", append([]string{"run"}, args...)...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	log.Println(stdout.String())
	log.Println(stderr.String())
	if err = cmd.Run(); err != nil {
		return "", err // fmt.Errorf("%v%v", stderr.String(), err)
	}
	containerID = "" // strings.TrimSpace(stdout.String())
	// if containerID == "" {
	// 	return "", errors.New("unexpected empty output from `docker run`")
	// }
	return
}
