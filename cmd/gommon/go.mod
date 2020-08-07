go 1.14

require (
	github.com/dyweb/gommon v0.0.13
	github.com/spf13/cobra v0.0.6
	golang.org/x/tools v0.0.0-20200401192744-099440627f01
)

replace github.com/dyweb/gommon => ../..

// TODO: might name is gom
// NOTE: rename it to gommonbin to aviod ambiguous import
// can't load package: package github.com/dyweb/gommon/cmd/gommon: ambiguous import: found github.com/dyweb/gommon/cmd/gommon in multiple modules:
//           github.com/dyweb/gommon/cmd/gommon (/home/at15/w/src/github.com/dyweb/gommon/cmd/gommon)
//           github.com/dyweb/gommon v0.0.13 (/home/at15/w/pkg/mod/github.com/dyweb/gommon@v0.0.13/cmd/gommon)
module github.com/dyweb/gommon/cmd/gommonbin
