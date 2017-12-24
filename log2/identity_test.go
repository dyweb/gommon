package log2

import "testing"

func TestNewIdentityFromCaller(t *testing.T) {
	NewIdentityFromCaller(0)
	NewIdentityFromCaller(1)
	NewIdentityFromCaller(2)
}
