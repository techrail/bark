package dbLogWriter

import (
	"fmt"
	"time"

	"github.com/techrail/bark/appRuntime"
	"github.com/techrail/bark/channels"
	"github.com/techrail/bark/models"
)

var BarkLogDao *models.BarkLogDao

func init() {
	BarkLogDao = models.NewBarkLogDao()
}

// StartWritingLogs is a go routine to check channel length and commit to DB
func StartWritingLogs() {
	logChannelLength := 0
	for {
		logChannelLength = len(channels.LogChannel)
		var logBatch = []models.BarkLog{}
		logBatchSize := 100
		if logChannelLength > 100 {
			if appRuntime.ShutdownRequested.Load() {
				logBatchSize = len(channels.LogChannel)
			}
			// Bulk insert
			for i := 0; i < logBatchSize; i++ {
				elem, ok := <-channels.LogChannel
				if !ok {
					fmt.Println("Error occured while getting batch from channel")
					break // Something went wrong
				}
				logBatch = append(logBatch, elem)
			}
			err := BarkLogDao.InsertBatch(logBatch)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Batch inserted at ", time.Now().Format("2006-01-02 15:04:05"))
		} else if logChannelLength > 0 && logChannelLength < 100 {
			// Commit one at a time
			singleLog := <-channels.LogChannel
			err := BarkLogDao.Insert(singleLog)

			if err != nil {
				fmt.Println(err)
			}
		} else {
			// fmt.Println("in sleep")
			time.Sleep(1 * time.Second)
		}
	}
}
