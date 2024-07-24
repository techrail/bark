package models

import (
	"context"
	"fmt"
	"github.com/techrail/bark/appRuntime"
	"github.com/techrail/bark/config"
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
	query := fmt.Sprintf(`
	INSERT INTO %vapp_log (
		log_time, log_level, service_name,
		service_instance_name, code, msg, 
        more_data
	) 
	VALUES (
	    $1, $2, $3,
	    $4, $5, $6,
	    $7
	)`, config.DbSchemaNameWithDot)

	fmt.Println("------------------")
	fmt.Println(query)
	fmt.Println("------------------")

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

	if config.DbSchemaName == "" {
		_, err := resources.BarkDb.Client.CopyFrom(context.Background(), pgx.Identifier{"app_log"},
			[]string{"log_time", "log_level", "service_name", "service_instance_name", "code", "msg", "more_data"}, pgx.CopyFromRows(batchOfBarkLog))
		if err != nil {
			return fmt.Errorf("E#1KSPLS - error while inserting batch: %w", err)
		}
	} else {
		_, err := resources.BarkDb.Client.CopyFrom(context.Background(), pgx.Identifier{config.DbSchemaName, "app_log"},
			[]string{"log_time", "log_level", "service_name", "service_instance_name", "code", "msg", "more_data"}, pgx.CopyFromRows(batchOfBarkLog))
		if err != nil {
			return fmt.Errorf("E#20JFGC - error while inserting batch: %w", err)
		}
	}

	return nil
}
