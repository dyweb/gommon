package structure

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
	"bytes"
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
}
