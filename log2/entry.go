package log2

import "io"

type EntryType uint8

type Entry struct {
	Parent          *Entry
	Children        []*Entry
	IncludePackage  bool
	IncludeFunction bool
	IncludeLine     bool
	Package         string
	Function        string
	Line            string // TODO: maybe int?
}

func (e *Entry) PrintEntryTree(w io.Writer) {

}

func (e *Entry) PrintEntireEntryTree(w io.Writer) {
	// find the root
	root := e
	for root.Parent != nil {
		root = root.Parent
	}
	// print from root
	root.PrintEntryTree(w)
}
