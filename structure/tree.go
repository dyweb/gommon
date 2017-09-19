package structure

import (
	"io"
	"os"
)

var (
	vLineStart = []byte("├") // this is in extended ASCII http://www.theasciicode.com.ar/extended-ascii-code/box-drawings-single-line-vertical-right-character-ascii-code-195.html
	vLine      = []byte("│") // this is not |, it won't have space vertically
	corner     = []byte("└")
	hLine      = []byte("── ") // we got a space at the end of hLine
	space      = []byte(" ")
	hSpace     = []byte("    ")
	nextLine   = []byte("\n")
)

type StringTreeNode struct {
	Val      string
	Children []StringTreeNode
}

func (tree *StringTreeNode) Append(child StringTreeNode) {
	tree.Children = append(tree.Children, child)
}

func (tree *StringTreeNode) Print() {
	tree.PrintTo(os.Stdout)
}

/*
tree command example output
.
├── benchgraph
│   ├── echart.go
│   ├── echart_test.go
│   ├── fixture
│   │   ├── zap-no-delete-field.txt
│   │   └── zap.txt
│   ├── main.go
│   ├── Makefile
│   ├── tsdb-bench.html
│   └── zap.html
├── gommon
│   └── main.go
└── README.md
*/

func (tree *StringTreeNode) PrintTo(w io.Writer) {
	treePrintHelper(tree, 0, false, w)
}

func treePrintHelper(tree *StringTreeNode, level int, lastOfUs bool, w io.Writer) {
	// print the prefix before me, both vertically and horizontally
	for i := 0; i < level-1; i++ {
		w.Write(vLine)
		w.Write(hSpace)
	}
	if level != 0 {
		if !lastOfUs {
			w.Write(vLineStart)
		} else {
			w.Write(corner)
		}
		w.Write(hLine)
	}
	// my value
	w.Write([]byte(tree.Val))
	w.Write(nextLine)
	// children
	level++
	n := len(tree.Children)
	for i := 0; i < n-1; i++ {
		treePrintHelper(&tree.Children[i], level, false, w)
	}
	if n > 0 {
		treePrintHelper(&tree.Children[n-1], level, true, w)
	}
}
