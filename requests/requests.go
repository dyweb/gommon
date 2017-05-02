// Package requests wrap net/http like requests did for python
// it is easy to use, but not very efficient
package requests

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"bytes"
	"io"
	"strings"
)

const (
	ContentJSON = "application/json"
)

func makeRequest(method string, url string, body io.Reader) (*Response, error) {
	var res *http.Response
	var err error
	switch method {
	case http.MethodGet:
		res, err = http.Get(url)
	case http.MethodPost:
		// TODO: we only support JSON for now
		res, err = http.Post(url, ContentJSON, body)
	}
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

func Get(url string) (*Response, error) {
	return makeRequest(http.MethodGet, url, nil)
}

func PostJSONString(url string, data string) (*Response, error) {
	return makeRequest(http.MethodPost, url, ioutil.NopCloser(strings.NewReader(data)))
}

func PostJSONBytes(url string, data []byte) (*Response, error) {
	return makeRequest(http.MethodPost, url, ioutil.NopCloser(bytes.NewReader(data)))
}

func GetJSON(url string, data interface{}) error {
	res, err := Get(url)
	if err != nil {
		return errors.Wrap(err, "error getting response")
	}
	err = res.JSON(data)
	if err != nil {
		return errors.Wrap(err, "error parsing json")
	}
	return nil
}

func GetJSONStringMap(url string) (map[string]string, error) {
	var data map[string]string
	res, err := Get(url)
	if err != nil {
		return data, errors.Wrap(err, "error getting response")
	}
	data, err = res.JSONStringMap()
	if err != nil {
		return data, errors.Wrap(err, "error parsing json")
	}
	return data, nil
}
