// +build ignore

// TODO: enable this file after refactor on tree of logger is finished
package log

import (
	"io"
	"os"

	"github.com/dyweb/gommon/structure"
)

func ToStringTree(root *Logger) *structure.StringTreeNode {
	visited := make(map[*Logger]bool)
	return toStringTreeHelper(root, visited)
}

func toStringTreeHelper(root *Logger, visited map[*Logger]bool) *structure.StringTreeNode {
	if visited[root] {
		return nil
	}
	// TODO: might add logger level as well
	n := &structure.StringTreeNode{Val: root.id.String()}
	visited[root] = true
	for _, group := range root.children {
		for _, l := range group {
			p := toStringTreeHelper(l, visited)
			if p != nil {
				n.Append(*p)
			}
		}
	}
	return n
}

// FIXME: it seem without return value and extra parameter, we can't clone a tree?
func toStringTree(root *Logger) *structure.StringTreeNode {
	visited := make(map[*Logger]bool)
	PreOrderDFS(root, visited, func(l *Logger) {

	})
	return nil
}

func (l *Logger) PrintTree() {
	l.PrintTreeTo(os.Stdout)
}

// FIXME: print tree is still having problem ....
//⇒  icehubd log
//	app logger /home/at15/workspace/src/github.com/at15/go.ice/_example/github/pkg/util/logutil/pkg.go:8
//	└── lib logger /home/at15/workspace/src/github.com/at15/go.ice/ice/util/logutil/pkg.go:8
//			└── lib logger /home/at15/workspace/src/github.com/dyweb/gommon/util/logutil/pkg.go:7
//	│         └── pkg logger /home/at15/workspace/src/github.com/dyweb/gommon/config/pkg.go:8

// PrintTreeTo prints logger as a tree, using current logger as root
func (l *Logger) PrintTreeTo(w io.Writer) {
	st := ToStringTree(l)
	st.PrintTo(w)
}
