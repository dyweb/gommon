// +build ignore

// TODO: enable this file after refactor on tree of logger is finished
package log

import (
	"io"
	"os"

	"github.com/dyweb/gommon/structure"
)

func SetLevelRecursive(root *Logger, level Level) {
	visited := make(map[*Logger]bool)
	PreOrderDFS(root, visited, func(l *Logger) {
		// TODO: remove it after we have tested it ....
		//fmt.Println(l.Identity().String())
		l.SetLevel(level)
	})
}

func SetHandlerRecursive(root *Logger, handler Handler) {
	visited := make(map[*Logger]bool)
	PreOrderDFS(root, visited, func(l *Logger) {
		l.SetHandler(handler)
	})
}

// FIXME: this fixed typo requires update in go.ice
func EnableSourceRecursive(root *Logger) {
	visited := make(map[*Logger]bool)
	PreOrderDFS(root, visited, func(l *Logger) {
		l.EnableSource()
	})
}

func DisableSourceRecursive(root *Logger) {
	visited := make(map[*Logger]bool)
	PreOrderDFS(root, visited, func(l *Logger) {
		l.DisableSource()
	})
}

// TODO: test it .... map traverse order is random, we need radix tree, it is need for pretty print as well
func PreOrderDFS(root *Logger, visited map[*Logger]bool, cb func(l *Logger)) {
	if visited[root] {
		return
	}
	cb(root)
	visited[root] = true
	for _, group := range root.children {
		for _, l := range group {
			PreOrderDFS(l, visited, cb)
		}
	}
}

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
