package log

import "testing"

func TestEntrylog(t *testing.T) {
	logger := NewLogger()
	entry := logger.NewEntry()
	entry.AddField("pkg", "dummy.d")
	entry.AddField("name", "jack")
	entry.log(DebugLevel, "hahaha")
	entry.Infof("%s %d", "haha", 1)
}
