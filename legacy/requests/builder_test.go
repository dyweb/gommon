package requests

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestTransportBuilder_UseSocks5(t *testing.T) {
	assert := asst.New(t)
	b := TransportBuilder.UseSocks5("127.0.0.1:1080", "", "")
	assert.NotNil(TransportBuilder.auth)
	assert.Nil(b.auth)
}

func TestTransportBuilder_Build(t *testing.T) {
	assert := asst.New(t)
	b := TransportBuilder.UseSocks5("127.0.0.1:1080", "", "")
	// NOTE: it does not connect to the proxy, so the test should pass regardless of the running proxy
	tr, err := b.Build()
	assert.Nil(err)
	assert.NotNil(tr)
}
