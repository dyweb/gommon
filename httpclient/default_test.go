package httpclient_test

import (
	"testing"

	"github.com/dyweb/gommon/httpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefault(t *testing.T) {
	d, err := httpclient.NewDefault("http://localhost:8080")
	require.Nil(t, err)
	// can't have direct access to transport
	tr, ok := d.Transport()
	assert.False(t, ok)
	assert.Nil(t, tr)
}
