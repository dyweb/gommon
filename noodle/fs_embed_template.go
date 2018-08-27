package noodle

var embedTemplate = `
import (
	"time"

	"github.com/dyweb/gommon/noodle"
)

// GetNoodle{{ .name }} returns an extracted EmbedBowl
func GetNoodle{{ .name }} () (noodle.EmbedBowel, error){

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

	data := {{ printf "%#v" .data }}
	bowl := noodle.EmbedBowel{
		Dirs: dirs,
		Data: data,
	}
	if err := bowl.ExtractFiles(); err != nil {
		return bowl, err
	}
    return bowl, nil
}
`
