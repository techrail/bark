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
	"math"
	"os"
	"strconv"
	"time"
)

// ZapLog represents a single log in uber zap log file.
// There is no guarantee that the naming we're using here will be followed.
// The naming is inspired from zap.NewProduction() Logger's output.
type ZapLog struct {
	Level         string  `json:"level"`
	UnixTimestamp float64 `json:"timestamp"`
	TextTimestamp string  `json:"timestamp"`
	Caller        string  `json:"caller"` // Not used right now.
	Message       string  `json:"msg"`
	PID           int     `json:"pid"`        // Not used right now.
	Stacktrace    string  `json:"stacktrace"` // Not used right now.
}

// isUnixTimestamp checks if a string is unix timestamp or not
func isUnixTimestamp(str string) bool {
	if _, err := strconv.ParseInt(str, 10, 64); err != nil {
		return false
	}
	return true
}

func main() {
	var numberOfLogs = 0
	filePath := flag.String("file", "log.txt", "File path to read logs from")
	url := flag.String("server", "http://localhost:8080/", "Bark server url")
	service := flag.String("service", "No service name", "Service Name")
	session := flag.String("session", "No session name", "Session Name")
	format := flag.String("format", time.RFC3339, "Time format")
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
			fmt.Printf("E#1M1AZO - Error while sending a log. E: %s\n", err.Error())
		}
		logToSend := client.RawLog{
			Message:     logToBeSent.Message,
			LogLevel:    logToBeSent.Level,
			SessionName: *session,
			ServiceName: *service,
			MoreData:    jsonObject.Typ{},
		}

		var timestamp time.Time

		if logToBeSent.UnixTimestamp == 0 {
			timestamp, _ = time.Parse(*format, logToBeSent.TextTimestamp)
		} else {
			sec, dec := math.Modf(logToBeSent.UnixTimestamp)
			timestamp = time.Unix(int64(sec), int64(dec*(1e9)))
		}

		logToSend.LogTime = timestamp

		err = logger.Raw(logToSend, true)
		if err != nil {
			fmt.Printf("E#1M1AYF - Error while sending the log. E: %s\n", err.Error())
		}
		numberOfLogs++
	}
	logger.WaitAndEnd()
	fmt.Printf("Inserted %d logs into bark.\n", numberOfLogs)
}
