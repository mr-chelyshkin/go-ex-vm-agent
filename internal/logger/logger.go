package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*zerolog.Logger
}

func New(cfg Config) (*Logger, error) {
	if err := cfg.Validate(); err != nil {
		return nil, initError(err.Error())
	}

	level, err := parseLogLevel(cfg.Level)
	if err != nil {
		return nil, initError(err.Error())
	}
	writer, err := setupWriter(cfg)
	if err != nil {
		return nil, initError(err.Error())
	}

	var logger zerolog.Logger
	logger = logger.Level(level)
	switch cfg.Format {
	case FormatConsole:
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        writer,
			TimeFormat: time.RFC3339,
		})
	case FormatJSON:
		logger = zerolog.New(writer)
	default:
		return nil, initError("unsupported log format: %s", cfg.Format)
	}
	if cfg.Output != OutputJournal {
		logger = logger.With().Timestamp().Logger()
	}
	return &Logger{Logger: &logger}, nil
}

func parseLogLevel(level LogLevel) (zerolog.Level, error) {
	switch level {
	case LevelDebug:
		return zerolog.DebugLevel, nil
	case LevelInfo:
		return zerolog.InfoLevel, nil
	case LevelWarn:
		return zerolog.WarnLevel, nil
	case LevelError:
		return zerolog.ErrorLevel, nil
	case LevelFatal:
		return zerolog.FatalLevel, nil
	case LevelPanic:
		return zerolog.PanicLevel, nil
	case LevelDisabled:
		return zerolog.Disabled, nil
	default:
		return zerolog.InfoLevel, fmt.Errorf("unknown level: %s", level)
	}
}

func setupWriter(cfg Config) (io.Writer, error) {
	switch cfg.Output {
	case OutputStdout:
		return os.Stdout, nil
	case OutputStderr:
		return os.Stderr, nil
	case OutputJournal:
		return os.Stdout, nil
	case OutputFile:
		if err := os.MkdirAll(filepath.Dir(cfg.Path), 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %v", err)
		}

		return &lumberjack.Logger{
			Filename:   cfg.Path,
			MaxAge:     cfg.MaxAge,
			MaxSize:    cfg.MaxSize,
			Compress:   cfg.Compress,
			MaxBackups: cfg.MaxBackups,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported output type: %s", cfg.Output)
	}
}
