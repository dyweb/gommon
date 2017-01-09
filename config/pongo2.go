package config

import (
	"github.com/flosch/pongo2"
	"github.com/pkg/errors"
)

var (
	// TODO: need to implement the template loader to have better controlled template
	// include and import logic
	defaultLoader = pongo2.MustNewLocalFileSystemLoader("")
	defaultSet    = pongo2.NewSet("gommon", defaultLoader)
)

func init() {

}

func RenderDocument(str string) (string, error) {
	//pongo2.Context{} is just map[string]interface{}
	//FIXME: pongo2.FromString is not longer in the new API, must first create a set
	tpl, err := defaultSet.FromString(str)
	if err != nil {
		return "", errors.Wrap(err, "can't parse template")
	}
	out, err := tpl.Execute(pongo2.Context{"foo1": "bar", "foo2": 1})
	if err != nil {
		return "", errors.Wrap(err, "can'r render template")
	}
	return out, nil
}
