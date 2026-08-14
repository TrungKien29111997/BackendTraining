package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type LoggerConfig struct {
	Level      string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	IsDev      string
}
type contextKey string

const TraceIDKey contextKey = "trace_id"

func NewLogger(config LoggerConfig) *zerolog.Logger {

	zerolog.TimeFieldFormat = time.RFC3339

	lvl, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	var writer io.Writer

	if config.IsDev == "Dev" {
		writer = PrettyJSONLogger{write: os.Stdout}
	} else {
		writer = &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
	}
	logger := zerolog.New(writer).Level(lvl).With().Timestamp().Logger()
	return &logger
}

type PrettyJSONLogger struct {
	write io.Writer
}

func (w PrettyJSONLogger) Write(p []byte) (n int, err error) {
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, p, "", "  ")
	if err != nil {
		return w.write.Write(p)
	}
	return w.write.Write(prettyJSON.Bytes())
}

func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		return traceID
	}
	return ""
}
