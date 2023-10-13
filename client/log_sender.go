package client

import (
	"fmt"
	"github.com/techrail/bark/constants"
	"github.com/techrail/bark/models"
	"time"
)

const logBatchSizeSmall = 10
const logBatchSizeMedium = 100
const logBatchSizeLarge = 500

func keepSendingLogs(serverUrl string) {
	clientChannelLength := 0
	for {
		clientChannelLength = len(PendingLogsChan)
		var logBatch = []models.BarkLog{}
		if clientChannelLength >= logBatchSizeLarge {
			// Bulk insert large
			for i := 0; i < logBatchSizeLarge; i++ {
				elem, ok := <-PendingLogsChan
				if !ok {
					fmt.Println("E#1LUFJ8 - Error occurred while getting large batch from channel")
					break // Something went wrong
				}
				logBatch = append(logBatch, elem)
			}
			go func() {
				defer Wg.Add(-logBatchSizeLarge)
				_, err := PostLogArray(serverUrl+constants.BatchInsertUrl, logBatch)

				if err.Severity == 1 {
					fmt.Println("E#1LUH90 - Large Batch insertion failed. Error: " + err.Msg + "\n")
					for _, logEntry := range logBatch {
						fmt.Printf("E#1LUHCR - Log message: | %v\n", logEntry)
					}
					return
				}

				fmt.Println("L#1LUE2H - Large Batch sent at ", time.Now().Format("2006-01-02 15:04:05"))
			}()

		} else if clientChannelLength >= logBatchSizeMedium && clientChannelLength < logBatchSizeLarge {
			// Bulk insert large
			for i := 0; i < logBatchSizeMedium; i++ {
				elem, ok := <-PendingLogsChan
				if !ok {
					fmt.Println("E#1LUFJC - Error occurred while getting medium batch from channel")
					break // Something went wrong
				}
				logBatch = append(logBatch, elem)
			}

			go func() {
				defer Wg.Add(-logBatchSizeMedium)
				_, err := PostLogArray(serverUrl+constants.BatchInsertUrl, logBatch)

				if err.Severity == 1 {
					fmt.Println("E#1LUHEZ - Medium Batch insertion failed. Error: " + err.Msg + "\n")
					for _, logEntry := range logBatch {
						fmt.Printf("E#1LUHF3 - Log message: | %v\n", logEntry)
					}
					return
				}

				fmt.Println("L#1LUE2H - Medium Batch sent at ", time.Now().Format("2006-01-02 15:04:05"))
			}()
		} else if clientChannelLength >= logBatchSizeSmall && clientChannelLength < logBatchSizeMedium {
			// Bulk insert large
			for i := 0; i < logBatchSizeSmall; i++ {
				elem, ok := <-PendingLogsChan
				if !ok {
					fmt.Println("E#1LUFTF - Error occurred while getting small batch from channel")
					break // Something went wrong
				}
				logBatch = append(logBatch, elem)
			}

			go func() {
				defer Wg.Add(-logBatchSizeSmall)
				_, err := PostLogArray(serverUrl+constants.BatchInsertUrl, logBatch)

				if err.Severity == 1 {
					fmt.Println("E#1LUHFT - Small Batch insertion failed. Error: " + err.Msg + "\n")
					for _, logEntry := range logBatch {
						fmt.Printf("E#1LUHFV - Log message: | %v\n", logEntry)
					}
					return
				}

				fmt.Println("L#1LUE2H - Medium Batch sent at ", time.Now().Format("2006-01-02 15:04:05"))
			}()
		} else if clientChannelLength > 0 && clientChannelLength < logBatchSizeSmall {
			// Commit one at a time
			singleLog, ok := <-PendingLogsChan
			if !ok {
				fmt.Println("E#1LUHNW - Error occurred while getting single log from channel")
			}

			go func() {
				defer Wg.Done()
				_, err := PostLog(serverUrl+constants.SingleInsertUrl, singleLog)

				if err.Severity == 1 {
					fmt.Println("E#1LUHW8 - Single log entry insertion failed. Error: " + err.Msg + "\n")
					fmt.Printf("E#1LUHWI - Log message: | %v\n", singleLog)
					return
				}

				fmt.Println("L#1LUHWR - Single insertion sent at ", time.Now().Format("2006-01-02 15:04:05"))
			}()
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}
