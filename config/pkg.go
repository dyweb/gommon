package config

import (
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
