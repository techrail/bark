package client

import (
	"fmt"
	"github.com/techrail/bark/models"
)

func InsertSingleRequest(logEntry models.BarkLog) {
	if len(PendingLogsChan) > ChannelCapacity-1 {
		fmt.Printf("E#1LB9MN - Channel is full. Cannot push. Log received: | %v\n", logEntry)
		return
	}
	Wg.Add(1)
	PendingLogsChan <- logEntry
}
