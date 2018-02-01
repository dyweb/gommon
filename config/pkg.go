package config

import (
	dlog "github.com/dyweb/gommon/legacy/log"
	"github.com/dyweb/gommon/util/logutil"
)

var log = logutil.NewPackageLogger()

const (
	yamlDocumentSeparator = "---"
	defaultTemplateName   = "gommon yaml"
	defaultKeyDelimiter   = "."
)

type Path string

type Reader interface {
	Path() Path
	Content() string
}

type StructuredConfig interface {
	Validate() error
}

// FIXME: migrate the structured config for logger
// NOTE: the interface check is here to avoid import cycle
var _ StructuredConfig = (*dlog.Config)(nil)
