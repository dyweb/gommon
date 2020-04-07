package fsutil

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dyweb/gommon/errors"
)

var AcceptAll = NewIgnores(nil, nil)

// ReadIgnoreFile reads a .ignore file and parse the patterns.
func ReadIgnoreFile(path string) (*Ignores, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "can't open ignore file")
	}
	defer f.Close()
	return ReadIgnore(f)
}

// ReadIgnore reads ignore file and change to patterns
// reading is based on https://github.com/codeskyblue/dockerignore/blob/HEAD/ignore.go#L40
func ReadIgnore(reader io.Reader) (*Ignores, error) {
	ignores := NewIgnores(nil, nil)
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

func LineToPattern(line string) IgnorePattern {
	for _, c := range line {
		switch c {
		case '*', '?':
			return WildcardPattern(line)
		}
	}
	return ExactPattern(line)
}

type Ignores struct {
	names  []IgnorePattern
	paths  []IgnorePattern
	prefix string
}

func NewIgnores(names []IgnorePattern, paths []IgnorePattern) *Ignores {
	return &Ignores{
		names: names,
		paths: paths,
	}
}

func (is *Ignores) IgnoreName(name string) bool {
	if is == nil {
		return false
	}
	for _, p := range is.names {
		if p.ShouldIgnore(name) {
			return true
		}
	}
	return false
}

func (is *Ignores) IgnorePath(path string) bool {
	if is == nil {
		return false
	}
	if is.prefix != "" {
		path = strings.TrimLeft(path, is.prefix)
	}
	for _, p := range is.paths {
		if p.ShouldIgnore(path) {
			return true
		}
	}
	return false
}

func (is *Ignores) AddName(name IgnorePattern) {
	is.names = append(is.names, name)
}

func (is *Ignores) AddPath(path IgnorePattern) {
	is.paths = append(is.paths, path)
}

func (is *Ignores) SetPathPrefix(prefix string) {
	is.prefix = prefix
}

func (is *Ignores) Len() int {
	if is == nil {
		return 0
	}
	return len(is.names) + len(is.paths)
}

func (is *Ignores) Patterns() []IgnorePattern {
	if is == nil {
		return nil
	}
	var p []IgnorePattern
	p = append(p, is.names...)
	p = append(p, is.paths...)
	return p
}

type IgnorePattern interface {
	ShouldIgnore(path string) bool
}

type ExactPattern string

func (p ExactPattern) ShouldIgnore(path string) bool {
	return string(p) == path
}

// Deprecated it is not implemented yet
// NOTE: only * and ? is supported
// * matches any non empty sequence of non-separator character
// ? matches one non-separator character
// we are NOT expecting path to have * and ?
type WildcardPattern string

func (p WildcardPattern) ShouldIgnore(path string) bool {
	// the pattern would always be no greater than path due to our limited features
	if len(p) > len(path) {
		return false
	}
	s := -1 // * location
	i := 0
	j := 0
	for i < len(p) && j < len(path) {
		if p[i] == '*' {
			// abc.* abc.t, i = 4, j = 4
			if i == len(p)-1 {
				for ; j < len(path); j++ {
					if path[j] == filepath.Separator {
						return false
					}
				}
				return true
			}
			// ajax_*.html does not match ajax_/a.html
			if path[j] == filepath.Separator {
				return false
			}
			// abc*.html should match abc.ht.html
			s = i
			i++
			j++ // * should match at least one character
		} else if p[i] == '?' {
			if path[j] == filepath.Separator {
				return false
			}
			i++
			j++
		} else if p[i] == path[j] {
			i++
			j++
		} else if s != -1 && path[j] != filepath.Separator {
			i = s + 1
			j++
		} else {
			return false
		}
	}
	//log.Infof("len(p) %d i %d len(path) %d j %d", len(p), i, len(path), j)
	// both pattern and path reaches end
	return len(p) == i && len(path) == j
}
