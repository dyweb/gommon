package testutil

import "go/build"

// GOPATH returns build.Default.GOPATH
// https://stackoverflow.com/questions/32649770/how-to-get-current-gopath-from-code
func GOPATH() string {
	return build.Default.GOPATH
}

