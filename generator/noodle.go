package generator

import (
	"github.com/dyweb/gommon/errors"
)

type NoodleConfig struct {
	Name string `yaml:"name"`
}

// IsGo returns true TODO: until noodle and generate zip file and attach to end of binary
func (c *NoodleConfig) IsGo() bool {
	return true
}

// Render is not implemented
func (c *NoodleConfig) Render(root string) error {
	return errors.New("not implemented")
}

func (c *NoodleConfig) Imports() []string {
	return []string{
		"time",
		"github.com/dyweb/gommon/noodle",
	}
}

func (c *NoodleConfig) RenderBody() ([]byte, error) {
	return nil, nil
}
