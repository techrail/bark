package db

import (
	"fmt"

	_ "github.com/lib/pq" //needed for postgres driver
	"github.com/techrail/bark/models"
)

// Inserts a Bark log within a transaction
func (db *BarkPostgresDb) InsertLog(l models.BarkLog) error {

	// Start a transaction
	tx, err := db.Client.Beginx()
	if err != nil {
		return fmt.Errorf("error starting a transaction: %w", err)
	}

	query := `
	INSERT INTO app_log 
	(log_time,log_level,service_name,code,msg,more_data) 
	VALUES (:log_time,:log_level,:service_name,:code,:msg,:more_data) 
	RETURNING id`

	rows, err := tx.NamedQuery(query, l)
	if err != nil {
		return fmt.Errorf("error while inserting logs: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&l.Id)
		if err != nil {
			return fmt.Errorf("error while inserting logs: %w", err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return rows.Err()

}

// Fetches multiple logs, the number (default 10) is configurable through the limit parameter
func (db *BarkPostgresDb) FetchLimitedLogs(rowLimit int) ([]models.BarkLog, error) {
	var actualRowLimit int = 10 //default is 10, if its specified through params, the limit is overwritten
	if rowLimit > 0 {
		actualRowLimit = rowLimit
	}
	var barkLogs []models.BarkLog

	query := `
	SELECT id,log_time,log_level,service_name,code,msg,more_data
	FROM app_log LIMIT $1;`

	err := db.Client.Select(&barkLogs, query, actualRowLimit)
	if err != nil {
		return nil, fmt.Errorf("error fetching log rows: %w", err)
	}

	return barkLogs, nil

}
