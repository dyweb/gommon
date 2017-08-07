package requests

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

const AuthorizationHeader = "Authorization"

func GenerateBasicAuth(username string, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}

func ExtractBasicAuth(val string) (username string, password string, err error) {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Authorization
	s := strings.Split(val, " ")
	if len(s) != 2 {
		err = errors.New("invalid authorization header, must have type and base64 encoded value separated by space")
		return
	}
	tpe := strings.ToLower(s[0])
	if tpe != "basic" {
		err = errors.Errorf("got %s instead basic auth", tpe)
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		err = errors.New("invalid base64 encoding")
		return
	}
	ss := strings.Split(string(decoded), ":")
	if len(ss) != 2 {
		errors.Errorf("invalid username:password, got %s segments after split by ':'", len(ss))
	}
	username = ss[0]
	password = ss[1]
	err = nil
	return
}

func ExtractBasicAuthFromRequest(r *http.Request) (string, string, error) {
	val := r.Header.Get(AuthorizationHeader)
	if val == "" {
		return "", "", errors.New("Authorization header is not set")
	}
	return ExtractBasicAuth(val)
}
