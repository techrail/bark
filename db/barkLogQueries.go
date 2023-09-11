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

// Fetches all logs within a transaction
func (db *BarkPostgresDb) FetchAllLogs() ([]models.BarkLog, error) {
	var barkLogs []models.BarkLog

	query := `
	SELECT id,log_time,log_level,service_name,code,msg,more_data
	FROM app_log`

	err := db.Client.Select(&barkLogs, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching log rows: %w", err)
	}

	return barkLogs, nil

}
