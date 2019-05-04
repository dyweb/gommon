package httpclient_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/httpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContext_SetBase(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, map[string]string{
			"version": "1.0",
		})
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, map[string]string{
			"health": "green",
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	client, err := httpclient.New(srv.URL+"/api/", httpclient.UseJSON())
	require.Nil(t, err)

	ver := make(map[string]string)
	require.Nil(t, client.Get(httpclient.Bkg(), "version", &ver))
	assert.Equal(t, ver["version"], "1.0")
	require.NotNil(t, client.Get(httpclient.Bkg(), "health", &ver)) // api/health does not exist
	// override base in context
	ctx := httpclient.Bkg().SetBase(srv.URL)
	health := make(map[string]string)
	require.Nil(t, client.Get(ctx, "health", &health))
	assert.Equal(t, health["health"], "green")
}

func TestContext_SetParam(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/param", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		dump := make(map[string][]string)
		for key, v := range q {
			dump[key] = v
		}
		writeJSON(t, w, dump)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	client, err := httpclient.New(srv.URL, httpclient.UseJSON())
	require.Nil(t, err)

	ctx := httpclient.Bkg().SetParam("foo", "bar")
	dump := make(map[string][]string)
	assert.Nil(t, client.Get(ctx, "/param", &dump))
	assert.Equal(t, "bar", dump["foo"][0])
}

func TestContext_SetErrorHandler(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/404/nobody", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	client, err := httpclient.New(srv.URL, httpclient.UseJSON())
	require.Nil(t, err)

	h := httpclient.ErrorHandlerFunc(func(status int, body []byte, res *http.Response) error {
		return errors.Errorf("custom %d", status)
	})
	ctx := httpclient.Bkg().SetErrorHandler(h)
	res, err := client.GetRaw(ctx, "/404/nobody")
	assert.NotNil(t, res, "application error has no nil response")
	assert.Equal(t, "custom 404", err.Error())
}
