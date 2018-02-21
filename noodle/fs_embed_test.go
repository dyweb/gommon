package noodle

import (
	"fmt"
	"os"
	"testing"
)

func TestGenerateEmbed(t *testing.T) {
	GenerateEmbed("_examples/embed/assets")
}

func TestFileMode(t *testing.T) {
	// 2147483648
	fmt.Printf("%#0d", os.ModeDir)
}
