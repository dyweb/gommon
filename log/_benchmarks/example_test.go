package benchmarks

import (
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// this file tests what it looks like when logging using different library and its different handlers

func TestZapJson(t *testing.T) {
	// https://github.com/sandipb/zap-examples

	t.Run("production logger", func(t *testing.T) {
		logger, err := zap.NewProduction()
		if err != nil {
			t.Fatal(err)
			return
		}
		logger.Info("hi")
		// {"level":"info","ts":1542779858.448864,"caller":"_benchmarks/example_test.go:19","msg":"hi"}
	})

	t.Run("production config", func(t *testing.T) {
		// TODO: this config is from zap's own benchmark, it does not have caller enabled, strange, seems need key in config
		// https://github.com/sandipb/zap-examples/tree/master/src/customlogger#customizing-the-encoder
		// it seems you need to add they keys in config ...
		ec := zap.NewProductionEncoderConfig()
		ec.EncodeDuration = zapcore.NanosDurationEncoder
		ec.EncodeTime = zapcore.EpochNanosTimeEncoder
		enc := zapcore.NewJSONEncoder(ec)
		logger := zap.New(zapcore.NewCore(
			enc,
			os.Stderr,
			zapcore.InfoLevel,
		))
		logger.Info("this is a message")
		// {"level":"info","ts":1542778510834696318,"msg":"this is a message"}
		logger.Named("jack").Info("this is a named message")
		// {"level":"info","ts":1542778510834708500,"logger":"jack","msg":"this is a named message"}
		logger.Named("jack").Named("marry").Info("what's my name")
		// {"level":"info","ts":1542781224287444070,"logger":"jack.marry","msg":"what's my name"}

		t.Run("context", func(t *testing.T) {
			logger.With(zap.Int("count", 1), zap.String("str", `"need escape"`)).Info("ha")
			// {"level":"info","ts":1542779234131875417,"msg":"ha","count":1,"str":"\"need escape\""}
			ctxLogger := logger.With(zap.Int("count", 1))
			ctxLogger.Info("this is the log")
			// {"level":"info","ts":1542779282702624079,"msg":"this is the log","count":1}
			ctxLogger.With(zap.Bool("b", true)).Info("yep")
			// {"level":"info","ts":1542779322646305013,"msg":"yep","count":1,"b":true}

			// NOTE: it will have duplicated key
			ctxLogger.With(zap.Bool("a", true), zap.Bool("a", false)).Info("dup?")
			// {"level":"info","ts":1542783708445651389,"msg":"dup?","count":1,"a":true,"a":false}

			// TODO: https://github.com/sandipb/zap-examples/tree/master/src/customlogger#changing-logger-behavior-on-the-fly
			// things like AddCaller, AddStacktrace seems no longer exists
		})
	})
}

func TestZapConsole(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
		return
	}
	logger.Info("text based")
	// 2018-11-20T22:56:08.273-0800	INFO	_benchmarks/example_test.go:64	text based
	// NOTE: fields are still encoded as json ...
	logger.With(zap.Int("a", 1), zap.String("b", "aaa")).Info("ho ho ho")
	// 2018-11-20T22:56:08.273-0800	INFO	_benchmarks/example_test.go:65	ho ho ho	{"a": 1, "b": "aaa"}
}
