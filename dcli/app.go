package dcli

import (
	"context"
	"os"
	"runtime"

	"github.com/dyweb/gommon/errors"
)

// app.go defines application struct and build info.

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

type Application struct {
	Description string
	Version     string

	name    string  // binary name
	command Command // entry command, its Name should be same as Application.Name but it is ignored when execute.
}

// RunApplication creates a new application and run it directly.
// It logs and exit with 1 if application creation or execution failed.
func RunApplication(name string, cmd Command) {
	app, err := NewApplication(name, cmd)
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}

func NewApplication(name string, cmd Command) (*Application, error) {
	if err := validate(cmd); err != nil {
		return nil, errors.Wrap(err, "command validation failed")
	}
	return &Application{
		Description: "",
		Version:     "",
		name:        name,
		command:     cmd,
	}, nil
}

// Run calls RunArgs with command line arguments (os.Args[1:]) and exit 1 when there is error.
func (a *Application) Run() {
	if err := a.RunArgs(context.Background(), os.Args[1:]); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func (a *Application) RunArgs(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return a.command.GetRun()(ctx)
	}
	return errors.New("not implemented")
}
