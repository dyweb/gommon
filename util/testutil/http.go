package testutil

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func GetBody(t *testing.T, c *http.Client, url string) []byte {
	res, err := c.Get(url)
	if err != nil {
		t.Fatalf("error GET %s: %s", url, err)
		return nil
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("error read body of %s: %s", url, err)
	}
	res.Body.Close()
	return b
}
