package main

import (
	"fmt"
	"github.com/techrail/bark/client"
	"os"
)

func main() {
	err := os.Setenv("DEBUG", "true")
	if err != nil {
		fmt.Println("This")
		return
	}
	log := client.NewClientWithServerWithSchema("postgres://vaibhav:vaibhav@127.0.0.1:5432/bark", "audit", "INFO", "brktest", "load test session", false)

	for i := 0; i < 500; i++ {
		log.Printf("20JFH1 - schema default message - %v", i)
	}

	// Now we will disable the Debug printing
	log.DisableDebugLogs()
	log.Debug("20JF30 - This message should not be saved")
	log.EnableDebugLogs()
	log.Debug("20JF3P - This message should be saved in the table in audit schema")

	log.WaitAndEnd()
}
