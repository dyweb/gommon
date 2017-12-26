package main

import (
	"fmt"

	// you can use package with _ prefix, but it is ignored by go build?
	"github.com/dyweb/gommon/log2/_examples/uselib/foo"
	"github.com/dyweb/gommon/log2/_examples/uselib/service"
	"github.com/dyweb/gommon/log2/_examples/uselib/storage"

	_ "github.com/dyweb/gommon/log2/_examples/uselib/storage/mem"
)



// TODO: make it a full example, like user auth with multiple mock backends
func main() {
	fmt.Println(foo.FOO)
	auth := service.NewAuth(storage.Get("mem"))
	auth.Check("jack", "123")
}
