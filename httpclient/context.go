package httpclient

import (
	"context"
	"time"
)

var _ context.Context = (*Context)(nil)

// Context implements context.Context and provides HTTP request specific helpers.
// It is lazy initialized, only call `make` when there is write to internal maps.
// User (including this package itself) should use setter when set value.
type Context struct {
	// base overrides base path set in client if it is not empty
	base string
	// headers is request specific headers, headers configured in client will be override
	headers map[string]string
	// params is the query parameters attached to url, i.e. query?name=foo&type=bar
	params map[string]string
	// errHandler is per request error handler, it will overwrite client level error handler
	// if it has non nil value
	errHandler ErrorHandler

	// wraps a standard context implementation for value and deadline
	stdCtx context.Context
}

// Bkg returns a Context that does not embed context.Context,
// it behaves like context.Background(), however we can't use a singleton like context package
// because Context in httpclient can be modified in place to store req/response body etc.
// So we always return pointer to a fresh new instance.
//
// We return pointer because it is meant to be modified along the way, it is not immutable like context.Context
// Also the context.Context is implemented using pointer receiver
func Bkg() *Context {
	return &Context{}
}

// NewContext returns a context that embed a context.Context
func NewContext(ctx context.Context) *Context {
	return &Context{
		stdCtx: ctx,
	}
}

// ConvertContext returns the original context if it is already *httpclient.Context,
// Otherwise it wraps the context using NewContext
func ConvertContext(ctx context.Context) *Context {
	c, ok := ctx.(*Context)
	if ok {
		return c
	}
	return NewContext(ctx)
}

// SetBase allow a single request to override client level request base.
// This is useful when most request is /api/bla and suddenly there is a /bla/api.
func (c *Context) SetBase(s string) *Context {
	c.base = s
	return c
}

func (c *Context) SetHeader(k, v string) *Context {
	if c.headers == nil {
		c.headers = make(map[string]string)
	}
	c.headers[k] = v
	return c
}

func (c *Context) SetParam(k, v string) *Context {
	if c.params == nil {
		c.params = make(map[string]string)
	}
	c.params[k] = v
	return c
}

func (c *Context) SetErrorHandler(h ErrorHandler) *Context {
	c.errHandler = h
	return c
}

// Deadline returns Deadline() from underlying context.Context if set
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	if c != nil && c.stdCtx != nil {
		return c.stdCtx.Deadline()
	}
	// NOTE: we are using named return, so empty value will be returned
	// learned this from context.Context's emptyCtx implementation
	return
}

// Done returns Done() from underlying context.Context if set
func (c *Context) Done() <-chan struct{} {
	if c != nil && c.stdCtx != nil {
		return c.stdCtx.Done()
	}
	// Done may return nil if this context can never be canceled
	return nil
}

// Err returns Err() from underlying context.Context if set
func (c *Context) Err() error {
	if c != nil && c.stdCtx != nil {
		return c.stdCtx.Err()
	}
	return nil
}

// Value first checks the map[string]interface{},
// if not found, it use the underlying context.Context if is set
// if not set, it returns nil
func (c *Context) Value(key interface{}) interface{} {
	if c != nil && c.stdCtx != nil {
		return c.stdCtx.Value(key)
	}
	return nil
}
