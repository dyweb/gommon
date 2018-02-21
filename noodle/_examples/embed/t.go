package main

import (
	"time"

	"github.com/dyweb/gommon/noodle"
)

func init() {

	dirs := map[string]noodle.EmbedDir{"_examples/embed/assets": {
		FileInfo: noodle.FileInfo{
			FileName:    "assets",
			FileSize:    4096,
			FileMode:    2147484157,
			FileModTime: time.Unix(1518819613, 0),
			FileIsDir:   true,
		},
		Entries: []noodle.FileInfo{{
			FileName:    "404",
			FileSize:    4096,
			FileMode:    2147484157,
			FileModTime: time.Unix(1518678160, 0),
			FileIsDir:   true,
		}, {
			FileName:    "idx",
			FileSize:    4096,
			FileMode:    2147484157,
			FileModTime: time.Unix(1519187199, 0),
			FileIsDir:   true,
		}, {
			FileName:    "noidx",
			FileSize:    4096,
			FileMode:    2147484157,
			FileModTime: time.Unix(1518677629, 0),
			FileIsDir:   true,
		}, {
			FileName:    ".noodleignore",
			FileSize:    147,
			FileMode:    436,
			FileModTime: time.Unix(1518819613, 0),
			FileIsDir:   false,
		}, {
			FileName:    "index.html",
			FileSize:    0,
			FileMode:    436,
			FileModTime: time.Unix(1518677256, 0),
			FileIsDir:   false,
		}},
	}, "_examples/embed/assets/404": {
		FileInfo: noodle.FileInfo{
			FileName:    "404",
			FileSize:    4096,
			FileMode:    2147484157,
			FileModTime: time.Unix(1518678160, 0),
			FileIsDir:   true,
		},
		Entries: []noodle.FileInfo{},
	}, "_examples/embed/assets/idx": {
		FileInfo: noodle.FileInfo{
			FileName:    "idx",
			FileSize:    4096,
			FileMode:    2147484157,
			FileModTime: time.Unix(1519187199, 0),
			FileIsDir:   true,
		},
		Entries: []noodle.FileInfo{{
			FileName:    "sub",
			FileSize:    4096,
			FileMode:    2147484157,
			FileModTime: time.Unix(1519187227, 0),
			FileIsDir:   true,
		}, {
			FileName:    "index.html",
			FileSize:    180,
			FileMode:    436,
			FileModTime: time.Unix(1518677561, 0),
			FileIsDir:   false,
		}, {
			FileName:    "main.css",
			FileSize:    38,
			FileMode:    436,
			FileModTime: time.Unix(1518677599, 0),
			FileIsDir:   false,
		}, {
			FileName:    "main.js",
			FileSize:    51,
			FileMode:    436,
			FileModTime: time.Unix(1518677599, 0),
			FileIsDir:   false,
		}},
	}, "_examples/embed/assets/idx/sub": {
		FileInfo: noodle.FileInfo{
			FileName:    "sub",
			FileSize:    4096,
			FileMode:    2147484157,
			FileModTime: time.Unix(1519187227, 0),
			FileIsDir:   true,
		},
		Entries: []noodle.FileInfo{{
			FileName:    "index.html",
			FileSize:    115,
			FileMode:    436,
			FileModTime: time.Unix(1519187227, 0),
			FileIsDir:   false,
		}},
	}, "_examples/embed/assets/noidx": {
		FileInfo: noodle.FileInfo{
			FileName:    "noidx",
			FileSize:    4096,
			FileMode:    2147484157,
			FileModTime: time.Unix(1518677629, 0),
			FileIsDir:   true,
		},
		Entries: []noodle.FileInfo{{
			FileName:    "main.css",
			FileSize:    37,
			FileMode:    436,
			FileModTime: time.Unix(1518677629, 0),
			FileIsDir:   false,
		}, {
			FileName:    "main.js",
			FileSize:    50,
			FileMode:    436,
			FileModTime: time.Unix(1518677629, 0),
			FileIsDir:   false,
		}, {
			FileName:    "noindex.html",
			FileSize:    192,
			FileMode:    436,
			FileModTime: time.Unix(1518677561, 0),
			FileIsDir:   false,
		}},
	}}

}
