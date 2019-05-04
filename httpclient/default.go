package httpclient

import (
	"crypto/tls"
	"net/http"
	"sync"

	"github.com/dyweb/gommon/util/httputil"
)

var (
	pkgMu                      sync.Mutex
	defaultTransport           *http.Transport
	defaultTransportSkipVerify *http.Transport
)

// NewDefault creates a dev client that use a shared un pooled http.Transport
func NewDefault(base string) (*Client, error) {
	pkgMu.Lock()
	defer pkgMu.Unlock()

	if defaultTransport == nil {
		defaultTransport = httputil.NewUnPooledTransport()
	}
	return New(base, WithTransport(defaultTransport))
}

// NewDefaultSkipVerify is same as NewDefault but disabled tls certificate verify.
// You should NOT use it production, it is useful for testing with self signed certs
// TODO: add links on how to generate and use self signed cert
func NewDefaultSkipVerify(base string) (*Client, error) {
	pkgMu.Lock()
	defer pkgMu.Unlock()

	if defaultTransportSkipVerify == nil {
		defaultTransportSkipVerify = httputil.NewUnPooledTransport()
		defaultTransportSkipVerify.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	return New(base, WithTransport(defaultTransportSkipVerify))
}
