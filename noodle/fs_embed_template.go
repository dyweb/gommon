package noodle

var embedTemplate = `
package {{ .pkg }}

import (
	"os"
	"time"
)

type embedFile struct {
	FileInfo
	data []byte
}

type embedDir struct {
	FileInfo
	Entries []FileInfo
}

type FileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func init() {

dirs := map[string]embedDir{
{{ range .dir }}
	"{{ .FileInfo.Name }}": {
		FileInfo: FileInfo{
			name: "{{ .FileInfo.Name }}",
			size: {{ .FileInfo.Size }},
			mode: {{ printf "%#0d" .FileInfo.Mode }},
			modTime: time.Unix({{.FileInfo.ModTime.Unix }}, 0),
			isDir: {{ .FileInfo.IsDir }},
		},
		Entries: []FileInfo{
			{{ range .Entries }}
			{
				name: "{{ .Name }}",
				size: {{ .Size }},
				mode: {{ printf "%#0d" .Mode }},
				modTime: time.Unix({{.ModTime.Unix }}, 0),
				isDir: {{ .IsDir }},
			},
			{{ end }}
        },
	},
{{ end }}
}


}
`
