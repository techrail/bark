package ingestion

import (
	"fmt"
	"github.com/techrail/bark/client/channels"
	"github.com/techrail/bark/models"
)

func InsertSingleRequest(logEntry models.BarkLog) {
	if len(channels.ClientChannel) > channels.ClientChannelCapacity-1 {
		fmt.Printf("E# - Channel is full. Cannot push. Log received: | %v\n", logEntry)
		return
	}
	channels.ClientChannel <- logEntry

}
