// +build ignore

package structure

import (
	"bytes"
	"testing"

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

	// FIXED: there were extra vertical lines
	// main
	// └── http
	// │    └── auth
	root.Children = root.Children[0:1]
	expected =
		`root
└── level1-A
     └── level2-A
`
	b.Reset()
	root.PrintTo(&b)
	assert.Equal(expected, string(b.Bytes()))

	// TODO: test more complex situation like this
	//.
	//├── benchgraph
	//│   ├── echart.go
	//│   ├── echart_test.go
	//│   │   └── zap-no-delete-field.txt
	//│   │
	//│   └── fixture
	//│       ├── zap-no-delete-field.txt
	//│       └── zap.txt
	//├── gommon
	//│   └── main.go
	//└── README.md
}
