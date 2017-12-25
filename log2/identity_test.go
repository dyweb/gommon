package log2

import "testing"

var _ = NewIdentityFromCaller(0)

func TestNewIdentityFromCaller(t *testing.T) {
	NewIdentityFromCaller(0)
	NewIdentityFromCaller(1)
	NewIdentityFromCaller(2)
}
