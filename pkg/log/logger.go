package log

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/*
*
初始化zerolog配置
*/
func InitLogs(logOptions *LogConfig) {
	if logOptions == nil {
		return
	}
	var out io.Writer
	var err error
	if logOptions.Path != "" {
		out, err = os.OpenFile(logOptions.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	} else {
		out = os.Stdout
	}

	var level zerolog.Level
	switch logOptions.Level {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	default:
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	switch logOptions.Format {
	case "text":
		out = zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339}
	case "json":
		fallthrough
	default:
	}

	log.Logger = zerolog.New(out).With().Timestamp().Caller().Logger()
}
