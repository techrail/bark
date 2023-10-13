package client

import "github.com/techrail/bark/models"

const ChannelCapacity = 100000

var PendingLogsChan chan models.BarkLog

func init() {
	PendingLogsChan = make(chan models.BarkLog, ChannelCapacity)
}
