package requests

import (
	"crypto/tls"
	"net/http"

	"golang.org/x/net/proxy"

	"github.com/dyweb/gommon/errors"
)

// TODO: might switch to functional options https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
// TODO: http proxy? might just use environment variable? https://stackoverflow.com/questions/14661511/setting-up-proxy-for-http-client
// TransportBuilder is the initial builder, its method use value receiver and return a new copy for chaining and keep itself unchanged
var TransportBuilder transportBuilder

type transportBuilder struct {
	skipKeyVerify bool
	socks5Host    string
	auth          *proxy.Auth
}

func (b transportBuilder) SkipKeyVerify() transportBuilder {
	b.skipKeyVerify = true
	return b
}

// UseShadowSocks uses the default config for shadowsocks local
func (b transportBuilder) UseShadowSocks() transportBuilder {
	return b.UseSocks5(ShadowSocksLocal, "", "")
}

func (b transportBuilder) UseSocks5(host string, username string, password string) transportBuilder {
	b.socks5Host = host
	if username == "" && password == "" {
		b.auth = nil
	} else {
		b.auth.User = username
		b.auth.Password = password
	}
	return b
}

func (b transportBuilder) Build() (*http.Transport, error) {
	transport := NewDefaultTransport()
	if b.skipKeyVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	if b.socks5Host != "" {
		// NOTE: actually the implementation does not connect at all, the err is always nil
		dialer, err := proxy.SOCKS5("tcp", b.socks5Host, b.auth, proxy.Direct)
		if err != nil {
			return nil, errors.Wrapf(err, "can't create socks5 dialer to %s with %s:%s", b.socks5Host, b.auth.User, b.auth.Password)
		}
		// TODO: Dial is deprecated and we should use DialContext, maybe contribute to golang/x/net
		transport.Dial = dialer.Dial
	}
	return transport, nil
}

func init() {
	TransportBuilder = transportBuilder{
		skipKeyVerify: false,
		socks5Host:    "",
		auth:          &proxy.Auth{},
	}
}
