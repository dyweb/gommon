// +build ignore

package requests

import (
	asst "github.com/stretchr/testify/assert"

	"testing"
)

func TestExtractBasicAuth(t *testing.T) {
	assert := asst.New(t)
	username := "Aladdin"
	password := "OpenSesame"
	// get from Postman
	header := "Basic QWxhZGRpbjpPcGVuU2VzYW1l"
	assert.Equal(header, GenerateBasicAuth(username, password))
	u, p, err := ExtractBasicAuth(header)
	assert.Nil(err)
	assert.Equal(username, u)
	assert.Equal(password, p)
}
