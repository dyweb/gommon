package testutil

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/dyweb/gommon/util/netutil"
	"github.com/stretchr/testify/assert"
)

func TestContainerConfig_ToDockerArgs(t *testing.T) {
	cfg := ContainerConfig{
		Image: "nginx",
		Ports: []PortMapping{
			{
				Host:      8093,
				Container: 80,
			},
		},
	}
	args, err := cfg.ToDockerArgs()
	assert.Nil(t, err)
	assert.Equal(t, "run -p 8093:80 nginx", strings.Join(args, " "))
}

func TestContainer_Stop(t *testing.T) {
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
	// TODO: wait until container is ready, this needs to be provided by the client ...
	time.Sleep(2 * time.Second)
	res, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
	assert.Nil(t, err)
	b, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Contains(t, string(b), "Welcome to nginx")
	err = c.Stop()
	assert.Nil(t, err)
}
