package dcli

import "runtime"

type Application struct {
	Name        string
	Description string
	Version     string
}

var (
	// set using -ldflags "-X github.com/dyweb/gommon/dcli.buildVersion=0.0.1"
	buildVersion string
	buildCommit  string
	buildBranch  string
	buildTime    string
	buildUser    string
)

// BuildInfo contains information that should be set at build time.
// e.g. go install ./cmd/myapp -ldflags "-X github.com/dyweb/gommon/dcli.buildVersion=0.0.1"
// You can use DefaultBuildInfo and copy paste its Makefile rules.
type BuildInfo struct {
	Version   string
	Commit    string
	Branch    string
	Time      string
	User      string
	GoVersion string
}

// DefaultBuildInfo returns a info based on ld flags sets to github.com/dyweb/gommon/dcli.*
// You can copy the following rules in your Makefile
//
// DCLI_PKG = github.com/dyweb/gommon/dcli.
// DCLI_LDFLAGS = -X $(DCLI_PKG)buildVersion=$(VERSION) -X $(DCLI_PKG)buildCommit=$(BUILD_COMMIT) -X $(DCLI_PKG)buildBranch=$(BUILD_BRANCH) -X $(DCLI_PKG)/buildTime=$(BUILD_TIME) -X $(DCLI_PKG)buildUser=$(CURRENT_USER)
//
// install:
// 	go install -ldflags $(DCLI_LDFLAGS) ./cmd/myapp
func DefaultBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   buildVersion,
		Commit:    buildCommit,
		Branch:    buildBranch,
		Time:      buildTime,
		User:      buildUser,
		GoVersion: runtime.Version(),
	}
}
