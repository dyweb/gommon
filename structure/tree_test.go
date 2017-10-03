package structure

import (
	"testing"

	"bytes"

	asst "github.com/stretchr/testify/assert"
)

func TestStringTreeNode_PrintTo(t *testing.T) {
	assert := asst.New(t)
	expected :=
		`root
├── level1-A
│    └── level2-A
└── level1-B
`
	root := StringTreeNode{Val: "root"}
	root.Append(StringTreeNode{Val: "level1-A"})
	root.Children[0].Append(StringTreeNode{Val: "level2-A"})
	root.Append(StringTreeNode{Val: "level1-B"})
	var b bytes.Buffer
	root.PrintTo(&b)
	assert.Equal(expected, string(b.Bytes()))

	// FIXME: there are extra vertical lines
	// main
	// └── http
	// │    └── auth
}
