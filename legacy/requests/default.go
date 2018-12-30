// +build ignore

package requests

import (
	"net"
	"net/http"
	"time"
)

// Default Transport Client that is same as https://golang.org/src/net/http/transport.go
// It's similar to https://github.com/hashicorp/go-cleanhttp

// NewDefaultTransport is copied from net/http/transport.go
func NewDefaultTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// NewDefaultClient returns a client using NewDefaultTransport
func NewDefaultClient() *http.Client {
	return &http.Client{
		Transport: NewDefaultTransport(),
	}
}
