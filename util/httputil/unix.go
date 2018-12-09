package httputil

import (
	"context"
	"net"
	"net/http"
	"time"
)

func NewPooledUnixTransport(sockFile string) *http.Transport {
	return &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			d := net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}
			return d.DialContext(ctx, "unix", sockFile)
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

func ListenAndServeUnix(srv *http.Server, addr string) error {
	ln, err := net.Listen("unix", addr)
	if err != nil {
		return err
	}
	// TODO: do we need tcpKeepAliveListener like the default tcp ListenAndServe?
	return srv.Serve(ln)
}
