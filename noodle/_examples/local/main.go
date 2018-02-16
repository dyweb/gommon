package main

import (
	"fmt"
	"github.com/dyweb/gommon/noodle"
	"net/http"
	"os"
)

func main() {
	fs := "default"
	if len(os.Args) > 1 {
		fs = os.Args[1]
	}
	fmt.Println(fs)
	addr := ":8080"
	//addr := ":6667" # https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers 6665-6669 is for IRC ...
	fmt.Printf("listen on %s\n", addr)
	var root http.FileSystem
	switch fs {
	case "default":
		root = http.Dir("./assets")
	case "local":
		root = noodle.NewLocal("./assets")
	case "local-unsafe":
		root = noodle.NewLocalUnsafe("./assets")
	default:
		panic("unknown fs")
	}
	http.ListenAndServe(addr, http.FileServer(root))

}
