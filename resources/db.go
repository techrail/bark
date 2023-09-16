package resources

import (
	`context`
	`fmt`

	`github.com/jmoiron/sqlx`
	_ `github.com/lib/pq`
)

// BarkPostgresDb wraps the sqlx.DB in a custom struct to use it as a receiver for query functions
type BarkPostgresDb struct {
	Client *sqlx.DB
}

var BarkDb *BarkPostgresDb

func InitDatabase() error {
	// Connect to Postgres DB instance
	var err error
	BarkDb, err = OpenDatabase()
	if err != nil {
		return fmt.Errorf("E#1KDZOZ - Opening database failed. Error: %v\n", err)
	}

	// Ping DB
	if err = BarkDb.Ping(context.Background()); err != nil {
		return fmt.Errorf("E#1KDZPY - Opening database failed. Error: %v\n", err)
	}

	fmt.Println("E#1KDZG7 - successfully connected to database")
	return nil
}

func OpenDatabase() (*BarkPostgresDb, error) {
	dbConn, err := sqlx.Open("postgres", "postgres://vaibhavkaushal:vaibhavkaushal@127.0.0.1:5432/bark?sslmode=disable")

	if err != nil {
		return &BarkPostgresDb{}, fmt.Errorf("E#1KDW57 - error connecting to db: %w", err)
	}

	return &BarkPostgresDb{Client: dbConn}, nil
}

func (d *BarkPostgresDb) Ping(ctx context.Context) error {
	return d.Client.DB.PingContext(ctx)
}
