package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dyweb/gommon/noodle/_examples/embed/gen"
)

// use local folder dev, use embed when prod
func main() {
	mode := "dev"
	if len(os.Args) > 1 {
		if strings.HasPrefix(os.Args[1], "p") {
			mode = "prod"
		}
	}
	var root http.FileSystem
	if mode == "dev" {
		// NOTE: if you are running it in IDE, use the following full path instead of relative path
		//localDir := os.Getenv("GOPATH") + "/src/github.com/dyweb/gommon/noodle/_examples/embed/assets"
		localDir := "assets"
		root = http.Dir(localDir)
		//root = noodle.NewLocal("assets")
	} else {
		bowel1, err := gen.GetNoodleYangChunMian()
		if err != nil {
			log.Fatal(err)
		}
		root = &bowel1
	}
	addr := ":8080"
	fmt.Printf("listen on %s in %s mode\n", addr, mode)
	fmt.Printf("use http://localhost:8080/index.html")
	log.Fatal(http.ListenAndServe(addr, http.FileServer(root)))
}
