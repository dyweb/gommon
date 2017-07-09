package config

import (
	"github.com/dyweb/gommon/util"
)

var log = util.Logger.NewEntryWithPkg("gommon.config")

const (
	yamlDocumentSeparator = "---"
	defaultTemplateName   = "gommon yaml"
	defaultKeyDelimiter   = "."
)
