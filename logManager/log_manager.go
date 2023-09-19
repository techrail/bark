// Package logManager provides a logging framework for managing log messages.
package logManager

import (
	. "encoding/json"
	"fmt"
	"github.com/google/uuid"
	. "github.com/techrail/bark/models"
	"github.com/techrail/bark/services/ingestion"
	"time"
)

// LoggerSession represents a session for logging messages.
type LoggerSession struct {
	dbConnectionStr string // Database connection string
	serviceName     string // Name of the service
	sessionID       string // Unique session identifier
}

func GetLogger(dbConnectionStr, serviceName string) *LoggerSession {
	return &LoggerSession{dbConnectionStr: dbConnectionStr, serviceName: serviceName, sessionID: uuid.NewString()}
}

func (logger *LoggerSession) log(level, message, code string, customFields ...map[string]string) {
	var moreData = RawMessage{}

	if len(customFields) > 0 {
		moreData = getJSONData(customFields[0])
	}

	log := BarkLog{
		LogTime:     time.Now(),
		LogLevel:    level,
		ServiceName: logger.serviceName,
		SessionName: logger.sessionID,
		Code:        code,
		Message:     message,
		MoreData:    moreData,
	}

	//go ingestion.InsertSingle(log)
	NewBarkLogDao().Insert(log)
}

func (logger *LoggerSession) Info(message, code string, customFields ...map[string]string) {
	logger.log("INFO", message, code, customFields...)
}

func (logger *LoggerSession) Warn(message, code string, customFields ...map[string]string) {
	logger.log("WARN", message, code, customFields...)
}

func (logger *LoggerSession) Debug(message, code string, customFields ...map[string]string) {
	logger.log("DEBUG", message, code, customFields...)
}

func (logger *LoggerSession) Error(message, code string, customFields ...map[string]string) {
	logger.log("ERROR", message, code, customFields...)
}

func (logger *LoggerSession) Fatal(message, code string, customFields ...map[string]string) {
	logger.log("FATAL", message, code, customFields...)
}

func (logger *LoggerSession) Custom(barkLog BarkLog) {
	go ingestion.InsertSingle(barkLog)
}

func getJSONData(data map[string]string) RawMessage {
	jsonData, err := Marshal(data)
	if err != nil {
		_ = fmt.Errorf("Error while parsing custom data: %v\n", err)
		return RawMessage("{}")
	}
	return jsonData
}
