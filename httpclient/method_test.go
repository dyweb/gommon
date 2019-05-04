package httpclient_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/httpclient"
	"github.com/dyweb/gommon/util/httputil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make sure we don't put wrong method due to copy and paste

func TestMethod(t *testing.T) {
	var m httputil.Method
	mux := http.NewServeMux()
	mux.HandleFunc("/method", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, string(m), r.Method)
		errors.Ignore2(w.Write([]byte(`{"foo": "bar"}`)))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	client, err := httpclient.New(srv.URL, httpclient.UseJSON())
	require.Nil(t, err)

	// well ... this test is much longer than the code being tested ...
	cases := []struct {
		noReqBody func(ctx *httpclient.Context, path string, resBody interface{}) error
		raw       func(ctx *httpclient.Context, path string) (*http.Response, error)
		ignoreRes func(ctx *httpclient.Context, path string) error
		m         httputil.Method
	}{
		{
			noReqBody: client.Get,
			raw:       client.GetRaw,
			ignoreRes: client.GetIgnoreRes,
			m:         httputil.Get,
		},
		{
			noReqBody: client.Delete,
			raw:       client.DeleteRaw,
			ignoreRes: client.DeleteIgnoreRes,
			m:         httputil.Delete,
		},
		{
			noReqBody: func(ctx *httpclient.Context, path string, resBody interface{}) error {
				return client.Post(ctx, path, nil, resBody)
			},
			raw: func(ctx *httpclient.Context, path string) (response *http.Response, e error) {
				return client.PostRaw(ctx, path, nil)
			},
			ignoreRes: func(ctx *httpclient.Context, path string) error {
				return client.PostIgnoreRes(ctx, path, nil)
			},
			m: httputil.Post,
		},
		{
			noReqBody: func(ctx *httpclient.Context, path string, resBody interface{}) error {
				return client.Patch(ctx, path, nil, resBody)
			},
			raw: func(ctx *httpclient.Context, path string) (response *http.Response, e error) {
				return client.PatchRaw(ctx, path, nil)
			},
			ignoreRes: func(ctx *httpclient.Context, path string) error {
				return client.PatchIgnoreRes(ctx, path, nil)
			},
			m: httputil.Patch,
		},
	}
	for _, cs := range cases {
		m = cs.m
		d := make(map[string]string)
		assert.Nil(t, cs.noReqBody(httpclient.Bkg(), "/method", &d))
		_, err := cs.raw(httpclient.Bkg(), "/method")
		assert.NoError(t, err)
		assert.Nil(t, cs.ignoreRes(httpclient.Bkg(), "/method"))
	}
}
