package db

import (
	"log"

	_ "github.com/lib/pq" //needed for postgres driver
	"github.com/techrail/bark/barklog"
)

func (db *Database) InsertLog(l barklog.BarkLog) error {

	// Start a transaction
	tx, err := db.Client.Beginx()
	if err != nil {
		log.Fatal(err)
		// log.Fatalf("Failed to start transaction: %v", err)
	}

	query := `INSERT INTO app_log (log_time,log_level,service_name,code,msg,more_data) VALUES (
		:log_time,:log_level,:service_name,:code,:msg,:more_data
	) RETURNING id`
	rows, err := tx.NamedQuery(query, l)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&l.Id)
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return rows.Err()

}
