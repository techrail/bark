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
	"time"
)

// ZapLog represents a single log entry in an Uber Zap log file.
type ZapLog struct {
	Level         string  `json:"level"`      // Log level, e.g., "info", "error", etc.
	UnixTimestamp float64 `json:"timestamp"`  // Unix timestamp.
	TextTimestamp string  `json:"timestamp"`  // Text-based timestamp (alternative to Unix timestamp).
	Caller        string  `json:"caller"`     // Not used at the moment.
	Message       string  `json:"msg"`        // Log message.
	PID           int     `json:"pid"`        // Not used at the moment.
	Stacktrace    string  `json:"stacktrace"` // Not used at the moment.
}

func main() {
	var numberOfLogs = 0

	// Define command-line flags and set default values.
	filePath := flag.String("file", "log.txt", "File path to read logs from")
	server := flag.String("server", "", "Bark server URL")
	db := flag.String("db", "", "Bark database connection URL")
	service := flag.String("service", "No service name", "Service Name")
	session := flag.String("session", "No session name", "Session Name")
	format := flag.String("format", time.RFC3339, "Time format")
	help := flag.Bool("help", false, "Show help")

	// Parse command-line arguments.
	flag.Parse()

	// Read the log file specified by the -file flag.
	file, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Println("E#1M2OQ5 - Error reading file: No log file found!!")
		*help = true
	}

	// Check if either -server or -db is provided, or display the help message.
	if *server == "" && *db == "" {
		fmt.Println("E#1M4Y1P - Error connecting to the server: Need either -server or -db connection URL.")
		*help = true
	}

	// Display help instructions if -help is provided or any invalid input.
	if *help {
		fmt.Println("Usage: zap-parser [options]")
		fmt.Println("Options:")
		fmt.Println("  -file <filepath>   : Specify the path to the log file (default: log.txt).")
		fmt.Println("  -server <url>     : Set the Bark server URL (required if -db is not provided).")
		fmt.Println("  -db <url>         : Set the Bark database connection URL (required if -server is not provided).")
		fmt.Println("  -service <name>   : Specify the service name (default: No service name).")
		fmt.Println("  -session <name>   : Specify the session name (default: No session name).")
		fmt.Println("  -format <timefmt> : Set the time format (default: RFC3339).")
		fmt.Println("  -help             : Show this help message.")
		return
	}

	// Initialize the Bark client based on the provided -server or -db flag.
	var logger *client.Config
	if *server != "" {
		logger = client.NewClient(*server, constants.Info, *service, *session, false, true)
	} else {
		logger = client.NewClientWithServer(*db, constants.Info, *service, *session, false)
	}

	// Process and insert logs from the specified log file.
	scanner := bufio.NewScanner(bytes.NewReader(file))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		logToBeSent := ZapLog{}
		err := json.Unmarshal(scanner.Bytes(), &logToBeSent)
		if err != nil {
			fmt.Printf("E#1M1AZO - Error while sending a log. Error: %s\n", err.Error())
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
			fmt.Printf("E#1M1AYF - Error while sending the log. Error: %s\n", err.Error())
		}
		numberOfLogs++
	}

	// Wait for pending logs to be sent.
	logger.WaitAndEnd()

	fmt.Printf("Inserted %d logs into Bark.\n", numberOfLogs)
}
