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
	log := client.NewClientWithServer("postgres://vaibhav:vaibhav@127.0.0.1:5432/bark", "INFO", "brktest", "load test session", false)

	for i := 0; i < 500; i++ {
		log.Printf("1M2UBT - Default message - %v", i)
	}

	// Now we will disable the Debug printing
	log.DisableDebugLogs()
	log.Debug("1M7J7X - This message should not be saved")
	log.EnableDebugLogs()
	log.Debug("1M7J7Y - This message should be saved")

	log.WaitAndEnd()
}
