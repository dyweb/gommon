package benchmarks

import (
	"io/ioutil"
	"testing"

	stdlog "log"

	// zerolog
	"github.com/rs/zerolog"

	// zap
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	// logrus
	"github.com/sirupsen/logrus"

	// apex
	apexlog "github.com/apex/log"
	apexlogconsole "github.com/apex/log/handlers/cli" // TODO: this relies on so many color packages ....
	apexlogjson "github.com/apex/log/handlers/json"

	dlog "github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/log/handlers/json"
)

type ZapDiscard struct {
}

func (z *ZapDiscard) Write(b []byte) (int, error) {
	return ioutil.Discard.Write(b)
}

func (z *ZapDiscard) Sync() error {
	return nil
}

func newZapJsonLogger(lvl zapcore.Level) *zap.Logger {
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeDuration = zapcore.NanosDurationEncoder
	ec.EncodeTime = zapcore.EpochNanosTimeEncoder
	enc := zapcore.NewJSONEncoder(ec)
	return zap.New(zapcore.NewCore(
		enc,
		&ZapDiscard{},
		lvl,
	))
}

func newZapConsoleLogger(lvl zapcore.Level) *zap.Logger {
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeDuration = zapcore.NanosDurationEncoder
	ec.EncodeTime = zapcore.EpochNanosTimeEncoder
	enc := zapcore.NewConsoleEncoder(ec)
	return zap.New(zapcore.NewCore(
		enc,
		&ZapDiscard{},
		lvl,
	))
}

func newZerologJsonLogger() zerolog.Logger {
	// TODO: this may not be the ideal way to init zero logger, see the author's benchmark
	// https://github.com/rs/logbench/blob/master/zerolog_test.go
	return zerolog.New(ioutil.Discard).With().Timestamp().Logger()
}

// https://github.com/rs/zerolog/tree/master#pretty-logging
func newZerologConsoleLogger() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{Out: ioutil.Discard}).With().Timestamp().Logger()
}

func newApexJsonLogger(lvl apexlog.Level) *apexlog.Logger {
	return &apexlog.Logger{
		Handler: apexlogjson.New(ioutil.Discard),
		Level:   lvl,
	}
}

func newApexConsoleLogger(lvl apexlog.Level) *apexlog.Logger {
	return &apexlog.Logger{
		Handler: apexlogconsole.New(ioutil.Discard),
		Level:   lvl,
	}
}

// TODO: use logrus entry might be more reasonable?
func newLogrusJsonLogger(lvl logrus.Level) *logrus.Logger {
	return &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: &logrus.JSONFormatter{},
		Level:     lvl,
	}
}

func newLogrusConsoleLogger(lvl logrus.Level) *logrus.Logger {
	return &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: &logrus.TextFormatter{DisableColors: true},
		Level:     lvl,
	}
}

//func BechmarkDisabledWithoutFieldsJSON(b *testing.B) {
//
//}

// TODO: this also don't call *f method
func BenchmarkWithoutFieldsText(b *testing.B) {
	b.ReportAllocs()
	b.Log("logging a single line text like stdlog without format")
	msg := "TODO: is fixed length msg really a good idea, we should give dynamic length with is more real world"

	b.Run("gommon", func(b *testing.B) {
		logger := dlog.NewTestLogger(dlog.InfoLevel)
		logger.SetHandler(dlog.NewIOHandler(ioutil.Discard))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("std", func(b *testing.B) {
		logger := stdlog.New(ioutil.Discard, "", stdlog.Ldate|stdlog.Ltime)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Print(msg)
			}
		})
	})
	b.Run("zap", func(b *testing.B) {
		logger := newZapConsoleLogger(zap.InfoLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("zap.sugar", func(b *testing.B) {
		logger := newZapConsoleLogger(zap.InfoLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("zerolog", func(b *testing.B) {
		logger := newZerologConsoleLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(msg)
			}
		})
	})
	b.Run("apex", func(b *testing.B) {
		logger := newApexConsoleLogger(apexlog.InfoLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("logrus", func(b *testing.B) {
		logger := newLogrusConsoleLogger(logrus.InfoLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})

}

func BenchmarkWithoutFieldsJSON(b *testing.B) {
	b.ReportAllocs()
	b.Log("logging without fields and without printf, use json output")
	msg := "TODO: is fixed length msg really a good idea, we should give dynamic length with is more real world"
	b.Run("gommon", func(b *testing.B) {
		logger := dlog.NewTestLogger(dlog.InfoLevel)
		logger.SetHandler(json.New(ioutil.Discard))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("zap", func(b *testing.B) {
		logger := newZapJsonLogger(zap.InfoLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("zap.sugar", func(b *testing.B) {
		logger := newZapJsonLogger(zap.InfoLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("zerolog", func(b *testing.B) {
		logger := newZerologJsonLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(msg)
			}
		})
	})
	b.Run("apex", func(b *testing.B) {
		logger := newApexJsonLogger(apexlog.InfoLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("logrus", func(b *testing.B) {
		logger := newLogrusJsonLogger(logrus.InfoLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
}
