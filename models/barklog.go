package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/techrail/bark/resources"
)

// BarkLog is a struct representing a log in Bark
type BarkLog struct {
	Id          int64           `db:"id" json:"id"`
	LogTime     time.Time       `db:"log_time" json:"logTime"`
	LogLevel    string          `db:"log_level" json:"logLevel"`
	ServiceName string          `db:"service_name" json:"serviceName"`
	SessionName string          `db:"session_name" json:"sessionName"`
	Code        string          `db:"code" json:"code"`
	Message     string          `db:"msg" json:"msg"`
	MoreData    json.RawMessage `db:"more_data" json:"moreData"`
}

func (b BarkLog) ValidateForInsert() (BarkLog, error) {
	if b.LogTime.IsZero() {
		b.LogTime = time.Now().UTC()
	}
	if strings.TrimSpace(b.LogLevel) == "" {
		b.LogLevel = "info"
	}
	if strings.TrimSpace(b.ServiceName) == "" {
		b.ServiceName = "def_svc"
	}
	if strings.TrimSpace(b.SessionName) == "" {
		b.SessionName = "def_sess"
	}

	if strings.TrimSpace(b.Code) == "" && strings.TrimSpace(b.Message) == "" {
		b.Code = "000000"
		b.Message = "_no_msg_supplied_"
		return b, fmt.Errorf("message and Code empty in Log")
	}

	if strings.TrimSpace(b.Code) == "" {
		b.Code = "000000"
	}
	if strings.TrimSpace(b.Message) == "" {
		b.Message = "_no_msg_supplied_"
	}

	if len(b.MoreData) == 0 {
		b.MoreData = json.RawMessage("{\"test\":\"value\"}")
	}

	return b, nil
}

func (b BarkLog) String() string {
	return fmt.Sprintf("Id: %v | LogTime: %v | LogLevel: %v | ServiceName: %v | SessionName: %v | Code: %v | Message: %v | MoreData: %v \n",
		b.Id, b.LogTime, b.LogLevel, b.ServiceName, b.SessionName, b.Code, b.Message, b.MoreData)
}

type BarkLogDao struct{}

func NewBarkLogDao() *BarkLogDao {
	return &BarkLogDao{}
}

// Insert inserts a Bark log in the database
func (bld *BarkLogDao) Insert(l BarkLog) error {
	query := `
	INSERT INTO app_log (
		log_time, log_level, service_name,
		session_name, code, msg, 
        more_data
	) 
	VALUES (
	    $1, $2, $3,
	    $4, $5, $6,
	    $7
	)`

	_, err := resources.BarkDb.Client.Queryx(query, l.LogTime, l.LogLevel, l.ServiceName,
		l.SessionName, l.Code, l.Message,
		l.MoreData)

	if err != nil {
		return fmt.Errorf("E#1KGY97 - error while inserting log: %w", err)
	}
	return nil
}

func (bld *BarkLogDao) InsertBatch(logs []BarkLog) error {

	query := `
	INSERT INTO app_log (
		log_time, log_level, service_name,
		session_name, code, msg, 
        more_data
	) 
	VALUES `

	numOfLogs := len(logs)
	logsToInsert := make([]interface{}, numOfLogs*7)

	for i, log := range logs {
		pos := i * 7
		logsToInsert[pos+0] = log.LogTime
		logsToInsert[pos+1] = log.LogLevel
		logsToInsert[pos+2] = log.ServiceName
		logsToInsert[pos+3] = log.SessionName
		logsToInsert[pos+4] = log.Code
		logsToInsert[pos+5] = log.Message
		logsToInsert[pos+6] = log.MoreData

		query += "(?, ?, ?, ?, ?, ?, ?)"

		if i < numOfLogs-1 {
			query += ","
		}
	}

	_, err := resources.BarkDb.Client.Queryx(query, logsToInsert)

	if err != nil {
		return fmt.Errorf("Error while inserting multiple logs: %w", err)
	}

	return nil
}
