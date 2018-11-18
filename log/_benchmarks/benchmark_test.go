package benchmarks

import (
	"io/ioutil"
	"testing"

	"github.com/rs/zerolog"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

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

// TODO: those are all json logger ...
func newGommonLogger(lvl dlog.Level) *dlog.Logger {
	logger := dlog.NewTestLogger(lvl)
	logger.SetHandler(json.New(ioutil.Discard))
	return logger
}

func newZapLogger(lvl zapcore.Level) *zap.Logger {
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

func newZerolog() zerolog.Logger {
	// TODO: this may not be the ideal way to init zero logger, see the author's benchmark
	// https://github.com/rs/logbench/blob/master/zerolog_test.go
	return zerolog.New(ioutil.Discard).With().Timestamp().Logger()
}

//func BechmarkDisabledWithoutFieldsJSON(b *testing.B) {
//
//}

func BenchmarkWithoutFieldsJSON(b *testing.B) {
	b.ReportAllocs()
	b.Logf("logging without fields and without printf, use json output")
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
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.InfoLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.InfoLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("zerorlog", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(msg)
			}
		})
	})
}
