package models

import (
	`encoding/json`
	`fmt`
	"time"
)

// BarkLog is a struct representing a log in Bark
type BarkLog struct {
	Id          int64           `db:"id" json:"id"`
	LogTime     time.Time       `db:"log_time" json:"logTime"`
	LogLevel    int             `db:"log_level" json:"logLevel"`
	ServiceName string          `db:"service_name" json:"serviceName"`
	Code        string          `db:"code" json:"code"`
	Message     string          `db:"msg" json:"msg"`
	MoreData    json.RawMessage `db:"more_data" json:"moreData"`
}

func (b BarkLog) String() string {
	return fmt.Sprintf("Id: %v | LogTime: %v | LogLevel: %v | ServiceName: %v | Code: %v | Message: %v | MoreData: %v \n",
		b.Id, b.LogTime, b.LogLevel, b.ServiceName, b.Code, b.Message, b.MoreData)
}
