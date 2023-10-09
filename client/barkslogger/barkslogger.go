// Package barkslogger manages everything related slog handler that's to be used by bark client.
package barkslogger

import (
	"context"
	"github.com/techrail/bark/constants"
	"io"
	"log"
	"log/slog"
	"os"
)

// Constants for custom log levels in bark.
const (
	LvlPanic   = slog.Level(10)
	LvlAlert   = slog.Level(9)
	LvlError   = slog.Level(8)
	LvlWarning = slog.Level(4)
	LvlNotice  = slog.Level(3)
	LvlInfo    = slog.Level(0)
	LvlDebug   = slog.Level(-4)
)

// BarkSlogHandler implements interface slog.Handler.
type BarkSlogHandler struct {
	slog.Handler
	log *log.Logger
}

func (handle *BarkSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return slog.NewJSONHandler(os.Stdout, nil)
}

func (handle *BarkSlogHandler) WithGroup(name string) slog.Handler {
	return handle
}

func (handle *BarkSlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	if level == LvlPanic || level == LvlAlert || level == LvlError ||
		level == LvlWarning || level == LvlNotice || level == LvlInfo ||
		level == LvlDebug {
		return true
	}
	return false
}

// Handle is an implementation of slog.Handler interface's methods for BarkSlogHandler.
func (handle *BarkSlogHandler) Handle(ctx context.Context, record slog.Record) error {
	level := record.Level.String()
	switch record.Level {
	case LvlPanic:
		level = constants.Panic
	case LvlAlert:
		level = constants.Alert
	case LvlError:
		level = constants.Error
	case LvlWarning:
		level = constants.Warning
	case LvlNotice:
		level = constants.Notice
	case LvlInfo:
		level = constants.Info
	case LvlDebug:
		level = constants.Debug
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
				case LvlPanic:
					attr.Value = slog.StringValue(constants.Panic)
				case LvlAlert:
					attr.Value = slog.StringValue(constants.Alert)
				case LvlError:
					attr.Value = slog.StringValue(constants.Error)
				case LvlWarning:
					attr.Value = slog.StringValue(constants.Warning)
				case LvlNotice:
					attr.Value = slog.StringValue(constants.Notice)
				case LvlInfo:
					attr.Value = slog.StringValue(constants.Info)
				case LvlDebug:
					attr.Value = slog.StringValue(constants.Debug)
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
