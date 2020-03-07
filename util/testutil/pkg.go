// Package testutil defines helper functions like condition, golden file, docker container etc.
//
// Condition allows you to skip/run test with message explaining why it is skipped, e.g. env RUN_ABC=false, skip ABC.
//
// Container is used to run one off test containers by shelling out to docker cli.
//
// Golden switch golden file generation and validation based on environment variable using condition.
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
