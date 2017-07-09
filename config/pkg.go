package config

import (
	"github.com/dyweb/gommon/util"
)

var log = util.Logger.NewEntryWithPkg("gommon.config")

const (
	yamlDocumentSeparator = "---"
	pongo2DefaultBaseDir  = ""
	pongo2DefaultSetName  = "gommon-yaml"
	defaultKeyDelimiter   = "."
)
