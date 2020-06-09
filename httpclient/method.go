package httpclient

import (
	"net/http"

	"github.com/dyweb/gommon/util/httputil"
)

// method.go contains wrapper for common http verbs, GET, POST, PUT, PATCH, DELETE

// GET

func (c *Client) Get(ctx *Context, path string, resBody interface{}) error {
	return c.FetchTo(ctx, httputil.Get, path, nil, resBody)
}

func (c *Client) GetRaw(ctx *Context, path string) (*http.Response, error) {
	return c.Do(ctx, httputil.Get, path, nil)
}

func (c *Client) GetIgnoreRes(ctx *Context, path string) error {
	return c.FetchToNull(ctx, httputil.Get, path, nil)
}

// POST

func (c *Client) Post(ctx *Context, path string, reqBody interface{}, resBody interface{}) error {
	return c.FetchTo(ctx, httputil.Post, path, reqBody, resBody)
}

func (c *Client) PostRaw(ctx *Context, path string, reqBody interface{}) (*http.Response, error) {
	return c.Do(ctx, httputil.Post, path, reqBody)
}

func (c *Client) PostIgnoreRes(ctx *Context, path string, reqBody interface{}) error {
	return c.FetchToNull(ctx, httputil.Post, path, reqBody)
}

// PUT

func (c *Client) Put(ctx *Context, path string, reqBody interface{}, resBody interface{}) error {
	return c.FetchTo(ctx, httputil.Put, path, reqBody, resBody)
}

func (c *Client) PutRaw(ctx *Context, path string, reqBody interface{}) (*http.Response, error) {
	return c.Do(ctx, httputil.Put, path, reqBody)
}

func (c *Client) PutIgnoreRes(ctx *Context, path string, reqBody interface{}) error {
	return c.FetchToNull(ctx, httputil.Put, path, reqBody)
}

// PATCH

func (c *Client) Patch(ctx *Context, path string, reqBody interface{}, resBody interface{}) error {
	return c.FetchTo(ctx, httputil.Patch, path, reqBody, resBody)
}

func (c *Client) PatchRaw(ctx *Context, path string, reqBody interface{}) (*http.Response, error) {
	return c.Do(ctx, httputil.Patch, path, reqBody)
}

func (c *Client) PatchIgnoreRes(ctx *Context, path string, reqBody interface{}) error {
	return c.FetchToNull(ctx, httputil.Patch, path, reqBody)
}

// DELETE

func (c *Client) Delete(ctx *Context, path string, resBody interface{}) error {
	return c.FetchTo(ctx, httputil.Delete, path, nil, resBody)
}

func (c *Client) DeleteRaw(ctx *Context, path string) (*http.Response, error) {
	return c.Do(ctx, httputil.Delete, path, nil)
}

func (c *Client) DeleteIgnoreRes(ctx *Context, path string) error {
	return c.FetchToNull(ctx, httputil.Delete, path, nil)
}
