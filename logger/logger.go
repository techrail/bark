package logger

import (
	"encoding/json"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/techrail/bark/db"
	. "github.com/techrail/bark/models"
)

func Init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type Logger struct {
	ServiceName string
}

func NewLogger(serviceName string) *Logger {
	return &Logger{
		ServiceName: serviceName,
	}
}

func (l *Logger) Log(logLevel int, code, message string, moreData json.RawMessage, database *db.BarkPostgresDb) {
	logEntry := BarkLog{
		LogTime:     time.Now(),
		LogLevel:    logLevel,
		ServiceName: l.ServiceName,
		Code:        code,
		Message:     message,
		MoreData:    moreData,
	}
	database.InsertLog(logEntry)
}

func GetLogger(serviceName string) func(int, string, string, json.RawMessage) {
	logger := NewLogger(serviceName)
	Init()
	db, err := db.ConnectToDatabase()
	if err != nil {
		log.Fatal(err)
	}
	return func(logLevel int, code, message string, moreData json.RawMessage) {
		logger.Log(logLevel, code, message, moreData, db)
	}
}

func ReadFromChannel(logChannel <-chan BarkLog) []BarkLog {
	var logsFromChannel []BarkLog
	for val := range logChannel {
		logsFromChannel = append(logsFromChannel, val)
	}
	return logsFromChannel
}

func WriteToChannel(logChannel chan BarkLog, logRecord BarkLog) {
	logChannel <- logRecord
}

// go routine to check channel length and commit to DB
func BatchCommit(logChannel chan BarkLog) string {
	logChannelLength := 0
	for {
		logChannelLength = len(logChannel)
		if logChannelLength > 100 {
			//commit in batches of 100
			logsToCommit := ReadFromChannel(logChannel)
			for i := 0; i < len(logsToCommit); i++ {
				// call bulk insert function
			}
		} else if logChannelLength > 0 && logChannelLength < 100 {
			// commit one at a time
			logsToCommit := ReadFromChannel(logChannel)
			for i := 0; i < len(logsToCommit); i++ {
				// call insert
			}
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}
