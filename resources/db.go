package resources

import (
	`context`
	`fmt`

	`github.com/techrail/bark/db`
)

func InitDatabase() error {
	// Connect to Postgres DB instance
	barkDb, err := db.OpenDatabase()
	if err != nil {
		return fmt.Errorf("E#1KDZOZ - Opening database failed. Error: %v\n", err)
	}
	// Ping DB
	if err = barkDb.Ping(context.Background()); err != nil {
		return fmt.Errorf("E#1KDZPY - Opening database failed. Error: %v\n", err)
	}

	fmt.Println("E#1KDZG7 - successfully connected to database")
	return nil
}
