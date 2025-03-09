package client

import (
	"time"
)

type RawLog struct {
	LogTime             time.Time
	LogLevel            string
	ServiceName         string
	ServiceInstanceName string
	Code                string
	Message             string
	MoreData            any
}
