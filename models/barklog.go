package models

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/techrail/bark/appRuntime"
	"github.com/techrail/bark/constants"
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
		b.LogLevel = constants.DefaultLogLevel
	}
	if strings.TrimSpace(b.ServiceName) == "" {
		b.ServiceName = constants.DefaultLogServiceName
	}
	if strings.TrimSpace(b.SessionName) == "" {
		b.SessionName = appRuntime.SessionName
	}

	if strings.TrimSpace(b.Code) == "" && strings.TrimSpace(b.Message) == "" {
		b.Code = constants.DefaultLogCode
		b.Message = constants.DefaultLogMessage
		return b, fmt.Errorf("E#1L3VJG - message and Code empty in Log")
	}

	if strings.TrimSpace(b.Code) == "" {
		b.Code = constants.DefaultLogCode
	}
	if strings.TrimSpace(b.Message) == "" {
		b.Message = constants.DefaultLogMessage
	}

	if len(b.MoreData) == 0 {
		b.MoreData = json.RawMessage("{}")
	}

	return b, nil
}

func (b BarkLog) String() string {
	return fmt.Sprintf("Id: %v | LogTime: %v | LogLevel: %v | AppName: %v | SessionName: %v | Code: %v | Message: %v | MoreData: %v \n",
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

	_, err := resources.BarkDb.Client.Exec(context.Background(), query, l.LogTime, l.LogLevel, l.ServiceName,
		l.SessionName, l.Code, l.Message,
		l.MoreData)

	if err != nil {
		return fmt.Errorf("E#1KGY97 - error while inserting log: %w", err)
	}
	return nil
}

func (bld *BarkLogDao) InsertBatch(l []BarkLog) error {
	batchOfBarkLog := [][]any{}
	for i := 0; i < len(l); i++ {
		batchElement := []any{l[i].LogTime, l[i].LogLevel, l[i].ServiceName, l[i].SessionName,
			l[i].Code, l[i].Message, l[i].MoreData}
		batchOfBarkLog = append(batchOfBarkLog, batchElement)
	}

	_, err := resources.BarkDb.Client.CopyFrom(context.Background(), pgx.Identifier{"app_log"},
		[]string{"log_time", "log_level", "service_name", "session_name", "code", "msg", "more_data"}, pgx.CopyFromRows(batchOfBarkLog))

	if err != nil {
		return fmt.Errorf("E#1KSPLS - error while inserting batch: %w", err)
	}

	return nil
}
