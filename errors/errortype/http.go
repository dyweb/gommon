package errortype

import "net/http"

// ToHTTPStatus converts error to http code based on their type
// It is similar to docker's https://github.com/moby/moby/blob/master/errdefs/http_helpers.go
// TODO: only a few types are supported and the assertions are actually resource consuming due to using reflection
func ToHTTPStatus(err error) int {
	switch {
	case IsNotFound(err):
		return http.StatusNotFound
	case IsAlreadyExists(err):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
