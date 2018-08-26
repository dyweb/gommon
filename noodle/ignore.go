package noodle

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/fsutil"
)

func ReadIgnoreFile(path string) (*fsutil.Ignores, error) {
	if f, err := os.Open(path); err != nil {
		return nil, errors.Wrap(err, "can't open ignore file")
	} else {
		defer f.Close()
		return ReadIgnore(f)
	}
}

// ReadIgnore reads ignore file and change to patterns
// reading is based on https://github.com/codeskyblue/dockerignore/blob/HEAD/ignore.go#L40
func ReadIgnore(reader io.Reader) (*fsutil.Ignores, error) {
	ignores := fsutil.NewIgnores(nil, nil)
	if reader == nil {
		return ignores, nil
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := CleanLine(scanner.Text())
		if len(line) == 0 {
			continue
		}
		if strings.Contains(line, string(filepath.Separator)) {
			ignores.AddPath(LineToPattern(line))
		} else {
			ignores.AddName(LineToPattern(line))
		}
	}
	if err := scanner.Err(); err != nil {
		return ignores, errors.Wrap(err, "scanner error")
	}
	return ignores, nil
}

// CleanLine remove trailing space and \n, and anything following the first #
func CleanLine(line string) string {
	var buf []rune
Outer:
	for _, c := range line {
		switch c {
		case '#', '\r', '\n':
			break Outer
		default:
			buf = append(buf, c)
		}
	}
	if len(buf) == 0 {
		return ""
	}
	return strings.TrimSpace(string(buf))
}

func LineToPattern(line string) fsutil.IgnorePattern {
	for _, c := range line {
		switch c {
		case '*', '?':
			return fsutil.WildcardPattern(line)
		}
	}
	return fsutil.ExactPattern(line)
}
