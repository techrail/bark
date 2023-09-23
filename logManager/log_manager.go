// Package logManager provides a logging framework for managing log messages.
package logManager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	. "github.com/techrail/bark/models"
	"net/http"
	"time"
)

// LoggerSession represents logging object.
type LoggerSession struct {
	serverUrl   string
	serviceName string
	sessionID   string
}

type JsonLogBody struct {
	log_time     string
	log_level    string
	service_name string
	session_name string
	code         string
	msg          string
	more_data    string
}

func GetLogger(serverUrl, serviceName string) *LoggerSession {
	return &LoggerSession{serverUrl: serverUrl, serviceName: serviceName, sessionID: uuid.NewString()}
}

func (logger *LoggerSession) log(level string, message *string, fields ...any) {

	log := JsonLogBody{
		log_time:     time.Now().String(),
		log_level:    level,
		msg:          *message,
		service_name: logger.serviceName,
		session_name: logger.sessionID,
	}

	switch len(fields) {
	case 1:
		if code, ok := fields[0].(string); ok {
			log.code = code
		} else {
			fmt.Println("E#4Z4PWS unable to read code")
			return
		}
	case 2:
		if code, ok := fields[0].(string); ok {
			log.code = code
		} else {
			fmt.Println("E#TCCT6G unable to read code")
			return
		}
		if more, ok := fields[1].(json.RawMessage); ok {
			log.more_data = string(more)
		} else {
			fmt.Println("E#D8KY3C unable to read json data")
			return
		}
	}

	payload, err := json.Marshal(log)

	if err != nil {
		fmt.Printf("E#4Z4PWS %v\n", err)
	}

	fmt.Println(string(payload))

	post, err := http.Post(logger.serverUrl+"/insertSingle", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Println(post)
}

func (logger *LoggerSession) Info(message string, fields ...any) {
	logger.log("INFO", &message, fields...)
}

func (logger *LoggerSession) Warn(message string, fields ...any) {
	logger.log("WARN", &message, fields...)
}

func (logger *LoggerSession) Debug(message string, fields ...any) {
	logger.log("DEBUG", &message, fields...)
}

func (logger *LoggerSession) Error(message string, fields ...any) {
	logger.log("ERROR", &message, fields...)
}

func (logger *LoggerSession) Fatal(message string, fields ...any) {
	logger.log("FATAL", &message, fields...)
}

func (logger *LoggerSession) Custom(barkLog *BarkLog) {
	NewBarkLogDao().Insert(*barkLog)
}
