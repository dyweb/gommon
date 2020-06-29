// gommon2 is a test binary for using dcli, it will be renamed to gommon once we have most functionality in spf13/cobra
package main

import (
	"context"

	"github.com/dyweb/gommon/dcli"
	dlog "github.com/dyweb/gommon/log"
)

var logReg = dlog.NewRegistry()
var log = logReg.NewLogger()

func main() {
	root := &dcli.Cmd{
		Name: "gommon2",
		Run: func(ctx context.Context) error {
			log.Info("gommon2 does nothing")
			return nil
		},
	}
	dcli.RunApplication(root)
}
