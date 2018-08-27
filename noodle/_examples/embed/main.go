package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dyweb/gommon/noodle/_examples/embed/gen"
)

func main() {
	bowel1, err := gen.GetNoodleYangchunMian()
	if err != nil {
		log.Fatal(err)
	}
	addr := ":8080"
	var root http.FileSystem
	root = &bowel1
	fmt.Printf("listen on %s\n", addr)
	http.ListenAndServe(addr, http.FileServer(root))
}
