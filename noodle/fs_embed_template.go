package noodle

var embedTemplate = `
package noodle

import "time"

var dirs = map[string]embedDir{
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
`
