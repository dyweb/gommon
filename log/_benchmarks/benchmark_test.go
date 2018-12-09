package benchmarks

import (
	"errors"
	"io/ioutil"
	"testing"

	stdlog "log"

	// zap
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	// zerolog
	"github.com/rs/zerolog"

	// apex
	apexlog "github.com/apex/log"
	apexlogconsole "github.com/apex/log/handlers/cli" // TODO: this relies on so many color packages ....
	apexlogjson "github.com/apex/log/handlers/json"

	// logrus
	"github.com/sirupsen/logrus"

	// klog, fork of glog by k8s
	"k8s.io/klog"

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

// disabled should have not allocation
func BenchmarkDisabledLevelNoFormat(b *testing.B) {
	b.Log("logging at a disabled level")
	msg := "If you support level you should not see me and should not cause allocation, I know I talk too much"
	b.Run("gommon", func(b *testing.B) {
		logger := dlog.NewTestLogger(dlog.ErrorLevel)
		logger.SetHandler(dlog.NewIOHandler(ioutil.Discard))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// TODO: it has 16B allocation due to parameter is interface, size of interface is int64(type), int64(ptr)
				// https://research.swtch.com/interfaces
				logger.Info(msg)
			}
		})
	})
	b.Run("gommon.F", func(b *testing.B) {
		logger := dlog.NewTestLogger(dlog.ErrorLevel)
		logger.SetHandler(dlog.NewIOHandler(ioutil.Discard))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.InfoF(msg)
			}
		})
	})
	b.Run("gommon.check", func(b *testing.B) {
		logger := dlog.NewTestLogger(dlog.ErrorLevel)
		logger.SetHandler(dlog.NewIOHandler(ioutil.Discard))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if logger.IsInfoEnabled() {
					logger.Info(msg)
				}
			}
		})
	})
	b.Run("zap", func(b *testing.B) {
		logger := newZapConsoleLogger(zap.ErrorLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("zap.check", func(b *testing.B) {
		logger := newZapConsoleLogger(zap.ErrorLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if m := logger.Check(zap.InfoLevel, msg); m != nil {
					m.Write()
				}
			}
		})
	})
	b.Run("zap.sugar", func(b *testing.B) {
		logger := newZapConsoleLogger(zap.ErrorLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("zerolog", func(b *testing.B) {
		logger := newZerologConsoleLogger().Level(zerolog.ErrorLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(msg)
			}
		})
	})
	b.Run("apex", func(b *testing.B) {
		logger := newApexConsoleLogger(apexlog.ErrorLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("logrus", func(b *testing.B) {
		logger := newLogrusConsoleLogger(logrus.ErrorLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	//b.Run("klog", func(b *testing.B) {
	//	// TODO: it seems glog can't create individual logger instance?
	//	klog.SetOutput(ioutil.Discard)
	//	//klog.InitFlags()
	//	b.RunParallel(func(pb *testing.PB) {
	//		for pb.Next() {
	//			klog.Info(msg)
	//		}
	//	})
	//})
}

// no fields and don't call *f method for Printf style text formatting
func BenchmarkWithoutFieldsText(b *testing.B) {
	b.ReportAllocs()
	b.Log("logging a single line text like stdlog without format and fields")
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
	b.Run("gommon.F", func(b *testing.B) {
		logger := dlog.NewTestLogger(dlog.InfoLevel)
		logger.SetHandler(dlog.NewIOHandler(ioutil.Discard))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.InfoF(msg)
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
	b.Run("klog", func(b *testing.B) {
		// TODO: it seems glog can't create individual logger instance?
		klog.SetOutput(ioutil.Discard)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				klog.Info(msg) // I think klog's buffer pool is the reason
			}
		})
	})
}

func BenchmarkWithoutFieldsTextFormat(b *testing.B) {
	b.ReportAllocs()
	b.Log("logging a single line text like stdlog with format but without fields")
	format := "TODO: is fixed length msg really a good idea? we should give dynamic length with is more real world %d %s %s"
	i1 := 10086
	s1 := "sub str aaaaa"
	err := errors.New("some error")

	b.Run("gommon", func(b *testing.B) {
		logger := dlog.NewTestLogger(dlog.InfoLevel)
		logger.SetHandler(dlog.NewIOHandler(ioutil.Discard))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infof(format, i1, s1, err)
			}
		})
	})
	b.Run("zap.sugar", func(b *testing.B) {
		logger := newZapConsoleLogger(zap.InfoLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infof(format, i1, s1, err)
			}
		})
	})
	// TODO: seems zerolog console logger also don't have *f variant
	b.Run("apex", func(b *testing.B) {
		logger := newApexConsoleLogger(apexlog.InfoLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infof(format, i1, s1, err)
			}
		})
	})
	b.Run("logrus", func(b *testing.B) {
		logger := newLogrusConsoleLogger(logrus.InfoLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infof(format, i1, s1, err)
			}
		})
	})
	b.Run("klog", func(b *testing.B) {
		// TODO: it seems glog can't create individual logger instance?
		klog.SetOutput(ioutil.Discard)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				klog.Infof(format, i1, s1, err)
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
	b.Run("gommon.F", func(b *testing.B) {
		logger := dlog.NewTestLogger(dlog.InfoLevel)
		logger.SetHandler(json.New(ioutil.Discard))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.InfoF(msg)
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

func BenchmarkCallerJSON(b *testing.B) {
	b.ReportAllocs()
	b.Log("logging without fields and without printf, use json output and enable log file line")
	msg := "TODO: is fixed length msg really a good idea, we should give dynamic length with is more real world"
	b.Run("gommon", func(b *testing.B) {
		logger := dlog.NewTestLogger(dlog.InfoLevel)
		logger.SetHandler(json.New(ioutil.Discard))
		logger.EnableSource()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("gommon.F", func(b *testing.B) {
		logger := dlog.NewTestLogger(dlog.InfoLevel)
		logger.SetHandler(json.New(ioutil.Discard))
		logger.EnableSource()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.InfoF(msg)
			}
		})
	})
}

func BenchmarkWithContextNoFieldsJSON(b *testing.B) {
	b.ReportAllocs()
	b.Log("logging with context attached to logger (entry/event) no text format, no fields, use json output")
	msg := "TODO: is fixed length msg really a good idea, we should give dynamic length with is more real world"
	b.Run("gommon", func(b *testing.B) {
		logger := dlog.NewTestLogger(dlog.InfoLevel)
		logger.SetHandler(json.New(ioutil.Discard))
		// TODO: generate unified fields for all logging libraries
		logger.AddFields(dlog.Int("i1", 1), dlog.Str("s1", "v1"))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("gommon.F", func(b *testing.B) {
		logger := dlog.NewTestLogger(dlog.InfoLevel)
		logger.SetHandler(json.New(ioutil.Discard))
		// TODO: generate unified fields for all logging libraries
		logger.AddFields(dlog.Int("i1", 1), dlog.Str("s1", "v1"))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.InfoF(msg)
			}
		})
	})
	b.Run("zap", func(b *testing.B) {
		logger := newZapJsonLogger(zap.InfoLevel).
			With(zap.Int("i1", 1), zap.String("s1", "v1"))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("zap.sugar", func(b *testing.B) {
		logger := newZapJsonLogger(zap.InfoLevel).
			With(zap.Int("i1", 1), zap.String("s1", "v1")).
			Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("zerolog", func(b *testing.B) {
		logger := newZerologJsonLogger().
			With().Int("i1", 1).Str("s1", "v1").Logger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(msg)
			}
		})
	})
	b.Run("apex", func(b *testing.B) {
		logger := newApexJsonLogger(apexlog.InfoLevel).
			WithFields(apexlog.Fields{
				"i1": 1,
				"s1": "v1",
			})
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
	b.Run("logrus", func(b *testing.B) {
		logger := newLogrusJsonLogger(logrus.InfoLevel).
			WithFields(logrus.Fields{
				"i1": 1,
				"s1": "v1",
			})
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg)
			}
		})
	})
}

func BenchmarkNoContextWithFieldsJSON(b *testing.B) {
	b.ReportAllocs()
	b.Log("logging with fields at log site, no context attached to logger (entry/event) no text format, use json output")
	msg := "TODO: is fixed length msg really a good idea, we should give dynamic length with is more real world"
	b.Run("gommon.F", func(b *testing.B) {
		logger := dlog.NewTestLogger(dlog.InfoLevel)
		logger.SetHandler(json.New(ioutil.Discard))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// TODO: generate unified fields for all logging libraries
				logger.InfoF(msg,
					dlog.Int("i1", 1),
					dlog.Str("s1", "v1"),
				)
			}
		})
	})
	b.Run("zap", func(b *testing.B) {
		logger := newZapJsonLogger(zap.InfoLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(msg, zap.Int("i1", 1), zap.String("s1", "v1"))
			}
		})
	})
	b.Run("zap.sugar", func(b *testing.B) {
		logger := newZapJsonLogger(zap.InfoLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infow(msg, "i1", 1, "s1", "v1")
			}
		})
	})
	b.Run("zerolog", func(b *testing.B) {
		logger := newZerologJsonLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Int("i1", 1).Str("s1", "v1").Msg(msg)
			}
		})
	})
	b.Run("apex", func(b *testing.B) {
		logger := newApexJsonLogger(apexlog.InfoLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(apexlog.Fields{
					"i1": 1,
					"s1": "v1",
				}).Info(msg)
			}
		})
	})
	b.Run("logrus", func(b *testing.B) {
		logger := newLogrusJsonLogger(logrus.InfoLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(logrus.Fields{
					"i1": 1,
					"s1": "v1",
				}).Info(msg)
			}
		})
	})
}
