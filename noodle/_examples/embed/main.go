package main

import (
	"log"

	"github.com/dyweb/gommon/noodle"
)

// need to include t.go
// go run main.go t.go
func main() {
	box, err := noodle.GetEmbedBox("test")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(box.Dirs))
}
