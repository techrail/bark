package models

import (
	"context"
	"fmt"
	"github.com/techrail/bark/appRuntime"
	"github.com/techrail/bark/internal/jsonObject"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/techrail/bark/constants"
	"github.com/techrail/bark/resources"
)

// BarkLog is a struct representing a log in Bark
type BarkLog struct {
	Id                  int64          `db:"id" json:"id"`
	LogTime             time.Time      `db:"log_time" json:"logTime"`
	LogLevel            string         `db:"log_level" json:"logLevel"`
	ServiceName         string         `db:"service_name" json:"serviceName"`
	ServiceInstanceName string         `db:"service_instance_name" json:"serviceInstanceName"`
	Code                string         `db:"code" json:"code"`
	Message             string         `db:"msg" json:"msg"`
	MoreData            jsonObject.Typ `db:"more_data" json:"moreData"`
}

// ValidateForInsert checks for missing values in the incoming BarkLog's fields.
// In case a missing value is encountered, a default value is assigned to it.
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
	if strings.TrimSpace(b.ServiceInstanceName) == "" {
		b.ServiceInstanceName = constants.DefaultLogSessionName
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

	if b.MoreData.IsEmpty() {
		b.MoreData = jsonObject.EmptyNotNullJsonObject()
	}

	return b, nil
}

func (b BarkLog) String() string {
	return fmt.Sprintf("Id: %v | LogTime: %v | LogLevel: %v | AppName: %v | SessionName: %v | Code: %v | Message: %v | MoreData: %v \n",
		b.Id, b.LogTime, b.LogLevel, b.ServiceName, b.ServiceInstanceName, b.Code, b.Message, b.MoreData)
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
		service_instance_name, code, msg, 
        more_data
	) 
	VALUES (
	    $1, $2, $3,
	    $4, $5, $6,
	    $7
	)`

	_, err := resources.BarkDb.Client.Exec(context.Background(), query, l.LogTime, l.LogLevel, l.ServiceName,
		l.ServiceInstanceName, l.Code, l.Message,
		l.MoreData)

	if err != nil {
		return fmt.Errorf("E#1KGY97 - error while inserting log: %w", err)
	}
	return nil
}

// InsertServerStartedLog inserts a log entry in the postgres DB stating that bark server has started successfully.
// This acts as a checkpoint that everything is working as expected in the DB connection department.
func (bld *BarkLogDao) InsertServerStartedLog() error {
	return bld.Insert(BarkLog{
		LogTime:             time.Now().UTC(),
		LogLevel:            constants.Info,
		ServiceName:         "Bark Server",
		ServiceInstanceName: appRuntime.SessionName,
		Code:                "1LQ2X3",
		Message:             "Server started",
		MoreData:            jsonObject.EmptyNotNullJsonObject(),
	})
}

// InsertBatch sends a batch of logs to the DB.
func (bld *BarkLogDao) InsertBatch(l []BarkLog) error {
	batchOfBarkLog := [][]any{}
	for i := 0; i < len(l); i++ {
		batchElement := []any{l[i].LogTime, l[i].LogLevel, l[i].ServiceName, l[i].ServiceInstanceName,
			l[i].Code, l[i].Message, l[i].MoreData}
		batchOfBarkLog = append(batchOfBarkLog, batchElement)
	}

	_, err := resources.BarkDb.Client.CopyFrom(context.Background(), pgx.Identifier{"app_log"},
		[]string{"log_time", "log_level", "service_name", "service_instance_name", "code", "msg", "more_data"}, pgx.CopyFromRows(batchOfBarkLog))

	if err != nil {
		return fmt.Errorf("E#1KSPLS - error while inserting batch: %w", err)
	}

	return nil
}

func (bld *BarkLogDao) FetchLogs(logLevel, serviceName, sessionName, startDate, endDate string) ([]BarkLog, error) {
	query := "SELECT * FROM logs WHERE 1=1" // Base query
	if logLevel != "" {
		query += fmt.Sprintf(" AND log_level = '%s'", logLevel)
	}
	if serviceName != "" {
		query += fmt.Sprintf(" AND service_name = '%s'", serviceName)
	}
	if sessionName != "" {
		query += fmt.Sprintf(" AND session_name = '%s'", sessionName)
	}
	if startDate != "" {
		query += fmt.Sprintf(" AND log_time >= '%s'", startDate)
	}
	if endDate != "" {
		query += fmt.Sprintf(" AND log_time <= '%s'", endDate)
	}
	rows, err := resources.BarkDb.Client.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	logs := make([]BarkLog, 0)
	for rows.Next() {
		var log BarkLog
		if err := rows.Scan(&log.Id, &log.LogTime, &log.LogLevel, &log.ServiceName, &log.SessionName, &log.Code, &log.Message, &log.MoreData); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}
