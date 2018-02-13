package main

import (
	"fmt"
	"net/http"
)

func handWrittenFs() *FileSystem {
	indexContent := `
<html>
<head>
<title>I am title</title>
<link rel="stylesheet" href="main.css">
</head>
<body>
	I am body
	<script src="a.js"></script>
</body>
</html>
`
	mainCss := `
body {
    background-color: lightcyan;
    margin: 20px;
}
`
	aJs := `
console.log({'a': 123});
`
	// TODO: index.html is not supported .... http server redirect browser to /
	index := NewFile("index.html", []byte(indexContent), false)
	index2 := NewFile("index2.html", []byte(indexContent), false)
	css := NewFile("main.css", []byte(mainCss), false)
	js := NewFile("a.js", []byte(aJs), false)
	return NewFs(index, index2, css, js)
}

// NOTE: need go run main.go data_hand_written.go ....
func main() {

	addr := ":8080"
	var root http.FileSystem
	root = handWrittenFs()
	fmt.Printf("listen on %s\n", addr)
	http.ListenAndServe(addr, http.FileServer(root))
}
