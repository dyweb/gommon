package config

import (
	"github.com/flosch/pongo2"
	"github.com/pkg/errors"
)

var (
	// TODO: may need to implement the template loader to have better controlled over template include and import logic
	defaultLoader = pongo2.MustNewLocalFileSystemLoader("")
	defaultSet    = pongo2.NewSet("gommon", defaultLoader)
)

// RenderDocumentString uses defaultSet due to pongo2's strange API
func RenderDocumentString(tplStr string, context pongo2.Context) (string, error) {
	//pongo2.Context{} is just map[string]interface{}
	//FIXME: pongo2.FromString is not longer in the new API, must first create a set
	tpl, err := defaultSet.FromString(tplStr)
	if err != nil {
		return "", errors.Wrap(err, "can't parse template")
	}
	out, err := tpl.Execute(context)
	if err != nil {
		return "", errors.Wrap(err, "can'r render template")
	}
	return out, nil
}

// RenderDocumentBytes uses defaultSet and since pongo2 have two function for String and Bytes, the wrapper also has two function
func RenderDocumentBytes(tplBytes []byte, context pongo2.Context) ([]byte, error) {
	tpl, err := defaultSet.FromBytes(tplBytes)
	var out []byte
	if err != nil {
		return out, errors.Wrap(err, "can't parse template")
	}
	out, err = tpl.ExecuteBytes(context)
	if err != nil {
		return out, errors.Wrap(err, "can'r render template")
	}
	return out, nil
}
