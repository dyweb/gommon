# vfsgen

- https://github.com/shurcooL/vfsgen
- generate struct for each file, some are compressed, some are not
  - [ ] how does it determine which to compress? check the size of compressed and raw?
  - making use of the compressed data relies on other libraries
- it supports dir (unlike statik)

Template (in generator.go)

````text
{{define "Time"}}
{{- if .IsZero -}}
	time.Time{}
{{- else -}}
	time.Date({{.Year}}, {{printf "%d" .Month}}, {{.Day}}, {{.Hour}}, {{.Minute}}, {{.Second}}, {{.Nanosecond}}, time.UTC)
{{- end -}}
{{end}}
````

````go
func (d *vfsgen۰Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}
````