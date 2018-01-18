package generator

import "testing"

// TODO: assert
func Test_Walk(t *testing.T) {
	//files := walk("../config", defaultIgnores)
	files := Walk("testdata", DefaultIgnores)
	t.Log(files)
}
