# net/http

Static file serving in standard library

- `net/http/fs.go` [serveFile](https://golang.org/src/net/http/fs.go#L182)
   - the logic the path need to be a directory, then it will call `Open` with `name

````go
func (d Dir) Open(name string) (File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	dir := string(d)
	if dir == "" {
		dir = "."
	}
	fullName := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
	f, err := os.Open(fullName)
	if err != nil {
		return nil, mapDirOpenError(err, fullName)
	}
	return f, nil
}

func serveFile(w ResponseWriter, r *Request, fs FileSystem, name string, redirect bool) {
    const indexPage = "/index.html"
    // ... omitted
	// use contents of index.html for directory, if present
	if d.IsDir() {
		index := strings.TrimSuffix(name, "/") + indexPage
		ff, err := fs.Open(index)
		if err == nil {
			defer ff.Close()
			dd, err := ff.Stat()
			if err == nil {
				name = index
				d = dd
				f = ff
			}
		}
	}
    // ... omitted
}
````