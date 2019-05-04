package httpclient

import (
	"fmt"
	"net/http"

	"github.com/dyweb/gommon/util/httputil"
)

type ErrorHandler interface {
	IsError(status int, res *http.Response) bool
	DecodeError(status int, body []byte, res *http.Response) (decodedError error)
}

// ErrorHandlerFunc use status code for IsError
type ErrorHandlerFunc func(status int, body []byte, res *http.Response) (decodedError error)

func IsSuccess(status int) bool {
	return status >= 200 && status <= 299
}

func (f ErrorHandlerFunc) IsError(status int, res *http.Response) bool {
	return !IsSuccess(status)
}

func (f ErrorHandlerFunc) DecodeError(status int, body []byte, res *http.Response) error {
	if f != nil {
		return f(status, body, res)
	}
	return defaultHandler.DecodeError(status, body, res)
}

var defaultHandler ErrorHandler = BasicErrorHandler{}

func DefaultHandler() ErrorHandler {
	return defaultHandler
}

type ErrApplication struct {
	Status int
	Method httputil.Method
	// Url is the full url including schema, host and query parameters
	Url string
	// Path is just the api path without host etc.
	Path string
	// Body is string instead of []byte because it's immutable
	Body string
}

func (e *ErrApplication) Error() string {
	return fmt.Sprintf("%d %s %s %s", e.Status, e.Method, e.Path, e.Body)
}

type ErrDecoding struct {
	Codec string
	Err   error
	Body  string
}

func (e *ErrDecoding) Error() string {
	return fmt.Sprintf("decode %s got %s body %s", e.Codec, e.Err, e.Body)
}

// BasicErrorHandler use value receiver because it does not have any fields
type BasicErrorHandler struct {
}

func (h BasicErrorHandler) IsError(status int, res *http.Response) bool {
	return !IsSuccess(status)
}

func (h BasicErrorHandler) DecodeError(status int, body []byte, res *http.Response) error {
	return &ErrApplication{
		Status: res.StatusCode,
		Method: httputil.Method(res.Request.Method),
		Url:    res.Request.URL.String(),
		Path:   res.Request.URL.Path,
		Body:   string(body),
	}
}
