package clientLogSender

import (
	"fmt"
	"github.com/techrail/bark/client/channels"
	"github.com/techrail/bark/client/network"
	"github.com/techrail/bark/constants"
	"github.com/techrail/bark/models"
	"time"
)

const logBatchSizeStandard = 100

func StartSendingLogs(serverUrl string) {
	clientChannelLength := 0
	for {
		clientChannelLength = len(channels.ClientChannel)
		var logBatch = []models.BarkLog{}
		if clientChannelLength > logBatchSizeStandard {
			// Bulk insert
			for i := 0; i < logBatchSizeStandard; i++ {
				elem, ok := <-channels.ClientChannel
				if !ok {
					fmt.Println("E# - Error occurred while getting batch from channel")
					break // Something went wrong
				}
				logBatch = append(logBatch, elem)
			}
			_, err := network.PostLogs(serverUrl+constants.BatchInsertUrl, logBatch)

			if err.Severity == 1 {
				fmt.Println(err.Msg)
			}
			fmt.Println("L# - Batch sent at ", time.Now().Format("2006-01-02 15:04:05"))
		} else if clientChannelLength > 0 && clientChannelLength < logBatchSizeStandard {
			// Commit one at a time
			singleLog := <-channels.ClientChannel
			_, err := network.PostLog(serverUrl+constants.SingleInsertUrl, singleLog)
			if err.Severity == 1 {
				fmt.Println(err.Msg)
			}
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}
