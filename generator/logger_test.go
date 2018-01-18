package generator

import (
	"os"
	"testing"
)

func TestLoggerConfig_RenderTo(t *testing.T) {
	c := LoggerConfig{Struct: "*YAMLConfig", Receiver: "c"}
	c.RenderTo(os.Stdout)
}
