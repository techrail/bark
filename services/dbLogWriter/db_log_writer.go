package dbLogWriter

import (
	"fmt"
	"github.com/techrail/bark/constants"
	"github.com/techrail/bark/resources"
	"time"

	"github.com/techrail/bark/appRuntime"
	"github.com/techrail/bark/channels"
	"github.com/techrail/bark/models"
)

var BarkLogDao *models.BarkLogDao

func init() {
	BarkLogDao = models.NewBarkLogDao()
}

// KeepSavingLogs is a go routine to check channel length and commit to DB
// The routine decides whether a batch or single insert DB call of the logs is needed to be made.
// Further bifurcation of the batch sizes is done based on the incoming traffic and LogChannel capacity.
// If appRuntime.ShutdownRequested is set to true, the routine will send a batch of all the remaining logs in the LogChannel to the DB.
func KeepSavingLogs() {
	logChannelLength := 0
	for {
		logChannelLength = len(channels.LogChannel)
		//fmt.Println("ChanLen: ", logChannelLength)
		var logBatch = []models.BarkLog{}
		if logChannelLength >= constants.ServerLogInsertionBatchSizeLarge {
			//fmt.Println("Sending Large Batch")
			// Bulk insert
			for i := 0; i < constants.ServerLogInsertionBatchSizeLarge; i++ {
				elem, ok := <-channels.LogChannel
				if !ok {
					fmt.Println("E#1LVMFC - Error occured while getting batch from channel")
					break // Something went wrong
				}
				logBatch = append(logBatch, elem)
			}

			go func() {
				defer resources.ServerDbSaverWg.Add(-(len(logBatch)))
				err := BarkLogDao.InsertBatch(logBatch)
				if err != nil {
					fmt.Println("E#1LVMIR - Large Batch insertion failed. Error: " + err.Error() + "\n")
					for _, logEntry := range logBatch {
						fmt.Printf("E#1LVMJG - Log message: | %v\n", logEntry)
					}
					return
				}
				fmt.Println("L#1LVM50 - Large Batch inserted at ", time.Now().Format("2006-01-02 15:04:05"))
			}()
		} else if logChannelLength >= constants.ServerLogInsertionBatchSizeMedium && logChannelLength < constants.ServerLogInsertionBatchSizeLarge {
			//fmt.Println("Sending Medium Batch")
			// Bulk insert
			for i := 0; i < constants.ServerLogInsertionBatchSizeMedium; i++ {
				elem, ok := <-channels.LogChannel
				if !ok {
					fmt.Println("E#1LVMFF - Error occured while getting batch from channel")
					break // Something went wrong
				}
				logBatch = append(logBatch, elem)
			}

			go func() {
				defer resources.ServerDbSaverWg.Add(-(len(logBatch)))
				err := BarkLogDao.InsertBatch(logBatch)
				if err != nil {
					fmt.Println("E#1LVMKR - Medium Batch insertion failed. Error: " + err.Error() + "\n")
					for _, logEntry := range logBatch {
						fmt.Printf("E#1LVMKU - Log message: | %v\n", logEntry)
					}
					return
				}
				fmt.Println("L#1LVMKM - Medium Batch inserted at ", time.Now().Format("2006-01-02 15:04:05"))
			}()
		} else if logChannelLength >= constants.ServerLogInsertionBatchSizeSmall && logChannelLength < constants.ServerLogInsertionBatchSizeMedium {
			//fmt.Println("Sending Small Batch")
			// Bulk insert
			for i := 0; i < constants.ServerLogInsertionBatchSizeSmall; i++ {
				elem, ok := <-channels.LogChannel
				if !ok {
					fmt.Println("E#1LVMFL - Error occured while getting batch from channel")
					break // Something went wrong
				}
				logBatch = append(logBatch, elem)
			}

			go func() {
				defer resources.ServerDbSaverWg.Add(-(len(logBatch)))
				err := BarkLogDao.InsertBatch(logBatch)
				if err != nil {
					fmt.Println("E#1LVMLE - Small Batch insertion failed. Error: " + err.Error() + "\n")
					for _, logEntry := range logBatch {
						fmt.Printf("E#1LVMLI - Log message: | %v\n", logEntry)
					}
					return
				}
				fmt.Println("L#1LVMFR - Small Batch inserted at ", time.Now().Format("2006-01-02 15:04:05"))
			}()
		} else if logChannelLength > 0 && logChannelLength < constants.ServerLogInsertionBatchSizeSmall {
			//fmt.Println("Sending Single Log")
			// Commit one at a time
			singleLog := <-channels.LogChannel
			err := BarkLogDao.Insert(singleLog)

			resources.ServerDbSaverWg.Done()

			if err != nil {
				fmt.Println("E#1LVMIR - Individual log insertion failed. Error: " + err.Error() + "\n")
				fmt.Printf("E#1LVMML - Log message: | %v\n", singleLog)
			}
		} else {
			if appRuntime.ShutdownRequested.Load() == true {
				if len(channels.LogChannel) == 0 {
					return
				} else {
					for i := 0; i < len(channels.LogChannel); i++ {
						elem, ok := <-channels.LogChannel
						if !ok {
							fmt.Println("E#1LVMFW - Error occured while getting batch from channel")
							break // Something went wrong
						}
						logBatch = append(logBatch, elem)
					}
					err := BarkLogDao.InsertBatch(logBatch)
					if err != nil {
						fmt.Println("E#1LVMN5 - Remaining Batch insertion failed. Error: " + err.Error() + "\n")
						for _, logEntry := range logBatch {
							fmt.Printf("E#1LVMN7 - Log message: | %v\n", logEntry)
							resources.ServerDbSaverWg.Done()
						}
						return
					}
					resources.ServerDbSaverWg.Add(-(len(logBatch)))
					fmt.Println("L#1LVMN9 - Batch inserted at ", time.Now().Format("2006-01-02 15:04:05"))
				}
			} else {
				// fmt.Println("in sleep")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}
}
