package logger

import (
	"time"
)

// Struct representing a log
type Log struct {
	Id          int64     `db:"id"`
	LogTime     time.Time `db:"log_time"`
	LogLevel    int       `db:"log_level"`
	ServiceName string    `db:"service_name"`
	Code        string    `db:"code"`
	Message     string    `db:"msg"`
	MoreData    string    `db:"more_data"`
}
