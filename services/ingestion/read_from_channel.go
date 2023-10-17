package ingestion

import (
	"fmt"
	"github.com/techrail/bark/constants"
	"github.com/techrail/bark/resources"

	"github.com/techrail/bark/channels"
	"github.com/techrail/bark/models"
)

func InsertSingle(logEntry models.BarkLog) {
	if len(channels.LogChannel) > constants.ServerLogInsertionChannelCapacity-1 {
		fmt.Printf("E#1KDY0O - Channel is full. Cannot push. Log received: | %v\n", logEntry)
		return
	}
	logEntry, err := logEntry.ValidateForInsert()

	if err == nil {
		resources.ServerDbSaverWg.Add(1)
		channels.LogChannel <- logEntry
	}
}

func InsertMultiple(logEntries []models.BarkLog) {
	var err error
	for _, logEntry := range logEntries {
		if len(channels.LogChannel) > constants.ServerLogInsertionChannelCapacity-1 {
			fmt.Printf("E#1KDZD2 - Channel is full. Cannot push. Log received: | %v\n", logEntry)
			return
		}

		logEntry, err = logEntry.ValidateForInsert()

		if err == nil {
			resources.ServerDbSaverWg.Add(1)
			channels.LogChannel <- logEntry
		}
	}
}
