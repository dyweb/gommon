package noodle

var embedTemplate = `
package {{ .pkg }}

import (
	"time"

	"github.com/dyweb/gommon/noodle"
)

func init() {

dirs := map[string]noodle.EmbedDir{
{{- range .dir -}}
	"{{ .FileInfo.Name }}": {
		FileInfo: noodle.FileInfo{
			FileName: "{{ .FileInfo.Name }}",
			FileSize: {{ .FileInfo.Size }},
			FileMode: {{ printf "%#0d" .FileInfo.Mode }},
			FileModTime: time.Unix({{.FileInfo.ModTime.Unix }}, 0),
			FileIsDir: {{ .FileInfo.IsDir }},
		},
		Entries: []noodle.FileInfo{
			{{- range .Entries -}}
			{
				FileName: "{{ .Name }}",
				FileSize: {{ .Size }},
				FileMode: {{ printf "%#0d" .Mode }},
				FileModTime: time.Unix({{.ModTime.Unix }}, 0),
				FileIsDir: {{ .IsDir }},
			},
			{{- end -}}
        },
	},
{{- end -}}

}


}
`
