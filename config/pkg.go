package config

import (
	dlog "github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/util"
)

var log = util.Logger.RegisterPkg()

const (
	yamlDocumentSeparator = "---"
	defaultTemplateName   = "gommon yaml"
	defaultKeyDelimiter   = "."
)

type StructuredConfig interface {
	Validate() error
}

// NOTE: the interface check is here to avoid import cycle
var _ StructuredConfig = (*dlog.Config)(nil)
