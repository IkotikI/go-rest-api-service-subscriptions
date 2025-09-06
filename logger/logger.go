package logger

import (
	"flag"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger(level string, consoleOut bool) *zerolog.Logger {
	lvl, _ := filterStringLogLevel(level)

	return initLogger(lvl, consoleOut)
}

var logLevel = flag.String("log-level", "info", "log level: fatal, panic, error, warn, info, debug, trace")

func InitLoggerByFlag(defLogLevel string, consoleOut bool) *zerolog.Logger {
	flag.Parse()
	var logger *zerolog.Logger

	if *logLevel == "" {
		logger = InitLogger(defLogLevel, consoleOut)
		return logger
	}

	lvl, ok := filterStringLogLevel(*logLevel)
	// fmt.Printf("flag: %s, ok: %v\n", *logLevel, ok)
	if ok {
		logger = initLogger(lvl, consoleOut)
		log.Debug().Msgf("logger initialized with level \"%s\"", *logLevel)
	} else {
		logger = InitLogger(defLogLevel, consoleOut)
		logger.Warn().Msgf("wrong log level \"%s\" is provided, fall back to default \"%s\"", *logLevel, defLogLevel)
	}

	return logger
}

func initLogger(lvl zerolog.Level, consoleOut bool) *zerolog.Logger {
	zerolog.SetGlobalLevel(lvl)

	if consoleOut {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).With().Timestamp().Logger()
	}

	return &log.Logger
}

func filterStringLogLevel(level string) (zerolog.Level, bool) {
	var zlvl zerolog.Level
	switch {
	case level == "fatal":
		zlvl = zerolog.FatalLevel
	case level == "panic":
		zlvl = zerolog.PanicLevel
	case level == "error":
		zlvl = zerolog.ErrorLevel
	case level == "warn":
		zlvl = zerolog.WarnLevel
	case level == "info":
		zlvl = zerolog.InfoLevel
	case level == "debug":
		zlvl = zerolog.DebugLevel
	case level == "trace":
		zlvl = zerolog.TraceLevel
	default:
		return zerolog.WarnLevel, false
	}
	return zlvl, true
}
