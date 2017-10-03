package log2

import (
	"io"
	"os"

	"github.com/dyweb/gommon/structure"
)

type EntryType uint8

type Entry struct {
	Parent          *Entry
	Children        []*Entry
	IncludePackage  bool
	IncludeFunction bool
	IncludeFile     bool
	Package         string
	Function        string
	File            string // location in source code including line number
}

func newEntry() *Entry {
	// TODO: grab the caller information here?
	return &Entry{
		Parent: nil,
	}
}

func (e *Entry) PrintEntryTree() {
	e.PrintEntireEntryTreeTo(os.Stdout)
}

func (e *Entry) PrintEntryTreeTo(w io.Writer) {
	st := e.ToStringTree()
	st.PrintTo(w)
}

func (e *Entry) ToStringTree() *structure.StringTreeNode {
	// TODO: use package or function or line or all of them?
	root := &structure.StringTreeNode{Val: e.Package}
	for _, child := range e.Children {
		root.Append(*child.ToStringTree())
	}
	return root
}

func (e *Entry) PrintEntireEntryTreeTo(w io.Writer) {
	// find the root
	root := e
	for root.Parent != nil {
		root = root.Parent
	}
	// print from root
	root.PrintEntryTreeTo(w)
}
