// This file manages everything related slog handler that's to be used by bark client.
package client

import (
	"context"
	"github.com/techrail/bark/constants"
	"io"
	"log"
	"log/slog"
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
	log *log.Logger
}

// WithAttrs method is an implementation of slog.Handler interface's method for BarkSlogHandler
// This allows a set of attributes to be added to slog package,
// but right now we're not supporting additional slog attributes.
// This method returns the handler as is for now.
func (handle *BarkSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return handle
}

// WithGroup method is an implementation of slog.Handler interface's method for BarkSlogHandler
// This allows a group to be added to slog package, but right now we're not supporting slog groups.
// This method returns the handler as is for now.
func (handle *BarkSlogHandler) WithGroup(name string) slog.Handler {
	return handle
}

// Enabled method is an implementation of slog.Handler interface's method for BarkSlogHandler
// This method defines which log levels are supported.
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
		log: log.New(out, "", log.Ldate|log.Ltime),
	}
	return handler
}

// SlogHandlerOptions returns slog.HandlerOptions which defines custom log levels.
func SlogHandlerOptions() *slog.HandlerOptions {
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

// newSlogger creates a new logger of type slog.Logger.
func newSlogger(writer io.Writer) *slog.Logger {
	handler := NewBarkSlogHandler(writer)
	return slog.New(handler)
}

// newSlogWithCustomHandler creates a new logger of type slog.Logger with custom slog.Handler object.
func newSlogWithCustomHandler(handler slog.Handler) *slog.Logger {
	return slog.New(handler)
}
