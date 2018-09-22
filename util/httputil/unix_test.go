package httputil_test

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"

	dhttputil "github.com/dyweb/gommon/util/httputil"
	"github.com/dyweb/gommon/util/testutil"
)

func proxyDocker(prefix string) http.Handler {
	proxy := httputil.ReverseProxy{
		Transport: dhttputil.NewPooledUnixTransport("/var/run/docker.sock"),
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			r.URL.Host = "api"
			r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
			r.Host = "api"
		},
	}
	return &proxy
}

func TestNewPooledUnixTransport(t *testing.T) {
	t.Skip("onl/y runs on node with docker")

	t.Run("docker client", func(t *testing.T) {
		tr := dhttputil.NewPooledUnixTransport("/var/run/docker.sock")
		c := &http.Client{Transport: tr}
		t.Log(string(testutil.GetBody(t, c, "http://api/version")))
	})
	t.Run("docker proxy", func(t *testing.T) {
		mux := http.NewServeMux()
		//mux.Handle("/docker/proxy/", proxyDocker("/docker/proxy/")) 400 Bad Request
		mux.Handle("/docker/proxy/", proxyDocker("/docker/proxy"))
		srv := httptest.NewServer(mux)
		c := srv.Client()
		t.Log(string(testutil.GetBody(t, c, srv.URL+"/docker/proxy/version")))
		t.Log(string(testutil.GetBody(t, c, srv.URL+"/docker/proxy/info")))
	})
}
