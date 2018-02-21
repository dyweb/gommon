package main

import (
	"log"

	"github.com/dyweb/gommon/noodle"
	"archive/zip"
	"bytes"
)

// need to include t.go
// go run main.go t.go
func main() {
	box, err := noodle.GetEmbedBox("test")
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
}
