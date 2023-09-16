package ingestion

import (
	"fmt"

	"github.com/techrail/bark/channels"
	"github.com/techrail/bark/models"
)

func InsertSingle(logEntry models.BarkLog) {
	if len(channels.LogChannel) > channels.LogChannelCapacity-1 {
		fmt.Printf("E#1KDY0O - Channel is full. Cannot push. Log received: | %v\n", logEntry)
		return
	}
	_, err := logEntry.ValidateForInsert()

	if err == nil {
		channels.LogChannel <- logEntry
	}
}

func InsertMultiple(logEntries []models.BarkLog) {
	for _, logEntry := range logEntries {
		if len(channels.LogChannel) > channels.LogChannelCapacity-1 {
			fmt.Printf("E#1KDZD2 - Channel is full. Cannot push. Log received: | %v\n", logEntry)
			return
		}

		_, err := logEntry.ValidateForInsert()

		if err == nil {
			channels.LogChannel <- logEntry
		}
	}
}
