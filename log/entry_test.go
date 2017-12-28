package log2

import (
	"bytes"
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestEntry_PrintEntryTree(t *testing.T) {
	assert := asst.New(t)
	project := newEntry()
	project.Package = "main"
	http := newEntry()
	http.Package = "http"
	// TODO: better have an append function, so the Parent in child is updated as well
	project.Children = append(project.Children, http)
	auth := newEntry()
	auth.Package = "auth"
	http.Children = append(http.Children, auth)
	// FIXED: there was an extra vertical line in front of auth
	var b bytes.Buffer
	expected :=
		`main
└── http
     └── auth
`
	project.PrintEntryTreeTo(&b)
	assert.Equal(expected, string(b.Bytes()))
}
