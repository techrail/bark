// Package slogger manages everything related slog handler that's to be used by bark/client.
package slogger

import (
	"log/slog"
	os "os"
)

type Leveler interface {
	Level() slog.Level
}

const (
	LevelNotice = slog.Level(1)
	LevelAlert  = slog.Level(9)
	LevelPanic  = slog.Level(10)
)

var LevelNames = map[slog.Leveler]string{
	LevelAlert:  "ALERT",
	LevelNotice: "NOTICE",
	LevelPanic:  "PANIC",
}

var options = slog.HandlerOptions{
	Level: slog.LevelDebug,
	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.LevelKey {
			level := a.Value.Any().(slog.Level)
			switch level {
			case LevelNotice:
				a.Value = slog.StringValue(LevelNames[LevelNotice])
			case LevelPanic:
				a.Value = slog.StringValue(LevelNames[LevelPanic])
			case LevelAlert:
				a.Value = slog.StringValue(LevelNames[LevelAlert])
			}
		}
		return a
	},
}

func NewHandler() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &options)
	return slog.New(handler)
}
