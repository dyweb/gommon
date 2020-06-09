package httpclient

import (
	"net/http"

	"github.com/dyweb/gommon/errors"
)

type Option func(c *Client) error

func UseJSON() Option {
	return func(c *Client) error {
		c.json = true
		c.SetHeader("Content-Type", "application/json").SetHeader("Accept", "application/json")
		return nil
	}
}

func WithErrorHandler(h ErrorHandler) Option {
	return func(c *Client) error {
		c.errHandler = h
		return nil
	}
}

func WithErrorHandlerFunc(f ErrorHandlerFunc) Option {
	return func(c *Client) error {
		c.errHandler = f
		return nil
	}
}
func WithClient(h *http.Client) Option {
	return func(c *Client) error {
		c.h = h
		return nil
	}
}

func WithTransport(tr *http.Transport) Option {
	return func(c *Client) error {
		if c.h == nil {
			return errors.New("native http client is empty, can't set transport")
		}
		c.h.Transport = tr
		return nil
	}
}

// TODO: add skip verify

func applyOptions(c *Client, opts ...Option) error {
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return errors.Wrap(err, "error apply option")
		}
	}
	return nil
}
