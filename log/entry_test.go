package log

import "testing"

func TestEntrylog(t *testing.T) {
	logger := NewLogger()
	entry := logger.NewEntry()
	entry.AddField("pkg", "dummy.d")
	entry.AddField("name", "jack")
	entry.log(DebugLevel, "hahaha")

}

func TestEntryLevelLog(t *testing.T){
	logger := NewLogger()
	f := NewTextFormatter()
	f.EnableColor = true
	logger.Formatter = f
	entry := logger.NewEntry()
	// NOTE: when run it in IDEA, the terminal does not have color, run it in real terminal it will show
	entry.Debug("You should not see me!")
	entry.Infof("%s %d", "haha", 1)
	entry.Warnf("%s %d", "haha", 1)
	entry.Errorf("%s %d", "haha", 1)

	logger2 := NewLogger()
	f2 := NewTextFormatter()
	f2.EnableColor = true
	f2.EnableElapsedTime = false
	logger2.Formatter = f2
	entry2 := logger2.NewEntry()
	entry2.Info("I should have full timestamp")

	logger3 := NewLogger()
	logger3.Level = DebugLevel
	f3 := NewTextFormatter()
	f3.EnableColor = false
	f3.EnableTimeStamp = false
	logger3.Formatter = f3
	entry3 := logger3.NewEntry()
	entry3.Info("I should have no timestamp")
	entry3.Debug("You should see me!")
}