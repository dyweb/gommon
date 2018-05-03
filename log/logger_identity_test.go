package log

import (
	"fmt"
	"testing"

	"github.com/dyweb/gommon/util/testutil"
	asst "github.com/stretchr/testify/assert"
)

var lg = NewPackageLogger()

func foo() *Logger {
	funcLog := NewFunctionLogger(lg)
	return funcLog
}

type Foo struct {
	log *Logger
}

func (f *Foo) GetLogger() *Logger {
	return f.log
}

func (f *Foo) SetLogger(logger *Logger) {
	f.log = logger
}

func (f *Foo) LoggerIdentity(justCallMe func() *Identity) *Identity {
	return justCallMe()
}

func (f *Foo) method() *Logger {
	mlog := NewMethodLogger(f.log)
	return mlog
}

var dummyFoo = &Foo{} // used for get struct logger identity

func TestNewPackageLogger(t *testing.T) {
	assert := asst.New(t)
	id := lg.id
	assert.Equal(PackageLogger, id.Type)
	assert.Equal("pkg", id.Type.String())
	assert.Equal("init", id.Function)
	assert.Equal(testutil.GOPATH()+"/src/github.com/dyweb/gommon/log/logger_identity_test.go:11",
		fmt.Sprintf("%s:%d", id.File, id.Line))
}

func TestNewFunctionLogger(t *testing.T) {
	assert := asst.New(t)
	flog := foo()
	id := flog.id
	assert.Equal(FunctionLogger, id.Type)
}

func TestNewStructLogger(t *testing.T) {
	assert := asst.New(t)
	slog := NewStructLogger(lg, dummyFoo)
	id := slog.id
	assert.Equal(StructLogger, id.Type)
	assert.Equal("struct", id.Type.String())
	assert.Equal("Foo", id.Struct)
	assert.Equal(MagicStructLoggerFunctionName, id.Function)
	assert.Equal(testutil.GOPATH()+"/src/github.com/dyweb/gommon/log/logger_identity_test.go:31",
		fmt.Sprintf("%s:%d", id.File, id.Line))
}

func TestNewMethodLogger(t *testing.T) {
	assert := asst.New(t)
	slog := NewStructLogger(lg, dummyFoo)
	dummyFoo.log = slog
	mlog := dummyFoo.method()
	id := mlog.id
	assert.Equal(MethodLogger, id.Type)
	assert.Equal("method", id.Type.String())
	assert.Equal("Foo", id.Struct)
	assert.Equal("method", id.Function)
	assert.Equal(testutil.GOPATH()+"/src/github.com/dyweb/gommon/log/logger_identity_test.go:35",
		fmt.Sprintf("%s:%d", id.File, id.Line))
}
