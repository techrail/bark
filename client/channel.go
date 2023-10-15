package client

import "github.com/techrail/bark/models"

const ChannelCapacity = 100000

var PendingLogsChan chan models.BarkLog

// init allocates fixed memory for pending log channel
func init() {
	PendingLogsChan = make(chan models.BarkLog, ChannelCapacity)
}
