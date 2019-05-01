package httpclient_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

// writeJson is used for testing
func writeJSON(t *testing.T, res http.ResponseWriter, val interface{}) {
	b, err := json.Marshal(val)
	if err != nil {
		t.Fatalf("error encode json: %s", err)
		return
	}
	if _, err = res.Write(b); err != nil {
		t.Fatalf("error write to http response: %s", err)
	}
}
