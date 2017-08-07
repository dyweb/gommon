package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// Fields is key-value string pair to annotate the log and can be used by filter
type Fields map[string]string

// Logger is used to set output, formatter and filters, the real log operation is in Entry
type Logger struct {
	Out            io.Writer
	Formatter      Formatter
	Level          Level
	mu             sync.Mutex // FIXME: this mutex is never used, I guess I was following logrus when I wrote this
	showSourceLine bool
	Filters        map[Level]map[string]Filter
	Entries        map[string]*Entry
}

// NewLogger returns a new logger using StdOut and InfoLevel
func NewLogger() *Logger {
	f := make(map[Level]map[string]Filter, len(AllLevels))
	for _, level := range AllLevels {
		f[level] = make(map[string]Filter, 1)
	}
	l := &Logger{
		Out:            os.Stdout,
		Formatter:      NewTextFormatter(),
		Level:          InfoLevel,
		Filters:        f,
		Entries:        make(map[string]*Entry),
		showSourceLine: false,
	}
	return l
}

func (log *Logger) SetLevel(s string) error {
	newLevel, err := ParseLevel(s, false)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("can't set logging level to %s", s))
	}
	if log.Level == newLevel {
		return nil
	}
	// update all the registered entries
	for _, entry := range log.Entries {
		// only update entry's level if the new global level is more verbose
		if newLevel > entry.EntryLevel {
			entry.EntryLevel = newLevel
		}
	}
	log.Level = newLevel
	return nil
}

// EnableSourceLine add `source` field when logging, it use runtime.Caller(), the overhead has not been measured
func (log *Logger) EnableSourceLine() {
	log.showSourceLine = true
}

// DisableSourceLine does not show `source` field
func (log *Logger) DisableSourceLine() {
	log.showSourceLine = false
}

// AddFilter add a filter to logger, the filter should be simple string check on fields, i.e. PkgFilter check pkg field
func (log *Logger) AddFilter(filter Filter, level Level) {
	log.Filters[level][filter.FilterName()] = filter
}

// NewEntry returns an Entry with empty Fields
// Deprecated: use RegisterPkg instead
func (log *Logger) NewEntry() *Entry {
	// TODO: may use pool, but need benchmark to see if using pool provides improvement
	return &Entry{
		Logger:     log,
		Pkg:        "",
		EntryLevel: log.Level,
		Fields:     make(map[string]string, 1),
	}
}

// NewEntryWithPkg returns an Entry with pkg Field set to pkgName, should be used with PkgFilter
// Deprecated: use RegisterPkg instead
func (log *Logger) NewEntryWithPkg(pkgName string) *Entry {
	fields := make(map[string]string, 1)
	fields["pkg"] = pkgName
	e := &Entry{
		Logger:     log,
		Pkg:        pkgName,
		EntryLevel: log.Level,
		Fields:     fields,
	}
	log.Entries[pkgName] = e
	return e
}

// TODO: allow different level for each entry
// TODO: this is better than do filter in logger since we can apply the logging to each entry
func (log *Logger) RegisterPkg() *Entry {
	fields := make(map[string]string, 1)
	pkg := getCallerPackage(2)
	fields["pkg"] = pkg
	e := &Entry{
		Logger:     log,
		Pkg:        pkg,
		EntryLevel: log.Level,
		Fields:     fields,
	}
	log.Entries[pkg] = e
	return e
}

func (log *Logger) RegisteredPkgs() map[string]*Entry {
	return log.Entries
}

func (log *Logger) PrintEntries() {
	// TODO: sort and print it in a hierarchy
	// github.com/dyweb/Ayi/web
	// github.com/dyweb/Ayi/web/static
	for pkg := range log.Entries {
		fmt.Println(pkg)
	}
}

// FIXME: it should be in util package, but we put it here to avoid import cycle
func getCallerPackage(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	fnName := fn.Name()
	lastDot := strings.LastIndex(fnName, ".")
	return fnName[:lastDot]
}
