package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/techrail/bark/client"
	"github.com/techrail/bark/constants"
	"github.com/techrail/bark/typs/jsonObject"
	"os"
	"time"
)

// ZapLog represents a single log in uber zap log file.
// There is no guarantee that the naming we're using here will be followed.
// The naming is inspired from zap.NewProduction() Logger's output.
type ZapLog struct {
	Level      string `json:"level"`
	Timestamp  string `json:"timestamp"`
	Caller     string `json:"caller"` // Not used right now.
	Message    string `json:"msg"`
	PID        int    `json:"pid"`        // Not used right now.
	Stacktrace string `json:"stacktrace"` // Not used right now.
}

func main() {
	var numberOfLogs = 0
	filePath := flag.String("file", "log.txt", "File path to read logs from")
	url := flag.String("server", "http://localhost:8080/", "Bark server url")
	service := flag.String("service", "No service name", "Service Name")
	session := flag.String("session", "No session name", "Session Name")
	layout := flag.String("layout", "2006-01-02T15:04:05.999-0700", "Time format")
	file, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Println("E#1M2OQ5 - No log file found!!")
	}
	scanner := bufio.NewScanner(bytes.NewReader(file))
	scanner.Split(bufio.ScanLines)
	logger := client.NewClient(*url, constants.Info, *service, *session, false, true)
	for scanner.Scan() {
		logToBeSent := ZapLog{}
		err := json.Unmarshal(scanner.Bytes(), &logToBeSent)
		if err != nil {
			err = fmt.Errorf("E#1M1AZO - Error while sending a log. E: %s\n", err.Error())
		}
		timestamp, err := time.Parse(*layout, logToBeSent.Timestamp)
		logToSend := client.RawLog{
			Message:     logToBeSent.Message,
			LogTime:     timestamp,
			LogLevel:    logToBeSent.Level,
			SessionName: *session,
			ServiceName: *service,
			MoreData:    jsonObject.Typ{},
		}
		err = logger.Raw(logToSend, true)
		if err != nil {
			fmt.Printf("E#1M1AYF - Error while sending the log. E: %s\n", err.Error())
		}
		numberOfLogs++
	}
	logger.WaitAndEnd()
	fmt.Printf("Inserted %d logs into bark.\n", numberOfLogs)
}
