package embed

import (
	"os"
	"time"

	"github.com/dyweb/gommon/noodle"
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

		"assets": {
			FileInfo: FileInfo{
				name:    "assets",
				size:    4096,
				mode:    2147484157,
				modTime: time.Unix(1518819613, 0),
				isDir:   true,
			},
			Entries: []FileInfo{

				{
					name:    "404",
					size:    4096,
					mode:    2147484157,
					modTime: time.Unix(1518678160, 0),
					isDir:   true,
				},

				{
					name:    "idx",
					size:    4096,
					mode:    2147484157,
					modTime: time.Unix(1518819613, 0),
					isDir:   true,
				},

				{
					name:    "noidx",
					size:    4096,
					mode:    2147484157,
					modTime: time.Unix(1518677629, 0),
					isDir:   true,
				},

				{
					name:    ".noodleignore",
					size:    147,
					mode:    436,
					modTime: time.Unix(1518819613, 0),
					isDir:   false,
				},

				{
					name:    "index.html",
					size:    0,
					mode:    436,
					modTime: time.Unix(1518677256, 0),
					isDir:   false,
				},
			},
		},

		"404": {
			FileInfo: FileInfo{
				name:    "404",
				size:    4096,
				mode:    2147484157,
				modTime: time.Unix(1518678160, 0),
				isDir:   true,
			},
			Entries: []FileInfo{},
		},

		"idx": {
			FileInfo: FileInfo{
				name:    "idx",
				size:    4096,
				mode:    2147484157,
				modTime: time.Unix(1518819613, 0),
				isDir:   true,
			},
			Entries: []FileInfo{

				{
					name:    "index.html",
					size:    180,
					mode:    436,
					modTime: time.Unix(1518677561, 0),
					isDir:   false,
				},

				{
					name:    "main.css",
					size:    38,
					mode:    436,
					modTime: time.Unix(1518677599, 0),
					isDir:   false,
				},

				{
					name:    "main.js",
					size:    51,
					mode:    436,
					modTime: time.Unix(1518677599, 0),
					isDir:   false,
				},
			},
		},

		"noidx": {
			FileInfo: FileInfo{
				name:    "noidx",
				size:    4096,
				mode:    2147484157,
				modTime: time.Unix(1518677629, 0),
				isDir:   true,
			},
			Entries: []FileInfo{

				{
					name:    "main.css",
					size:    37,
					mode:    436,
					modTime: time.Unix(1518677629, 0),
					isDir:   false,
				},

				{
					name:    "main.js",
					size:    50,
					mode:    436,
					modTime: time.Unix(1518677629, 0),
					isDir:   false,
				},

				{
					name:    "noindex.html",
					size:    192,
					mode:    436,
					modTime: time.Unix(1518677561, 0),
					isDir:   false,
				},
			},
		},
	}

}
