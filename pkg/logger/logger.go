package logger

import (
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Log is global logger
	Log *zap.Logger

	// timeFormat is custom Time format
	customTimeFormat string

	// onceInit guarantee initialize logger only once
	onceInit sync.Once
)

// customTimeEncoder encode Time to our custom format
// This example how we can customize zap default functionality
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(customTimeFormat))
}

// Init initializes log by input parameters
// lvl - global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
// timeFormat - custom time format for logger of empty string to use default
func Init(lvl int, timeFormat string, mode string) error {
	var err error

	onceInit.Do(func() {
		// First, define our level-handling logic.
		globalLevel := zapcore.Level(lvl)

		// High-priority output: standard error
		// low-priority output: standard out.
		// It is useful for Kubernetes deployment.
		// Kubernetes interprets os.Stdout log items as INFO and os.Stderr log items as ERROR by default.
		highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= globalLevel && lvl < zapcore.ErrorLevel
		})
		consoleInfos := zapcore.Lock(os.Stdout)
		consoleErrors := zapcore.Lock(os.Stderr)

		// Configure console output.
		var useCustomTimeFormat bool
		var ecfg zapcore.EncoderConfig

		if mode == "prod" {
			ecfg = zap.NewProductionEncoderConfig()

			if len(timeFormat) > 0 {
				customTimeFormat = timeFormat
				ecfg.EncodeTime = customTimeEncoder
				useCustomTimeFormat = true
			}
		} else {
			ecfg = zap.NewDevelopmentEncoderConfig()
			ecfg.EncodeCaller = zapcore.FullCallerEncoder
			ecfg.CallerKey = "caller"
			ecfg.MessageKey = "message"
			ecfg.LevelKey = "level"
			ecfg.TimeKey = "time"
		}

		consoleEncoder := zapcore.NewJSONEncoder(ecfg)

		// Join the outputs, encoders, and level-handling functions into zapcore.
		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
			zapcore.NewCore(consoleEncoder, consoleInfos, lowPriority),
		)

		// From a zapcore.Core, it's easy to construct a Logger.
		if mode == "prod" {
			Log = zap.New(core)
			Log.Info("Mode is prod")
		} else {
			Log = zap.New(core, zap.AddCaller())
			Log.Info("Mode is dev")
		}

		zap.RedirectStdLog(Log)

		if !useCustomTimeFormat {
			Log.Warn("time format for logger is not provided - use zap default")
		}
	})

	return err
}
