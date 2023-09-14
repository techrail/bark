package channels

import (
	`github.com/techrail/bark/models`
)

const LogChannelCapacity = 10000

var LogChannel chan models.BarkLog

func init() {
	LogChannel = make(chan models.BarkLog, LogChannelCapacity)
}
