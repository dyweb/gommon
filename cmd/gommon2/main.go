// gommon2 is a test binary for using dcli, it will be renamed to gommon once we have most functionality in spf13/cobra
package main

import (
	"fmt"

	"github.com/dyweb/gommon/dcli"
)

func main() {
	fmt.Printf("info %v", dcli.DefaultBuildInfo())
}
