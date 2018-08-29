package noodle_test

import (
	"fmt"
	"os"
	"testing"

	requir "github.com/stretchr/testify/require"

	"github.com/dyweb/gommon/noodle"
	"github.com/dyweb/gommon/util/fsutil"
	"github.com/dyweb/gommon/util/testutil"
)

func TestGenerateEmbedBytes(t *testing.T) {
	testutil.RunIf(t, testutil.IsTravis())
	//testutil.SkipIf(t, testutil.IsTravis())

	require := requir.New(t)

	b, err := noodle.GenerateEmbedBytes([]noodle.EmbedConfig{
		{
			Src:     "_examples/embed/assets",
			Name:    "YangchunMian",
			Package: "gen",
		},
		{
			Src:     "_examples/embed/third_party",
			Name:    "BieRenJiaDeMian",
			Package: "gen",
		},
	})
	require.Nil(err)
	fsutil.WriteFile("_examples/embed/gen/noodle.go", b)
}

func TestFileMode(t *testing.T) {
	// 2147483648
	fmt.Printf("%#o", os.ModeDir)
}
