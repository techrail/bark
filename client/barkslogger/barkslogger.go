// Package barkslogger manages everything related slog handler that's to be used by bark client.
package barkslogger

import (
	"context"
	"github.com/techrail/bark/constants"
	"io"
	"log"
	"log/slog"
)

// Constants for custom log levels in bark.
const (
	LvlNotice = slog.Level(1)
	LvlAlert  = slog.Level(9)
	LvlPanic  = slog.Level(10)
)

// BarkSlogHandler implements interface slog.Handler.
type BarkSlogHandler struct {
	slog.Handler
	log *log.Logger
}

// Handle is an implementation of slog.Handler interface's methods for BarkSlogHandler.
func (handle *BarkSlogHandler) Handle(ctx context.Context, record slog.Record) error {
	level := record.Level.String()
	switch record.Level {
	case LvlNotice:
		level = constants.Notice
	case LvlAlert:
		level = constants.Alert
	case LvlPanic:
		level = constants.Panic
	}
	message := record.Message
	handle.log.Println(level, message)
	return nil
}

// NewBarkSlogHandler returns an object of BarkSlogHandler
func NewBarkSlogHandler(out io.Writer) *BarkSlogHandler {
	handler := &BarkSlogHandler{
		Handler: slog.NewJSONHandler(out, nil),
		log:     log.New(out, "", log.Ldate|log.Ltime),
	}
	return handler
}

// Options returns slog.HandlerOptions which defines custom log levels.
func Options() *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.LevelKey {
				level := attr.Value.Any().(slog.Level)
				switch level {
				case LvlNotice:
					attr.Value = slog.StringValue(constants.Notice)
				case LvlPanic:
					attr.Value = slog.StringValue(constants.Panic)
				case LvlAlert:
					attr.Value = slog.StringValue(constants.Alert)
				}
			}
			return attr
		},
	}
}

// New creates a new logger of type slog.Logger.
func New(writer io.Writer) *slog.Logger {
	handler := NewBarkSlogHandler(writer)
	return slog.New(handler)
}

// NewWithCustomHandler creates a new logger of type slog.Logger with custom slog.Handler object.
func NewWithCustomHandler(handler slog.Handler) *slog.Logger {
	return slog.New(handler)
}
