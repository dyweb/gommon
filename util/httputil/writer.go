package httputil

import (
	"bufio"
	"net"
	"net/http"
)

// writer.go implements a http.ResponseWriter so it can be used to generate http access log
// it is copied from go.ice https://github.com/dyweb/go.ice/blob/v0.0.1/ice/transport/http/writer.go

var (
	_ http.ResponseWriter = (*TrackedWriter)(nil)
	_ http.Flusher        = (*TrackedWriter)(nil)
	_ http.Pusher         = (*TrackedWriter)(nil)
	_ http.Hijacker       = (*TrackedWriter)(nil)
)

// NOTE: CloseNotifier is deprecated in favor of context
//var _ http.CloseNotifier = (*TrackedWriter)(nil)

// TrackedWriter keeps track of status code and bytes written so it can be used by logger.
// It proxies all the interfaces except Hijacker, since it is not supported by HTTP/2.
// Most methods comments are copied from net/http
// It is based on https://github.com/gorilla/handlers but put all interface into one struct
type TrackedWriter struct {
	w           http.ResponseWriter
	status      int
	size        int
	writeCalled int
}

// NewTrackedWriter set the underlying writer based on argument,
// It returns a value instead of pointer so it can be allocated on stack.
// TODO: add benchmark to prove it ... might change it to pointer ...
func NewTrackedWriter(w http.ResponseWriter) TrackedWriter {
	return TrackedWriter{w: w, status: 200}
}

// Status return the tracked status code, returns 0 if WriteHeader has not been called
func (tracker *TrackedWriter) Status() int {
	return tracker.status
}

// Size return number of bytes written through Write, returns 0 if Write has not been called
func (tracker *TrackedWriter) Size() int {
	return tracker.size
}

// Header returns the header map of the underlying ResponseWriter
//
// Changing the header map after a call to WriteHeader (or
// Write) has no effect unless the modified headers are
// trailers.
func (tracker *TrackedWriter) Header() http.Header {
	return tracker.w.Header()
}

// Write keeps track of bytes written of the underlying ResponseWriter
//
// Write writes the data to the connection as part of an HTTP reply.
//
// If WriteHeader has not yet been called, Write calls
// WriteHeader(http.StatusOK) before writing the data. If the Header
// does not contain a Content-Type line, Write adds a Content-Type set
// to the result of passing the initial 512 bytes of written data to
// DetectContentType.
func (tracker *TrackedWriter) Write(b []byte) (int, error) {
	tracker.writeCalled++
	size, err := tracker.w.Write(b)
	tracker.size += size
	return size, err
}

// WriteHeader keep track of status code and call the underlying ResponseWriter
//
// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (tracker *TrackedWriter) WriteHeader(status int) {
	tracker.status = status
	tracker.w.WriteHeader(status)
}

// Flush calls Flush on underlying ResponseWriter if it implemented http.Flusher
//
// Flusher interface is implemented by ResponseWriters that allow
// an HTTP handler to flush buffered data to the client.
// The default HTTP/1.x and HTTP/2 ResponseWriter implementations
// support Flusher
func (tracker *TrackedWriter) Flush() {
	if f, ok := tracker.w.(http.Flusher); ok {
		f.Flush()
	}
}

// Push returns http.ErrNotSupported if underlying ResponseWriter does not implement http.Pusher
//
// Push initiates an HTTP/2 server push, returns ErrNotSupported if the client has disabled push or if push
// is not supported on the underlying connection.
func (tracker *TrackedWriter) Push(target string, opts *http.PushOptions) error {
	if p, ok := tracker.w.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return http.ErrNotSupported
}

// Hijack implements http.Hijacker, which is used by websocket etc.
// It returns http.ErrNotSupported with nil pointer if the underlying writer does not support it
// NOTE: HTTP/1.x supports it but HTTP/2 does not
func (tracker *TrackedWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := tracker.w.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}
