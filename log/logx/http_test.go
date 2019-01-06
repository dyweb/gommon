package logx_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	dlog "github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/log/logx"
	"github.com/dyweb/gommon/util/httputil"
	"github.com/stretchr/testify/assert"
)

func TestNewHttpAccessLogger(t *testing.T) {
	if testing.Short() {
		t.Skip("http test server is used")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bar"))
	})
	logger := dlog.NewPackageLogger()
	var buf bytes.Buffer
	logger.SetHandler(dlog.NewTextHandler(&buf))
	h := logx.NewHttpAccessLogger(logger, mux)
	srv := httptest.NewServer(h)
	defer srv.Close()

	res, err := http.Get(srv.URL + "/foo")
	assert.Nil(t, err)
	httputil.DiscardBody(res)

	// TODO: need better matching ...
	assert.Contains(t, buf.String(), "/foo status=200 size=3 method=GET duration=")
	t.Log(buf.String())
}
