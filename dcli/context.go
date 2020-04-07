package dcli

import (
	"context"
	"time"
)

// Context implements context.Context and provides cli specific helper func.
// A default implementation DefaultContext is provided.
type Context interface {
	context.Context
}

var _ Context = (*DefaultContext)(nil)

type DefaultContext struct {
	stdCtx context.Context
}

// TODO(generator): those default context wrapper should be generated. It is also used in httpclient package
// Deadline returns Deadline() from underlying context.Context if set
func (c *DefaultContext) Deadline() (deadline time.Time, ok bool) {
	if c != nil && c.stdCtx != nil {
		return c.stdCtx.Deadline()
	}
	// NOTE: we are using named return, so empty value will be returned
	// learned this from context.Context's emptyCtx implementation
	return
}

// Done returns Done() from underlying context.Context if set
func (c *DefaultContext) Done() <-chan struct{} {
	if c != nil && c.stdCtx != nil {
		return c.stdCtx.Done()
	}
	// Done may return nil if this context can never be canceled
	return nil
}

// Err returns Err() from underlying context.Context if set
func (c *DefaultContext) Err() error {
	if c != nil && c.stdCtx != nil {
		return c.stdCtx.Err()
	}
	return nil
}

// Value first checks the map[string]interface{},
// if not found, it use the underlying context.Context if is set
// if not set, it returns nil
func (c *DefaultContext) Value(key interface{}) interface{} {
	if c != nil && c.stdCtx != nil {
		return c.stdCtx.Value(key)
	}
	return nil
}
