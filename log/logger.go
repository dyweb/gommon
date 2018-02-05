package log

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/dyweb/gommon/structure"
	"time"
)

type Logger struct {
	mu       sync.RWMutex
	h        Handler
	level    Level
	fields   Fields
	children map[string][]*Logger
	id       *Identity
}

func (l *Logger) Level() Level {
	return l.level
}

func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	l.level = level
	l.mu.Unlock()
}

func (l *Logger) SetHandler(h Handler) {
	l.mu.Lock()
	l.h = h
	l.mu.Unlock()
}

func (l *Logger) Identity() *Identity {
	return l.id
}

// TODO: allow release a child logger
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

func (l *Logger) Panic(args ...interface{}) {
	s := fmt.Sprint(args...)
	l.h.HandleLog(PanicLevel, time.Now(), s)
	l.h.Flush()
	panic(s)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.Panic(fmt.Sprintf(format, args))
}

func (l *Logger) Fatal(args ...interface{}) {
	s := fmt.Sprint(args...)
	l.h.HandleLog(FatalLevel, time.Now(), s)
	l.h.Flush()
	// TODO: allow user to register hook to do cleanup before exit directly
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(format, args))
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

//// TODO: deal w/ http access log later
//type HttpAccessLogger struct {
//}
