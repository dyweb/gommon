package generator

import (
	"io/ioutil"
	"os"
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func Test_Walk(t *testing.T) {
	assert := asst.New(t)
	if testing.Short() {
		t.Skip()
		return
	}
	//files := walk("../config", defaultIgnores)
	files := Walk("testdata", DefaultIgnores())
	//t.Log(files)
	assert.Equal([]string{
		"testdata/pkg/client/gommon.yml",
		"testdata/pkg/gommon.yml",
		"testdata/pkg/server/gommon.yml",
	}, files)
}

// https://github.com/dyweb/gommon/issues/41
// NOTE: octal literal is needed for file permission, and it starts with 0 ...
func Test_WriteFile(t *testing.T) {
	t.Skip("breaks on travis") // https://travis-ci.org/dyweb/gommon/jobs/337322278
	assert := asst.New(t)
	ioutil.WriteFile("/tmp/mod_666", []byte("you can't see me"), 666)
	info, _ := os.Stat("/tmp/mod_666")
	assert.Equal("--w--wx---", info.Mode().String())
	ioutil.WriteFile("/tmp/mod_0666", []byte("you can see me now"), 0666)
	info, _ = os.Stat("/tmp/mod_0666")
	assert.Equal("-rw-rw-r--", info.Mode().String())
	WriteFile("/tmp/mod_0664", []byte("normal file mode"))
	assert.Equal("-rw-rw-r--", info.Mode().String())
}
