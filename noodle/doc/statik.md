# Statik

- https://github.com/rakyll/statik
- generate a go file with all the content as zipped bytes
- don't allow list dir by returning empty slice for `Readdir`
  - i.e. even you want to list dir, it is not allowed
  - [ ] it is using a map, may need to sort names and figure out which one is folder, but `archive/zip` might support this?

````go
func init() {
	data := `PK\x03\x04\x14\x00\x08\x00\x08\x00S\x07LL\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00	
\x00\x00\x00hello.txt\xf2H\xcd\xc9\xc9W\x08\xcf/\xcaI\xe1\xe2\x02\x04\x00\x00\xff\xffPK\x07\x08\xb4 
`
	fs.Register(data)
}
````

- `fs/fs.go`
  - handles `index.html` ... I was thinking `net/http` would handle this ... but it seems not ...
    - in `serveFile`, the file have to be a dir in order to have `/index.html` to be checked ...
  - disallow list directory 

````go
func (fs *statikFS) Open(name string) (http.File, error) {
	name = strings.Replace(name, "//", "/", -1)
	f, ok := fs.files[name]
	if ok {
		return newHTTPFile(f, false), nil
	}
	// The file doesn't match, but maybe it's a directory,
	// thus we should look for index.html
	indexName := strings.Replace(name+"/index.html", "//", "/", -1)
	f, ok = fs.files[indexName]
	if !ok {
		return nil, os.ErrNotExist
	}
	return newHTTPFile(f, true), nil
}

// Readdir returns an empty slice of files, directory
// listing is disabled.
func (f *httpFile) Readdir(count int) ([]os.FileInfo, error) {
	// directory listing is disabled.
	return make([]os.FileInfo, 0), nil
}
````
