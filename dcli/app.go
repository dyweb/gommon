package dcli

import (
	"context"
	"os"

	"github.com/dyweb/gommon/errors"
)

// app.go defines application struct, a wrapper for top level command.

type Application struct {
	Build BuildInfo
	Root  Command // entry command, its Name should be same as Application.Name but it is ignored when execute.
}

// RunApplication creates a new application and run it directly.
// It logs and exit with 1 if application creation or execution failed.
func RunApplication(cmd Command) {
	app, err := NewApplication(cmd)
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}

const versionCmd = "version"

// NewApplication validate root command and injects version command if not exists.
func NewApplication(cmd Command) (*Application, error) {
	if err := ValidateCommand(cmd); err != nil {
		return nil, errors.Wrap(err, "command validation failed")
	}
	info := DefaultBuildInfo()
	// Inject version command if it does not exist
	if !hasChildCommand(cmd, versionCmd) {
		// TODO: a better way is to wrap it so we don't modify original command
		// or a new interface for mutable command that allows adding command
		c, ok := cmd.(*Cmd)
		if ok {
			//log.Info("adding version command")
			c.Children = append(c.Children, &Cmd{
				Name: versionCmd,
				Run: func(_ context.Context) error {
					PrintBuildInfo(os.Stdout, info)
					return nil
				},
			})
		}
	}
	return &Application{
		Build: info,
		Root:  cmd,
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
	// TODO: extra arg and flags
	//log.Infof("args %v", args)
	c, err := FindCommand(a.Root, args)
	if err != nil {
		return err
	}
	return c.GetRunnable()(ctx)
}
