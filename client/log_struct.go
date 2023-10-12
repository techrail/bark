package client

import (
	"time"
)

type RawLog struct {
	Id          int64
	LogTime     time.Time
	LogLevel    string
	ServiceName string
	SessionName string
	Code        string
	Message     string
	MoreData    any
}
