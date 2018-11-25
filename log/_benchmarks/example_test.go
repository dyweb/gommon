package benchmarks

import (
	stdlog "log"
	"os"
	"testing"

	"k8s.io/klog"

	// zap
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/rs/zerolog"

	// apex
	apexlog "github.com/apex/log"
	apexlogconsole "github.com/apex/log/handlers/cli" // TODO: this relies on so many color packages ....
	apexlogjson "github.com/apex/log/handlers/json"

	"github.com/sirupsen/logrus"
)

// this file shows how to use different logging library and
// what it looks like when logging using different formats

// the order is
// - gommon
// - std
// - zap
// - zerolog
// - apex
// - logrus

// log in standard library, no level and field support
func TestStd(t *testing.T) {
	logger := stdlog.New(os.Stdout, "", stdlog.LstdFlags)
	logger.Print("a", 1)
	logger.Printf("%s %d", "a", 1)
	logger.Println("a", 1)
	//2018/11/23 21:13:59 a1
	//2018/11/23 21:13:59 a 1
	//2018/11/23 21:13:59 a 1
}

// zap https://github.com/uber-go/zap
func TestZap(t *testing.T) {
	// https://github.com/sandipb/zap-examples

	t.Run("json", func(t *testing.T) {
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

			// TODO: might move context to top level since it's an important feature
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
	})

	t.Run("console", func(t *testing.T) {
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
	})
}

// zerolog https://github.com/rs/zerolog
func TestZerolog(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Info().Msg("show some info")
		logger.Info().Str("f1", "v1").Msg("have key-value")
		//{"level":"info","time":"2018-11-23T21:25:39-08:00","message":"show some info"}
		//{"level":"info","f1":"v1","time":"2018-11-23T21:25:39-08:00","message":"have key-value"}

		t.Run("context", func(t *testing.T) {
			ctxLogger := logger.With().Str("base", "value").Logger()
			ctxLogger.Info().Msg("inherit context")
			ctxLogger.Info().Str("extra", "value").Msg("extra field")
			//{"level":"info","base":"value","time":"2018-11-23T21:25:39-08:00","message":"inherit context"}
			//{"level":"info","base":"value","extra":"value","time":"2018-11-23T21:25:39-08:00","message":"extra field"}
		})
	})
	t.Run("console", func(t *testing.T) {
		logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
		logger.Info().Msg("level have color? yes") // info is yellow
		logger.Warn().Msg("warn is yellow?")       // warn is red
		logger.Info().Str("f1", "v1").Msg("field has color? no")
		//9:29PM INF level have color? yes
		//9:29PM WRN warn is yellow?
		//9:29PM INF field has color? no f1=v1
	})
}

// apex/log https://github.com/apex/log
func TestApex(t *testing.T) {
	// TODO: its fields are not flat w/ level and message?
	t.Run("json", func(t *testing.T) {
		logger := apexlog.Logger{
			Handler: apexlogjson.New(os.Stdout),
			Level:   apexlog.InfoLevel,
		}
		logger.Info("hi")
		logger.WithField("f1", "v1").Info("have field")
		//{"fields":{},"level":"info","timestamp":"2018-11-23T21:43:00.637421923-08:00","message":"hi"}
		//{"fields":{"f1":"v1"},"level":"info","timestamp":"2018-11-23T21:43:00.637500665-08:00","message":"have filed"}

		t.Run("context", func(t *testing.T) {
			ctxEntry := logger.WithField("base", "value")
			ctxEntry.Info("inherit context")
			ctxEntry.WithField("extra", "value").Info("extra field")
			//{"fields":{"base":"value"},"level":"info","timestamp":"2018-11-23T21:43:00.637524216-08:00","message":"inherit context"}
			//{"fields":{"base":"value","extra":"value"},"level":"info","timestamp":"2018-11-23T21:43:00.637534403-08:00","message":"extra field"}
		})

	})
	t.Run("console", func(t *testing.T) {
		logger := apexlog.Logger{
			Handler: apexlogconsole.New(os.Stdout),
			Level:   apexlog.InfoLevel,
		}
		logger.Info("level have color?") // NOTE: it seems color is disabled when not tty and become a dot ....
		logger.WithField("f1", "v1").Info("field have color?")
		//   • level have color?
		//   • field have color?         f1=v1
	})
}

func TestLogrus(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		logger := logrus.New()
		logger.SetOutput(os.Stdout)
		logger.SetLevel(logrus.InfoLevel)
		logger.SetFormatter(&logrus.JSONFormatter{})
		//logger.SetReportCaller()

		logger.Info("hi")
		logger.WithField("f1", "v1").Info("have field")
		//{"level":"info","msg":"hi","time":"2018-11-23T22:51:07-08:00"}
		//{"f1":"v1","level":"info","msg":"have field","time":"2018-11-23T22:51:07-08:00"}

		t.Run("context", func(t *testing.T) {
			ctxEntry := logger.WithField("base", "value")
			ctxEntry.Info("inherit context")
			ctxEntry.WithField("extra", "value").Info("extra field")
			//{"base":"value","level":"info","msg":"inherit context","time":"2018-11-23T23:19:56-08:00"}
			//{"base":"value","extra":"value","level":"info","msg":"extra field","time":"2018-11-23T23:19:56-08:00"}
		})
	})

	t.Run("console", func(t *testing.T) {
		logger := logrus.New()
		logger.SetOutput(os.Stdout)
		logger.SetLevel(logrus.InfoLevel)
		logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})

		logger.Info("level have color? yes, when tty or forced")
		logger.WithField("f1", "v1").Info("field has color? yes")
		//INFO[0000] level have color? yes, when tty or forced
		//INFO[0000] field has color?                              f1=v1
	})
}

func TestKlog(t *testing.T) {
	t.Run("console", func(t *testing.T) {
		klog.SetOutput(os.Stderr)
		klog.Info("just log something")
	})
}
