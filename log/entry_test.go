package log

import "testing"

func TestEntrylog(t *testing.T) {
	logger := NewLogger()
	entry := logger.NewEntry()
	entry.log(DebugLevel, "hahaha")
}
