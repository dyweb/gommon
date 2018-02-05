package log

import (
	"fmt"
	"io"
	"os"

	"github.com/dyweb/gommon/structure"
)

// TODO: allow release a child logger, this will be a trouble if we created 1,000 Client struct with its own logger...
func (l *Logger) AddChild(child *Logger) {
	l.mu.Lock()
	if l.children == nil {
		l.children = make(map[string][]*Logger, 1)
	}
	// children are group by their identity, i.e a package logger may have many struct logger of same struct because
	// that struct is used in multiple goroutines, those loggers have different pointer for identity, but they should
	// have same source line, so we use SourceLocation as key
	k := child.id.SourceLocation()
	children := l.children[k]
	// avoid putting same pointer twice, though it should never happen if AddChild is called correctly
	exists := false
	for _, c := range children {
		if c == child {
			exists = true
			break
		}
	}
	if !exists {
		l.children[k] = append(children, child)
	}
	l.mu.Unlock()
}

func SetLevelRecursive(root *Logger, level Level) {
	visited := make(map[*Logger]bool)
	PreOrderDFS(root, visited, func(l *Logger) {
		// TODO: remove it after we have tested it ....
		fmt.Println(l.Identity().String())
		l.SetLevel(level)
	})
}

func SetHandlerRecursive(root *Logger, handler Handler) {
	visited := make(map[*Logger]bool)
	PreOrderDFS(root, visited, func(l *Logger) {
		l.SetHandler(handler)
	})
}

// TODO: test it ....
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
//app logger /home/at15/workspace/src/github.com/at15/go.ice/_example/github/pkg/util/logutil/pkg.go:8
//└── lib logger /home/at15/workspace/src/github.com/at15/go.ice/ice/util/logutil/pkg.go:8
//		└── lib logger /home/at15/workspace/src/github.com/dyweb/gommon/util/logutil/pkg.go:7
//│         └── pkg logger /home/at15/workspace/src/github.com/dyweb/gommon/config/pkg.go:8

func (l *Logger) PrintTreeTo(w io.Writer) {
	st := ToStringTree(l)
	st.PrintTo(w)
}
