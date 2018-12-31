package cli_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/log/handlers/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: test fields and caller

func TestNew(t *testing.T) {
	buf := bytes.Buffer{}
	h := cli.New(&buf, true)
	h.HandleLog(log.InfoLevel, time.Now(), "hi", log.EmptyCaller(), nil, nil)
	assert.Equal(t, "\x1b[34mINFO\x1b[0m 0000 hi\n", buf.String())
}

func TestNewNoColor(t *testing.T) {
	buf := bytes.Buffer{}
	h := cli.NewNoColor(&buf)
	tm, err := time.Parse(cli.DefaultTimeStampFormat, "2018-12-30T21:10:49-08:00")
	require.Nil(t, err)
	h.HandleLog(log.InfoLevel, tm, "hi", log.EmptyCaller(), nil, nil)
	assert.Equal(t, "INFO 2018-12-30T21:10:49-08:00 hi\n", buf.String())
}
