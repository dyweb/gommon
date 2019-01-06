package httputil_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/httputil"
	"github.com/stretchr/testify/assert"
)

func TestNewTrackedWriter(t *testing.T) {
	if testing.Short() {
		t.Skip("http test server is used")
	}

	mux := http.NewServeMux()
	response := "bar"
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		errors.Ignore2(
			w.Write([]byte(response)),
		)
	})
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tr := httputil.NewTrackedWriter(w)
		mux.ServeHTTP(&tr, r)
		// the size of the response body is tracked
		assert.Equal(t, len(response), tr.Size())
	})
	srv := httptest.NewServer(h)
	defer srv.Close()

	// send a request
	res, err := http.Get(srv.URL + "/foo")
	assert.Nil(t, err)
	httputil.DiscardBody(res)
}
