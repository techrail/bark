package db

import (
	"log"

	_ "github.com/lib/pq" //needed for postgres driver
	"github.com/techrail/bark/logger"
)

// Add inserting student record using the sqlx package
// func (db *Database) addLog(logger.Log) (int64, error) {
// 	query := "insert into app_log (fname, lname, date_of_birth, email, gender, address) values (?, ?, ?, ?, ?, ?);"
// 	result := db.Client.MustExec(query, s.Fname, s.Lname, s.DateOfBirth, s.Email, s.Gender, s.Address)
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, fmt.Errorf("addStudent Error: %v", err)
// 	}
// 	return id, nil
// }

func (db *Database) InsertLog(l logger.Log) error {

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
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&l.Id)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return rows.Err()

}
