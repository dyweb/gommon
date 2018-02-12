package main

import (
	"fmt"
	"net/http"
)

func handWrittenFs() *FileSystem {
	a := NewFile("a.txt", []byte("I am a.txt"), false)
	b := NewFile("dir/a.txt", []byte("I am dir/a.txt"), false)
	return NewFs(a, b)
}

// NOTE: need go run main.go data_hand_written.go ....
func main() {

	addr := ":8080"
	var root http.FileSystem
	root = handWrittenFs()
	fmt.Printf("listen on %s\n", addr)
	http.ListenAndServe(addr, http.FileServer(root))
}
