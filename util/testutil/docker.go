package testutil

import (
	"errors"
	"fmt"
	"os/exec"
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
}

type Container struct {
	cfg ContainerConfig
	cmd *exec.Cmd
}

// NewContainer calls docker run to start a container.
func NewContainer(cfg ContainerConfig) (*Container, error) {
	c := &Container{cfg: cfg}
	return c, c.run()
}

// ToDockerArgs validates and converts a container config to docker run arguments.
func (cfg *ContainerConfig) ToDockerArgs() ([]string, error) {
	args := []string{"run"}
	for _, port := range cfg.Ports {
		// https://docs.docker.com/config/containers/container-networking/#published-ports
		// -p 8080:80 Map TCP port 80 in the container to port 8080 on the Docker host
		args = append(args, "-p", fmt.Sprintf("%d:%d", port.Host, port.Container))
	}
	if cfg.Image == "" {
		return nil, errors.New("image is empty")
	}
	args = append(args, cfg.Image)
	return args, nil
}

func (c *Container) run() error {
	args, err := c.cfg.ToDockerArgs()
	if err != nil {
		return err
	}
	cmd := exec.Command("docker", args...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error start docker command %w", err)
	}
	c.cmd = cmd
	return nil
}

func (c *Container) Stop() error {
	// TODO: kill the process only works if we run container in foreground
	// FIXME: it seems the container is still running after sending kill, might need to call docker rm
	// might check https://github.com/docker/cli/blob/master/cli/command/container/run.go
	if err := c.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("error kill docker run %w", err)
	}
	return nil
}
