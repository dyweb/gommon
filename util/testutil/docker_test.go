package testutil

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dyweb/gommon/util/netutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContainer_DockerRunArgs(t *testing.T) {
	tStart := time.Unix(1583459680, 0)
	SetTestStart(tStart)
	tCreate := time.Unix(1583459690, 0)
	cfg := ContainerConfig{
		Image: "nginx",
		Ports: []PortMapping{
			{
				Host:      8093,
				Container: 80,
			},
		},
		now: tCreate,
	}
	c, err := NewContainerWithoutRun(cfg)
	require.Nil(t, err)
	args := c.DockerRunArgs()
	assert.Equal(t, "run -p 8093:80 -l gommon-container=1 -l gommon-test-start-time=1583459680000000000 -l gommon-create-time=1583459690000000000 nginx", strings.Join(args, " "))
}

func TestContainer_Stop(t *testing.T) {
	t.Skip("always skip docker test")

	// TODO: need more rules, the test always run for machine that has active docker running
	RunIf(t, HasDocker())

	port, err := netutil.AvailablePortBySystem()
	assert.Nil(t, err)
	cfg := ContainerConfig{
		Image: "nginx",
		Ports: []PortMapping{
			{
				Host:      port,
				Container: 80,
			},
		},
	}
	c, err := NewContainer(cfg)
	assert.Nil(t, err)
	// TODO: wait until container is ready, this needs to be provided in config
	time.Sleep(1 * time.Second)
	b := GetBody(t, nil, fmt.Sprintf("http://localhost:%d", port))
	assert.Contains(t, string(b), "Welcome to nginx")
	err = c.Stop()
	assert.Nil(t, err)
}
