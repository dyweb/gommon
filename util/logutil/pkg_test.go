package logutil

import "testing"

func TestNewPackageLoggerAndRegistry(t *testing.T) {
	l, reg := NewPackageLoggerAndRegistry()
	t.Log(l.Identity())
	t.Log(reg.Identity())
}
