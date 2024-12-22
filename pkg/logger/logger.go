package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

var once sync.Once

var log zerolog.Logger

var isDebug bool

func SetDebug(debug bool) {
	isDebug = debug
}

func Get() zerolog.Logger {
	once.Do(func() {

		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		var logLevel int
		var err error

		if isDebug {
			logLevel = int(zerolog.DebugLevel)
		} else {
			logLevel, err = strconv.Atoi(os.Getenv("LOG_LEVEL"))

			if err != nil {
				logLevel = int(zerolog.InfoLevel) // default to INFO
			}
		}

		var output io.Writer

		if os.Getenv("APP_ENV") != "development" {
			fileLogger := &lumberjack.Logger{
				Filename:   "./logs/sudoku.log",
				MaxSize:    5, //
				MaxBackups: 10,
				MaxAge:     14,
				Compress:   true,
			}

			output = zerolog.MultiLevelWriter(fileLogger)
		}

		var gitRevision string

		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, v := range buildInfo.Settings {
				if v.Key == "vcs.revision" {
					gitRevision = v.Value
					break
				}
			}
		}

		log = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Str("git_revision", gitRevision).
			Str("go_version", buildInfo.GoVersion).
			Logger()

	})

	return log
}
