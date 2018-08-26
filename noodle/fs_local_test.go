package noodle

import (
	"net/http"
	"net/http/httptest"
	"testing"

	asst "github.com/stretchr/testify/assert"

	"github.com/dyweb/gommon/util/testutil"
)

func TestNewLocal(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(NewLocal(".")))
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := srv.Client()

	t.Run("can access file", func(t *testing.T) {
		assert := asst.New(t)
		assert.Equal(testutil.ReadFixture(t, "README.md"),
			testutil.GetBody(t, c, srv.URL+"/README.md"))
	})

	t.Run("can NOT read dir", func(t *testing.T) {
		assert := asst.New(t)
		assert.Equal("404 page not found\n", string(testutil.GetBody(t, c, srv.URL+"/doc")))
	})
}

func TestNewLocalUnsafe(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(NewLocalUnsafe(".")))
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := srv.Client()

	t.Run("can access file", func(t *testing.T) {
		assert := asst.New(t)
		assert.Equal(testutil.ReadFixture(t, "README.md"),
			testutil.GetBody(t, c, srv.URL+"/README.md"))
	})

	t.Run("can read dir", func(t *testing.T) {
		assert := asst.New(t)
		b := testutil.GetBody(t, c, srv.URL+"/doc")
		assert.NotEqual("404 page not found\n", string(b))
		assert.Contains(string(b), `<a href="README.md">README.md</a>`)
	})
}
