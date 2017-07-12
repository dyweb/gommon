package config

import (
	"github.com/dyweb/gommon/util"
	dlog "github.com/dyweb/gommon/log"
)

var log = util.Logger.NewEntryWithPkg("gommon.config")

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
