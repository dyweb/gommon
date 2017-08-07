package stdlib

import (
	"testing"
	"runtime"
	"go/build"

	asst "github.com/stretchr/testify/assert"
	"strings"
)

func TestRuntime_GOROOT(t *testing.T) {
	t.Log(runtime.GOROOT())
	t.Log(build.Default.GOPATH)
}

func TestRuntime_GetPackage(t *testing.T) {
	assert := asst.New(t)
	pc, file, line, ok := runtime.Caller(0)
	assert.True(ok)
	t.Log(file)
	t.Log(line)
	fn := runtime.FuncForPC(pc)
	fnName := fn.Name()
	lastDot := strings.LastIndex(fnName, ".")
	pkgName := fnName[:lastDot]
	t.Log(pkgName)
}