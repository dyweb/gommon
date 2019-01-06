package logx

import (
	"net/http"
	"time"

	dlog "github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/util/httputil"
)

// http.go provides handler for log http access

// HttpAccessLogger is an example for how to use gommon/log for http access log
type HttpAccessLogger struct {
	h      http.Handler // underlying http handler, the one do the actual job, normally a mux for dispatching
	logger *dlog.Logger
}

// NewHttpAccessLogger returns a http.Handler that generates the following access log when using textHandler
// INFO 2019-01-05T18:04:03-08:00 /foo status=200 size=3 method=GET duration=4.225µs refer=
// If you need more customization, you should create a new struct
func NewHttpAccessLogger(logger *dlog.Logger, mux http.Handler) http.Handler {
	return &HttpAccessLogger{
		h:      mux,
		logger: logger,
	}
}

// ServeHTTP implements http.Handler interface
func (ac *HttpAccessLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tracker := httputil.NewTrackedWriter(w)
	start := time.Now()
	// real serving logic
	ac.h.ServeHTTP(&tracker, r)
	// count time and log
	duration := time.Now().Sub(start)
	ac.logger.InfoF(
		r.URL.Path, // TODO: what should the string message be? full url? might check apache and nginx
		dlog.Int("status", tracker.Status()),
		dlog.Int("size", tracker.Size()),
		dlog.Str("method", r.Method),
		dlog.Str("duration", duration.String()), // TODO: use numeric string? what should be the unit?
		dlog.Str("refer", r.Referer()),
	)
}
