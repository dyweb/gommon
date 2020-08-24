// Package httputil provides helper for net/http, i.e. unix domain socket client, http request logger
package httputil

import (
	"net/http"
)

// Method gives a type for http method
// See https://github.com/bradfitz/exp-httpclient/blob/master/http/method.go
type Method string

const (
	Get     Method = http.MethodGet
	Head    Method = http.MethodHead
	Post    Method = http.MethodPost
	Put     Method = http.MethodPut
	Patch   Method = http.MethodPatch
	Delete  Method = http.MethodDelete
	Connect Method = http.MethodConnect
	Options Method = http.MethodOptions
	Trace   Method = http.MethodTrace
)

// AllMethods returns all the http methods (defined in this package)
func AllMethods() []Method {
	return []Method{
		Get,
		Head,
		Post,
		Put,
		Patch,
		Delete,
		Connect,
		Options,
		Trace,
	}
}

// UserAgent data are from https://techblog.willshouse.com/2012/01/03/most-common-user-agents/
// For UserAgent spec, see MDN https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent
type UserAgent string

const (
	DefaultUA UserAgent = "Go" // TODO: should put real version number into the constant, need to use generator
	// TODO: should list more UAs, including mobile devices https://deviceatlas.com/blog/list-of-user-agent-strings
	// TODO: os and browser versions are too old in those UAs ...
	UACurl        UserAgent = "curl/7.60.0"
	UAChromeWin   UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"
	UAChromeLinux UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.62 Safari/537.36"
	UAChromeMac   UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"
	// TODO: add UA for mobile device
)
