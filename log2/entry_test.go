package log2

import (
	"testing"
)

func TestEntry_PrintEntryTree(t *testing.T) {
	project := newEntry()
	project.Package = "main"
	http := newEntry()
	http.Package = "http"
	// TODO: better have an append function, so the Parent in child is updated as well
	project.Children = append(project.Children, http)
	auth := newEntry()
	auth.Package = "auth"
	http.Children = append(http.Children, auth)
	// FIXME: there is extra vertical line
	project.PrintEntryTree()
}
