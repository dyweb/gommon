// Package requests wrap net/http like requests did for python
// it is easy to use, but not very efficient
package requests

import (
	"io/ioutil"
	"net/http"

	"bytes"
	"github.com/pkg/errors"
	"io"
	"strings"
)

const (
	ShadowSocksLocal = "127.0.0.1:1080"
	ContentJSON      = "application/json"
)

type Client struct {
	http *http.Client
}

func NewClient(options ...func(h *http.Client)) (*Client, error) {
	c := &Client{http: NewDefaultClient()}
	for _, option := range options {
		option(c.http)
	}
	return c, nil
}

func (client *Client) makeRequest(method string, url string, body io.Reader) (*Response, error) {
	if client.http == nil {
		return nil, errors.New("http client is not initialized")
	}
	var res *http.Response
	var err error
	switch method {
	case http.MethodGet:
		res, err = client.http.Get(url)
	case http.MethodPost:
		// TODO: we only support JSON for now
		res, err = client.http.Post(url, ContentJSON, body)
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

func (client *Client) Get(url string) (*Response, error) {
	return client.makeRequest(http.MethodGet, url, nil)
}

// TODO: accept io.reader

func (client *Client) PostJSONString(url string, data string) (*Response, error) {
	return client.makeRequest(http.MethodPost, url, ioutil.NopCloser(strings.NewReader(data)))
}

func (client *Client) PostJSONBytes(url string, data []byte) (*Response, error) {
	return client.makeRequest(http.MethodPost, url, ioutil.NopCloser(bytes.NewReader(data)))
}

func (client *Client) GetJSON(url string, data interface{}) error {
	res, err := client.Get(url)
	if err != nil {
		return errors.Wrap(err, "error getting response")
	}
	err = res.JSON(data)
	if err != nil {
		return errors.Wrap(err, "error parsing json")
	}
	return nil
}

func (client *Client) GetJSONStringMap(url string) (map[string]string, error) {
	var data map[string]string
	res, err := client.Get(url)
	if err != nil {
		return data, errors.Wrap(err, "error getting response")
	}
	data, err = res.JSONStringMap()
	if err != nil {
		return data, errors.Wrap(err, "error parsing json")
	}
	return data, nil
}
