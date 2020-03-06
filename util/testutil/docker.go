package testutil

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

// docker.go starts and stop container for testing by shell out to docker cli
// https://github.com/dyweb/gommon/issues/124

// PortMapping maps a container port to host
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#containerport-v1-core
type PortMapping struct {
	Container int
	Host      int
}

// ContainerConfig specifies how to create the container
type ContainerConfig struct {
	// Image contains image name with tag, it is what you put after docker run, e.g. dyweb/go-dev:1.14
	Image string
	Ports []PortMapping

	now time.Time
}

type Container struct {
	cfg        ContainerConfig
	createTime time.Time

	// TODO: do we really need to save the docker run process? we are not doing anything with it after switch to docker rm
	cmd    *exec.Cmd
	id     string
	labels []string
}

// NewContainer shell out to docker cli and runs a container in foreground.
func NewContainer(cfg ContainerConfig) (*Container, error) {
	now := time.Now()
	if !cfg.now.IsZero() {
		now = cfg.now
	}
	c := &Container{
		cfg:        cfg,
		createTime: now,
	}
	c.labels = []string{"gommon-container=1", fmt.Sprintf("gommon-test-start-time=%d", testStart.UnixNano()),
		fmt.Sprintf("gommon-create-time=%d", now.UnixNano())}
	return c, c.run()
}

// NewContainerWithoutRun validate config, generate labels and does not execute any docker command.
func NewContainerWithoutRun(cfg ContainerConfig) (*Container, error) {
	if cfg.Image == "" {
		return nil, errors.New("image is empty")
	}

	now := time.Now()
	if !cfg.now.IsZero() {
		now = cfg.now
	}
	c := &Container{
		cfg:        cfg,
		createTime: now,
	}
	c.labels = []string{"gommon-container=1", fmt.Sprintf("gommon-test-start-time=%d", testStart.UnixNano()),
		fmt.Sprintf("gommon-create-time=%d", now.UnixNano())}
	return c, nil
}

// DockerRunArgs converts a container config to docker run arguments.
func (c *Container) DockerRunArgs() []string {
	args := []string{"run"}
	for _, port := range c.cfg.Ports {
		// https://docs.docker.com/config/containers/container-networking/#published-ports
		// -p 8080:80 Map TCP port 80 in the container to port 8080 on the Docker host
		args = append(args, "-p", fmt.Sprintf("%d:%d", port.Host, port.Container))
	}
	for _, l := range c.labels {
		args = append(args, "-l", l)
	}
	args = append(args, c.cfg.Image)
	return args
}

// DockerPsArgs is used to filter out container by all of its labels
func (c *Container) DockerPsArgs() []string {
	args := []string{"ps"}
	for _, l := range c.labels {
		args = append(args, "--filter", "label="+l)
	}
	args = append(args, "-q")
	return args
}

// run calls docker run in foreground and collect its id.
func (c *Container) run() error {
	pullCmd := exec.Command("docker", "pull", c.cfg.Image)
	out, err := pullCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker pull failed: %w %s", err, string(out))
	}

	// docker run
	runCmd := exec.Command("docker", c.DockerRunArgs()...)
	if err := runCmd.Start(); err != nil {
		return fmt.Errorf("error start docker command %w", err)
	}
	c.cmd = runCmd

	// TODO(#126): manual retry until we have a retry package
	// The drawback of shell out is we don't know when the container will be ready especially when pulling is needed.
	for i := 0; i < 5; i++ {
		// docker ps to get id
		psCmd := exec.Command("docker", c.DockerPsArgs()...)
		out, err = psCmd.CombinedOutput()
		id := string(bytes.TrimSpace(out))
		if err != nil || id == "" {
			time.Sleep(1 * time.Duration(i) * time.Second)
			continue
		}
		c.id = id
		return nil
	}
	return fmt.Errorf("error get conatienr id : %w %s", err, string(out))
}

// Stop removes the (running) container by id
func (c *Container) Stop() error {
	// TODO: kill the foreground docker run process? though it will exit once the container is removed

	delCmd := exec.Command("docker", "rm", "-f", c.id)
	out, err := delCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error remove container %s: %w %s", c.id, err, string(out))
	}
	return nil
}
