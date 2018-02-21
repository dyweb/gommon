package noodle

import (
	"testing"
	"fmt"
	"os"
)

func TestGenerateEmbed(t *testing.T) {
	GenerateEmbed("_examples/embed/assets")
}


func TestFileMode(t *testing.T) {
	// 2147483648
	fmt.Printf("%#0d", os.ModeDir)
}