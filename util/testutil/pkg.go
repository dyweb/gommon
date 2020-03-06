// Package testutil defines helper functions like condition, golden file, docker container etc.
package testutil

import "time"

var testStart time.Time

// SetTestStart allows you to override the global test start time.
// Which is used as a simple identifier to distinguish different tests.
// For instance, it is used as label in container test.
func SetTestStart(t time.Time) {
	testStart = t
}

func init() {
	testStart = time.Now()
}
