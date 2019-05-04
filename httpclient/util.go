package httpclient

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dyweb/gommon/errors"
)

// DrainResponseBody reads entire response into a byte slice and close the underlying body
// so the connection can be reused.
//
// It replace res.Body with a bytes.Reader so it can be passed to func that reads res.Body
func DrainResponseBody(res *http.Response) ([]byte, error) {
	if res == nil {
		return nil, errors.New("http.Response is nil, can not drain and restore body")
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error read all response body")
	}
	if closeErr := res.Body.Close(); closeErr != nil {
		return b, errors.Wrap(closeErr, "error close drained body")
	}
	// restore body, so it can still be used
	res.Body = ioutil.NopCloser(bytes.NewReader(b))
	return b, nil
}

// DrainAndClose ignore all the errors when drain response body to /dev/null and close it
// It should be used when are using a raw http.Response and don't want to write server lines
// for draining and close the body
//
// You should use FetchToNull most of the time, this only useful when you need some header
// and don't care about the body
func DrainAndClose(res *http.Response) {
	if res == nil {
		return
	}
	errors.Ignore2(io.Copy(ioutil.Discard, res.Body))
	errors.Ignore(res.Body.Close())
}

// JoinPath does not sanitize path like path.Join, which would change https:// to https:/, it only remove duplicated
// slashes to avoid // in url i.e. http://myapi.com/api/v1//comments/1
func JoinPath(base string, sub string) string {
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(sub, "/")
}
