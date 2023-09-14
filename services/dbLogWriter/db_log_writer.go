package dbLogWriter

import (
	`fmt`
	`time`

	`github.com/techrail/bark/db`
	`github.com/techrail/bark/models`
)

// go routine to check channel length and commit to DB
func BatchCommit(logChannel chan models.BarkLog) string {
	logChannelLength := 0
	db := new(db.BarkPostgresDb)
	for {
		logChannelLength = len(logChannel)
		if logChannelLength > 100 {
			var logBatch = []models.BarkLog{}
			for i := 0; i < 100; i++ {
				elem, ok := <-logChannel
				if !ok {
					fmt.Println("Error occured while getting batch from channel")
					break // Something went wrong
				}
				logBatch = append(logBatch, elem)
			}
			err := db.InsertBatch(logBatch)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Batch inserted at ", time.Now().Format("2006-01-02 15:04:05"))

		} else if logChannelLength > 0 && logChannelLength < 100 {
			// commit one at a time
			singleLog := <-logChannel
			err := db.InsertLog(singleLog)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Log inserted at ", time.Now().Format("2006-01-02 15:04:05"))

		} else {
			time.Sleep(1 * time.Second)
		}
	}
}
