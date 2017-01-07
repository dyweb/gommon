package config

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"log"
)

func TestYAMLConfig_Parse(t *testing.T) {
	assert := assert.New(t)
	var dat = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`
	c := NewYAMLConfig()
	err := c.Parse([]byte(dat))
	assert.Nil(err)

	// NOTE： this is invalid yaml because when you use ` syntax to declare long string in Golang,
	// the indent are also included in the string, so this yaml has indent without any parent, which is invalid
	var invalidDat = `
	a: Easy!
	b:
	  c: 2
	  d: [3, 4]
	`
	// the print should show you the string has indent you may not be expecting
	log.Print(invalidDat)
	err = c.Parse([]byte(invalidDat))
	assert.NotNil(err)
}
