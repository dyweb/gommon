package noodle

var embedTemplate = `
import (
	"time"

	"github.com/dyweb/gommon/noodle"
)

func init() {

dirs := map[string]noodle.EmbedDir{
{{- range $path, $dir := .dir -}}
	"{{ $path }}": {
		FileInfo: noodle.FileInfo{
			FileName: "{{ $dir.FileInfo.Name }}",
			FileSize: {{ $dir.FileInfo.Size }},
			FileMode: {{ printf "%#o" $dir.FileInfo.Mode }},
			FileModTime: time.Unix({{$dir.FileInfo.ModTime.Unix }}, 0),
			FileIsDir: {{ $dir.FileInfo.IsDir }},
		},
		Entries: []noodle.FileInfo{
			{{- range $dir.Entries -}}
			{
				FileName: "{{ .Name }}",
				FileSize: {{ .Size }},
				FileMode: {{ printf "%#o" .Mode }},
				FileModTime: time.Unix({{.ModTime.Unix }}, 0),
				FileIsDir: {{ .IsDir }},
			},
			{{- end -}}
        },
	},
{{- end -}}
}

box := noodle.EmbedBox{
	Dirs: dirs,
	Data: {{ printf "%#v" .data }},
}

noodle.RegisterEmbedBox("test", box)

}
`
