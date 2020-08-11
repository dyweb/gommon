package httputil

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"time"
)

// http.go contains util for creating fresh transport and client that don't use default transport
// It is based on https://github.com/hashicorp/go-cleanhttp/blob/master/cleanhttp.go which is based
// on how default transport is created in standard library

// NewUnPooledTransport uses NewPooledTransport with keep alive and idle connection disabled
func NewUnPooledTransport() *http.Transport {
	tr := NewPooledTransport()
	tr.DisableKeepAlives = true
	tr.MaxIdleConnsPerHost = -1
	return tr
}

// NewPooledTransport is same as DefaultTransport in net/http.
// But it is not shared and won't be alerted by other library
func NewPooledTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxConnsPerHost:       runtime.GOMAXPROCS(0) + 1, // https://github.com/hashicorp/go-cleanhttp/blob/master/cleanhttp.go#L35
	}
}

// NewClient panics if transport is nil or http.DefaultTransport,
// You should always bring your own http.Transport instead of using the default one,
// because all the third party libraries can modify it without the application knowing
func NewClient(tr *http.Transport) *http.Client {
	if tr == nil {
		panic("transport is nil")
	}
	if tr == http.DefaultTransport {
		panic("stop using http.DefaultTransport")
	}
	return &http.Client{
		Transport: tr,
	}
}

// NewUnPooledClient returns a net/http client with a fresh http.Transport that has connection pooling disabled
func NewUnPooledClient() *http.Client {
	return NewClient(NewUnPooledTransport())
}

// NewPooledClient returns a net/http client with a fresh http.Transport using NewPooledTransport.
// The transport is a connection pool, you should reuse the client instead of creating new one
// with new transport, http.Client does not keep internal state in struct except cookie jar
//
// If you do need multiple clients to reuse same connections, you should use NewClient and pass a transport
func NewPooledClient() *http.Client {
	return NewClient(NewPooledTransport())
}

// DiscardBody drain and close the body and ignore all the errors,
// It should be used for test
func DiscardBody(res *http.Response) {
	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()
}
