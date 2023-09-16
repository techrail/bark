package models

import (
	`encoding/json`
	`fmt`
	"time"

	`github.com/techrail/bark/resources`
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
		log_time,log_level,service_name,
		code,msg,more_data
	) 
	VALUES (
	    $1, $2, $3,
	    $4, $5, $6
	)`

	_, err := resources.BarkDb.Client.Queryx(query, l.LogTime, l.LogLevel, l.ServiceName,
		l.Code, l.Message, l.MoreData)

	if err != nil {
		return fmt.Errorf("E#1KGY97 - error while inserting log: %w", err)
	}
	return nil
}

func (bld *BarkLogDao) InsertBatch(l []BarkLog) error {
	fmt.Println("E#1KGYSG - NOT YET IMPLEMENTED")
	return nil

	// // Start a transaction
	// tx, err := db.Client.Beginx()
	// if err != nil {
	// 	return fmt.Errorf("error starting a transaction: %w", err)
	// }
	//
	// query := `
	// INSERT INTO app_log
	// (log_time,log_level,service_name,code,msg,more_data)
	// VALUES (:log_time,:log_level,:service_name,:code,:msg,:more_data)
	// RETURNING id`
	//
	// result, err := tx.NamedExec(query, l)
	// if err != nil {
	// 	return fmt.Errorf("error while inserting logs: %w", err)
	// }
	//
	// // Commit the transaction
	// err = tx.Commit()
	// if err != nil {
	// 	return fmt.Errorf("error committing transaction: %w", err)
	// }
	// numRowsAffected, err := result.RowsAffected()
	// fmt.Println("Rows inserted ", numRowsAffected)
	// return err

}
