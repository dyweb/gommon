package httpclient_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dyweb/gommon/httpclient"
	"github.com/dyweb/gommon/util/httputil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultHandler(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/404/nobody", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	mux.HandleFunc("/404/body", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("nginx 0.0.1 404"))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	client, err := httpclient.New(srv.URL, httpclient.UseJSON())
	require.Nil(t, err)

	res, err := client.GetRaw(httpclient.Bkg(), "/404/nobody")
	assert.NotNil(t, res, "application error has no nil response")
	appErr, ok := err.(*httpclient.ErrApplication)
	require.True(t, ok, "default error handler gives ErrApplication")
	assert.Equal(t, 404, appErr.Status)
	assert.Equal(t, httputil.Get, appErr.Method)
	assert.Equal(t, "/404/nobody", appErr.Path)

	res, err = client.GetRaw(httpclient.Bkg(), "/404/body")
	appErr, ok = err.(*httpclient.ErrApplication)
	require.True(t, ok, "default error handler gives ErrApplication")
	assert.Equal(t, "nginx 0.0.1 404", appErr.Body)
}
