// Package httpclient is a high level wrapper around net/http with more types and easier to use interface
// It is loosely modeled after https://github.com/bradfitz/exp-httpclient
package httpclient

// UnixBasePath is used as a placeholder for unix domain socket client.
// Only the protocol http is needed, host can be anything because dialer use socket path
// TODO: what about https over unix domain socket?
const UnixBasePath = "http://localhost"
