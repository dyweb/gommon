package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/dyweb/gommon/noodle"
	_ "github.com/dyweb/gommon/noodle/_examples/embed/gen"
)

func main() {
	box, err := noodle.GetEmbedBowel("test")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(box.Dirs))
	zipReader, err := zip.NewReader(bytes.NewReader(box.Data), int64(len(box.Data)))
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range zipReader.File {
		log.Println(f.Name)
	}
	if err := box.ExtractFiles(); err != nil {
		log.Fatal(err)
	}
	addr := ":8080"
	var root http.FileSystem
	root = &box
	fmt.Printf("listen on %s\n", addr)
	http.ListenAndServe(addr, http.FileServer(root))
}
