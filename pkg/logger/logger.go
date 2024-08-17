package logger

import (
	"fmt"
	"os"
	"strings"

	"demo-temporal-project/configs"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/rs/zerolog"
)

var _ log.Logger = (*ZerologLogger)(nil)

type ZerologLogger struct {
	Logger *zerolog.Logger
}

func NewLogger(conf *configs.Log) log.Logger {

	var l zerolog.Level

	switch strings.ToLower(conf.Level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	skipFrameCount := 2
	loggerOutput := os.Stdout
	logger := zerolog.New(loggerOutput).With().Timestamp().CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).Logger()
	return &ZerologLogger{Logger: &logger}
}

func (l *ZerologLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 {
		l.Logger.Warn().Msg("No log message found.")
		return nil
	}
	var err error
	switch level {
	case log.LevelDebug:
		err = logWithArgs(l.Logger.Debug(), keyvals)
	case log.LevelInfo:
		err = logWithArgs(l.Logger.Info(), keyvals)
	case log.LevelWarn:
		err = logWithArgs(l.Logger.Warn(), keyvals)
	case log.LevelError:
		err = logWithArgs(l.Logger.Error(), keyvals)
	}

	if err != nil {
		return err
	}
	return nil
}

func logWithArgs(event *zerolog.Event, keyvals ...interface{}) error {
	if len(keyvals)%2 == 0 {
		for i := 0; i < len(keyvals); i += 2 {
			event = event.Str(fmt.Sprint(keyvals[i]), fmt.Sprint(keyvals[i+1]))
		}
		event.Msg("")
	} else {
		var i int
		for i = 0; i < len(keyvals)-1; i += 2 {
			event = event.Str(fmt.Sprint(keyvals[i]), fmt.Sprint(keyvals[i+1]))
		}
		event.Msg(fmt.Sprint(keyvals[i]))
	}
	return nil
}

var ProviderSet = wire.NewSet(NewLogger)
