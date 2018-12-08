// +build ignore

package log

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

var tPkgLogger = NewPackageLogger() // lv0

func TestPreOrderDFS(t *testing.T) {
	// TODO: now we use slice, the order is determined
	t.Skip("FIXME: map does not guarantee order, the traverse result is unstable")
	assert := asst.New(t)
	lv1a := NewFunctionLogger(tPkgLogger)
	lv1b := NewFunctionLogger(tPkgLogger)
	lv2a := NewFunctionLogger(lv1a)
	var ids []string
	// FIXME: the result it unstable ... we should use a radix tree, map does not guarantee the order ...
	var expected = []string{
		tPkgLogger.id.SourceLocation(), // 9
		lv1a.id.SourceLocation(),       // 13
		lv2a.id.SourceLocation(),       // 14
		lv1b.id.SourceLocation(),
	}
	PreOrderDFS(tPkgLogger, make(map[*Logger]bool), func(l *Logger) {
		ids = append(ids, l.id.SourceLocation())
	})
	assert.Equal(expected, ids)
}
