// Package requests wrap net/http like requests did for python
// it is easy to use, but not very efficient
package requests

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func Get(url string) (*Response, error) {
	res, err := http.Get(url)
	response := &Response{}
	if err != nil {
		return response, errors.Wrap(err, "error making request")
	}
	defer res.Body.Close()
	resContent, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, errors.Wrap(err, "error reading body")
	}
	response.Res = res
	response.Text = resContent
	return response, nil
}

func GetJSON(url string) (map[string]string, error) {
	var data map[string]string
	res, err := Get(url)
	if err != nil {
		return data, errors.Wrap(err, "error getting response")
	}
	data, err = res.JSON()
	if err != nil {
		return data, errors.Wrap(err, "error parsing json")
	}
	return data, nil
}
