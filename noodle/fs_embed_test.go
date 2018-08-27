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

func TestGenerateEmbed(t *testing.T) {
	//testutil.RunIf(t, testutil.IsTravis())
	testutil.SkipIf(t, testutil.IsTravis())

	require := requir.New(t)

	b, err := noodle.GenerateEmbeds([]noodle.EmbedConfig{
		{
			Root: "_examples/embed/assets",
			Name: "YangchunMian",
		},
		{
			Root: "_examples/embed/third_party",
			Name: "BieRenJiaDeMian",
		},
	}, "gen")
	require.Nil(err)
	fsutil.WriteFile("_examples/embed/gen/noodle.go", b)
}

func TestFileMode(t *testing.T) {
	// 2147483648
	fmt.Printf("%#o", os.ModeDir)
}
